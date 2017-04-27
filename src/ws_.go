package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"bufio"
	"bytes"
	"encoding/base64"
	"crypto/sha1"
	//"net/textproto"
	//"regexp"
	//"strings"
)

func main() {
	go startWss()
	//k := hand("dGhlIHNhbXBsZSBub25jZQ==" + magic_server_key)
	//fmt.Println(k)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3001"
	CONN_TYPE = "tcp"
	magic_server_key = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)
type web_sokker struct {
	//map[]
}

var p = fmt.Println

func startWss() {
	p("Listen for incoming connections.")
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		p("Error listening:", err.Error())
		os.Exit(1)
	}
	//Executed when the application closes.
	defer listener.Close()
	p("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			p("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new thread (goroutine)
		go handler(conn)
	}
}

// Handles incoming requests.
func handler(client net.Conn) {
	handshake(client)
}
func hand(str string)(keyz string){
	h:=sha1.New()
	h.Write([]byte(str))
	keyz = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

func recv_data(client net.Conn){
	p("LISTEN TO recv_data")
	reply := make([]byte, 32)
	client.Read(reply)
	decoded := decode(reply)
	fmt.Println("Message Received:", decoded)
	client.Write(reply)
	//client.Close()
}

func handshake(client net.Conn) {
	status, key := parseKey(client)
	if status != 101 {
		//reject
		reject(client)
	} else {
		//Complete handshake
		var t = hand(key + magic_server_key)
		var buff bytes.Buffer
		buff.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
		buff.WriteString("Connection: Upgrade\r\n")
		buff.WriteString("Upgrade: websocket\r\n")
		buff.WriteString("Sec-WebSocket-Accept:")
		buff.WriteString(t + "\r\n\r\n")
		client.Write(buff.Bytes())
		p(key)
		recv_data(client)
	}
}

func parseKey(client net.Conn) (code int, k string) {
	bufReader := bufio.NewReader(client)
	request, err := http.ReadRequest(bufReader)
	if err != nil {
		p(err)
	}
	if request.Header.Get("Upgrade") != "websocket" {
		return http.StatusBadRequest, ""
	} else {
		key := request.Header.Get("Sec-Websocket-Key")
		return http.StatusSwitchingProtocols, key
	}
}

func reject(client net.Conn) {
	reject := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nIncorrect request"
	client.Write([]byte(reject))
	//client.Close();
}

//Funnet på nett
/*
Første byte inneholder typ beskrivelse
Andre byte inneholder lengden på dataen fra(/til) klienten
either two or eight bytes if the length does not fit in the second byte (the second byte is then a code saying how many bytes are used for the length)
the actual (raw) data
 */
