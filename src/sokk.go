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

package main

import (
	"net"
	"net/http"
	"os"
	"bufio"
	"bytes"
	"encoding/base64"
	"crypto/sha1"
	"log"
)


const (
	CONN_HOST = "localhost"
	CONN_PORT = "3001"
	CONN_TYPE = "tcp"
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	OP_Continue = 0
	OP_Text = 1
	OP_Binary = 2
	OP_Close = 8
	OP_Ping = 9
	OP_Pong = 10
)

func main (){
	sokk := newSokk()

	go sokk.startWss()
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe(":3000", nil)
}


type sokk struct {
	clients []net.Conn

}

func newSokk() *sokk{
	sokk :=&sokk{
		clients:make([]net.Conn,0),
	}
	return sokk
}
/*
func (*web_sokker) Add(net.Conn)
takes a net connection and adds to the client list
 */
func (ws *sokk) Add(c net.Conn){
	//new_client:= new_client_(c)
	ws.clients = append(ws.clients, c) // new client
	log.Println(c.RemoteAddr()," connected, ",len(ws.clients)," client[s] connected now.")
}

func(ws *sokk) startWss() {
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	//Executed when the application closes.
	defer listener.Close()
	log.Println("WE LIVE " + CONN_HOST + ":" + CONN_PORT)
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
starts the read loop from the user
 */
func(ws *sokk) handler(c net.Conn) {
	ok := ws.handshake(c)
	if ok{
		for{
			buff := make([]byte,512)
			c.Read(buff)
			log.Println(string(buff[:]))
			opcode := opcode(buff)
			switch opcode {
			case OP_Continue: // 0
				log.Println("OP Code Continue")
				data := decode(buff)
				go ws.sendData(encode(data))

			/*case OP_Text: // 1
				log.Println("OP Code Text")

			case OP_Binary: // 2
				log.Println("OP Code Binary")

			*/
			case OP_Close: // 8
				log.Println("OP Code Close")
				c.Close()
				break
			/*

			case OP_Ping:  // 9
				log.Println("Op Code Ping")

			case OP_Pong: // 10
				log.Println("OP Code Pong")*/
			default:
				log.Println("No familiar op code from client:",opcode)
				data := decode(buff)

				go ws.sendData(encode(data))

			}

		}
	}
}

/*
func sendData(bytes[])
sends a websocket frame to all the clients
 */

func(ws *sokk) sendData(buff []byte){
	for i := range ws.clients{
		ws.clients[i].Write(buff)
	}
}

/*
func handshake(net.Conn)
Sends a client to parse the key, it either gets rejected(bad request) or accepted => 101 status code
101 statuscode is switching protocols, because we are going over to websockets
 */

func(ws *sokk) handshake(client net.Conn) bool{
	status, key := parseKey(client)
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
		log.Println(key)
		//recv_data(client)
		ws.Add(client)
		return true

	}
}
/*
func parseKey(net.Conn) (httpStatus int, errcode string)
Parses the header first sent by client
Returns http statuscodes and a string/errstring
 */

func parseKey(client net.Conn) (code int, k string) {
	bufReader := bufio.NewReader(client) // TODO Double trouble? this coud very well be a client instead
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
/*
func magic_str(in string) (key string)
takes the key from the clients request, and appends the wskey
sha1 that sum and returns that string
 */
func magic_str(str string)(keyz string){
	h:=sha1.New()
	h.Write([]byte(str))
	keyz = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}