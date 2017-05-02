/*
THE GOAL IS TO PARSE THIS
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/63)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
*/

package sokk

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"os"
	"fmt"
)

const (
	CONN_TYPE        = "tcp"
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	OP_Continue      = 0
	OP_Text          = 1
	OP_Binary        = 2
	OP_Close         = 8
	OP_Ping          = 9
	OP_Pong          = 10
)

func main() {
	sokk := NewSokk()
	
	sokk.OnClose = func(){
		fmt.Println("OnClose!")
	}
	sokk.OnConnection = func(c net.Conn){
		fmt.Println("NEW CONNECTION!")
		sokk.Clients = append(sokk.Clients,c) // add the user into the accepted client list
	}
	sokk.OnError = func(w string, e error){ // custom handle error
		fmt.Println(w, " ", e.Error())
		os.Exit(1)
		
	}
	sokk.OnMessage = func(b SokkMsg){
		fmt.Println(string(b.Payload[:b.PlLen])) // prints the data
		sokk.Send(&b)                             // sends to all Clients which exists in the sockets array of connections
		
	}
	go sokk.Start("127.0.0.1", "3001") // localhost:3000
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe("localhost:3000", nil)
}

type Sokk struct {
	Clients      []net.Conn
	OnMessage    func(msg SokkMsg)
	OnError      func (w string,e error)
	OnConnection func (c net.Conn)
	OnClose      func ()
}


func NewSokk() *Sokk {
	sokk := &Sokk{
		Clients: make([]net.Conn, 0),
	}
	return sokk
}

func (ws *Sokk) Start(ad string, port string) {
	listener, err := net.Listen(CONN_TYPE, ad+":"+port)
	if err != nil {
		ws.OnError("Init listen error.",err)
	}
	//Executed when the application closes.
	defer listener.Close()
	log.Println("WE LIVE " + ad + ":" + string(port))
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			ws.OnError("Accept connection error.",err)
		}
		// Handle connections in a new thread (goroutine)
		go ws.handler(&conn)
	}
}

func (ws *Sokk) Close() {
	for i := range ws.Clients {
		ws.Clients[i].Close()
	}
}

/*
	func handler(net.Conn)
	initiates the handshake process which is required
	handler does handshake
	starts the read loop from the user
*/
func (ws *Sokk)handler(c *net.Conn) {
	ok := ws.handshake((*c))
	if ok {
		for {
			buff := make([]byte, 512)
			(*c).Read(buff)
			var opC = int(0x7F & buff[0])
			//fmt.Println(opC)
			if opC == 8 {
				fmt.Println("CLOSE")
				ws.close_r(*c)
				break
			}else if opC == 9 {
				response := make([]byte,2)
				response[0] = byte(138)
				(*c).Write(response)
			}else{ // create wsframe
				go ws.prep_msg(buff)
			}
			
		}
	}
}

func (ws *Sokk) prep_msg(buff []byte){
	var frame = decode(buff)
	ws.OnMessage(*frame)
	//var buffTosend = encode(frame)
	//ws.Send(buffTosend)
}

/*
	func Send(bytes[])
	sends a websocket frame to all the Clients
*/

func (ws *Sokk) Send(m *SokkMsg) {
	var msg = encode(m)
	fmt.Println(len(ws.Clients))
	for i := range ws.Clients {
		ws.Clients[i].Write(msg)
	}
	
}

/*
	func handshake(net.Conn)
	Sends a client to parse the key, it either gets rejected(bad request) or accepted => 101 status code
	101 statuscode is switching protocols, because we are going over to websockets
*/

func (ws *Sokk) handshake(client net.Conn) bool {
	status, key := ws.parseKey(client)
	if status != 101 {
		//reject
		reject(client)
		return false
	} else {
		//Complete handshake
		var t = magic_str(key + magic_server_key)
		var buff bytes.Buffer
		buff.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
		buff.WriteString("Connection: Upgrade\r\n")
		buff.WriteString("Upgrade: websocket\r\n")
		buff.WriteString("Sec-WebSocket-Accept:")
		buff.WriteString(t + "\r\n\r\n")
		client.Write(buff.Bytes())
		ws.OnConnection(client)
		return true
	}
}

/*
	func parseKey(net.Conn) (httpStatus int, errcode string)
	Parses the header first sent by client
	Returns http statuscodes and a string/errstring
*/

func (ws *Sokk) parseKey(client net.Conn) (code int, k string) {
	bufReader := bufio.NewReader(client) // TODO Double trouble? this coud very well be a client instead
	request, err := http.ReadRequest(bufReader)
	if err != nil {
		ws.OnError("Parse HTTP request error. ",err)
	}
	if request.Header.Get("Upgrade") != "websocket" {
		return http.StatusBadRequest, ""
	} else {
		key := request.Header.Get("Sec-Websocket-Key")
		return http.StatusSwitchingProtocols, key
	}
}

/*
	Client did not pass upgrade handshake, so the client gets rejected
	Returns a standard request, with bad request as status
*/
func reject(client net.Conn) {
	var buff bytes.Buffer
	buff.WriteString("HTTP/1.1 400 Bad Request\r\n")
	buff.WriteString("Content-Type: text/plain\r\n")
	buff.WriteString("Connection: close\r\n\r\nIncorrect request")
	client.Write(buff.Bytes())
	client.Close()
}

/*
	func magic_str(in string) (key string)
	takes the key from the Clients request, and appends the wskey
	sha1 that sum and returns that string
*/
func magic_str(str string) (keyz string) {
	h := sha1.New()
	h.Write([]byte(str))
	keyz = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

/*
	func close_r(net.Conn)
	removes the connection from the list. If the user count is low,
*/
func (ws *Sokk) close_r(c net.Conn) {
	for i := range ws.Clients {
		if ws.Clients[i] != nil{
			if ws.Clients[i] == c {
				c.Close()
				ws.Clients = append(ws.Clients[:i], ws.Clients[i+1:]...) // delete from slice
				break
			}
		}
	}
	ws.OnClose()
}
