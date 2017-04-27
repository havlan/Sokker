package main

import (
	"net"
	"bufio"
	//"bytes"
	"log"
)
const (
	small_buff = 128
	normal_buff = 256
	big_buff = 512
)

type client_ struct {
	inc chan string
	out chan string
	conn net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

//create a new client
//listens to client
func new_client_ (c net.Conn)*client_{
	w := bufio.NewWriter(c)
	r := bufio.NewReader(c)
	n_cli:= &client_{
		inc: make(chan string),
		out: make(chan string),
		reader: r,
		writer:w,
	}
	n_cli.listen()
	return n_cli

}

func (c *client_) read (){ // \r\n\r\n
	for{
		bt, _ := c.reader.ReadString('\n')
		c.inc <- bt //channel that sends to
	}
}
//write to client
func (c *client_) write (){
	for data := range c.out{
		c.writer.WriteString(data)
		c.writer.Flush()
	}
}

//2 goroutines to which listens to users
func (c *client_) listen(){
	log.Println("Listen: ",c.conn)
	go c.read()
	go c.write()
}