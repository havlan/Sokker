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

package uferdig

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	sokk := new_web_sokker()

	go sokk.Init()
	//k := hand("dGhlIHNhbXBsZSBub25jZQ==" + magic_server_key)
	//fmt.Println(k)
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe(":3000", nil)
}

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "3001"
	CONN_TYPE        = "tcp"
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	OP_Continue      = 0
	OP_Text          = 1
	OP_Binary        = 2
	OP_Close         = 8
	OP_Ping          = 9
	OP_Pong          = 10
)

type web_sokker struct {
	clients []*client_
	joins   chan net.Conn
	inc     chan string
	out     chan string
}

func new_web_sokker() *web_sokker {
	ws := &web_sokker{
		clients: make([]*client_, 0),
		joins:   make(chan net.Conn),
		inc:     make(chan string),
		out:     make(chan string),
	}
	ws.Listen()
	return ws
}

func (ws *web_sokker) Init() {
	listener, err := net.Listen("tcp", "localhost:3001")
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	//Executed when the application closes.
	defer listener.Close()
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new go
		go ws.handler(conn)
	}

}

func (ws *web_sokker) Broadcast(d string) {

	//TODO IMPLEMENT WEBSOCKET FRAME HERE

	for _, c := range ws.clients {
		c.out <- d
	}
}

//adds a client to connected, and connects incoming ws messages to the client.
/*
func (*web_sokker) Add(net.Conn)
takes a net connection and adds to the client list
*/
func (ws *web_sokker) Add(c net.Conn) {
	new_client := new_client_(c)
	ws.clients = append(ws.clients, new_client) // new client
	//lambda
	go func() {
		for {
			ws.inc <- <-new_client.inc // add new client incoming to ws inc
		}
	}()
}

/*
func (ws *web_sokker) Listen()
listens to the channels, if the channel gets data. Send that data to decode => frame => broadcast
*/
func (ws *web_sokker) Listen() {
	go func() {
		for {
			select {
			//case http := <-  // ??????

			case data := <-ws.inc: // disse lytter til data
				log.Println(data)
				//ws.Broadcast(data)

				//case conn := <- ws.joins:
				//ws.Add(conn)
			}
		}
	}()
}

func (ws *web_sokker) startWss() {
	log.Println("Listen for incoming connections.")
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	//Executed when the application closes.
	defer listener.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new thread (goroutine)
		go ws.handler(conn)
	}
}

/*
func handler(net.Conn)
initiates the handshake process which is required
handler does handshake
*/
func (ws *web_sokker) handler(client net.Conn) {
	ws.handshake(client)
}

/*
func magic_str(in string) (key string)
takes the key from the clients request, and appends the wskey
sha1 that sum and returns that string
*/

func magic_str(str string) (keyz string) {
	h := sha1.New()
	h.Write([]byte(str))
	keyz = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

//unused?
func recv_data(client net.Conn) {
	reply := make([]byte, 64)
	client.Read(reply)
	upcodeInt := opcode(reply)
	if upcodeInt == 1 {
		log.Println("avbryter klient")
		client.Close()
	} else {
		log.Println("fortsetter klient")
		decoded := decode(reply)
		encoded := encode(decoded)
		client.Write(encoded)
		recv_data(client) // echo for debugging
	}
}

/*
func handshake(net.Conn)
Sends a client to parse the key, it either gets rejected(bad request) or accepted => 101 status code
101 statuscode is switching protocols, because we are going over to websockets
*/

func (ws *web_sokker) handshake(client net.Conn) {
	status, key := parseKey(client)
	if status != 101 {
		//reject
		reject(client)
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
		log.Println(key)
		//recv_data(client)
		ws.Add(client)

	}
}

/*
func parseKey(net.Conn) (httpStatus int, errcode string)
Parses the header first sent by client
Returns http statuscodes and a string/errstring
*/

func parseKey(client net.Conn) (code int, k string) {
	bufReader := bufio.NewReader(client)
	request, err := http.ReadRequest(bufReader)
	if err != nil {
		log.Println(err)
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
