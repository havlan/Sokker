/*
THE GOAL IS TO PARSE THIS
From https://tools.ietf.org/html/rfc6455
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
)

const (
	CONN_TYPE        = "tcp"
	Magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	OP_Continue      = 0
	OP_Text          = 1
	OP_Binary        = 2
	OP_Close         = 8
	OP_Ping          = 9
	OP_Pong          = 10
)

type Sokk struct {
	Clients      []net.Conn
	OnMessage    func(msg SokkMsg)
	OnError      func (w string,e error)
	OnConnection func (c net.Conn)
	OnClose      func (c net.Conn)
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
			_,err :=(*c).Read(buff)
			if(err != nil){
				ws.Close_wErr(*c)
				ws.OnError("Error reading from client. ",err)
				break
			}
			var opC = int(0x7F & buff[0])
			if opC == 8 {
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
/*
	func prep_msg(buff[]byte)
	creates a websocket frame of incoming data then calls the user specified OnMessage method.
*/
func (ws *Sokk) prep_msg(buff []byte){
	var frame = Decode(buff)
	ws.OnMessage(*frame)
	//var buffTosend = Encode(frame)
	//ws.Send(buffTosend)
}

/*
	func Send(bytes[])
	sends a websocket frame to all the Clients
*/
func (ws *Sokk) Send(m *SokkMsg) {
	var msg = Encode(m)
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
		var t = Magic_str(key + Magic_server_key)
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
	bufReader := bufio.NewReader(client)
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
	func Magic_str(in string) (key string)
	takes the key from the Clients request, and appends the wskey
	sha1 that sum and returns that string
*/
func Magic_str(str string) (keyz string) {
	h := sha1.New()
	h.Write([]byte(str))
	keyz = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
/*
	func Close_wErr(c net.Conn)
	removes a user from the user list. Without calling the user specified method OnClose
 */

func (ws *Sokk) Close_wErr(c net.Conn){
	for i := range ws.Clients {
		if ws.Clients[i] != nil{
			if ws.Clients[i] == c {
				c.Close()
				ws.Clients = append(ws.Clients[:i], ws.Clients[i+1:]...) // delete from slice
				break
			}
		}
	}
}

/*
	func close_r(net.Conn)
	calls the user specified onClose method
	removes the connection from the list.
*/
func (ws *Sokk) close_r(c net.Conn) {
	ws.OnClose(c)
	for i := range ws.Clients {
		if ws.Clients[i] != nil{
			if ws.Clients[i] == c {
				c.Close()
				ws.Clients = append(ws.Clients[:i], ws.Clients[i+1:]...) // delete from slice
				break
			}
		}
	}
}
