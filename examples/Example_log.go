package main

import (
	ws "github.com/havlan/Sokker/src"
	"fmt"
	"net"
	"os"
	"bufio"
	"time"
	"net/http"
)

func main() {
	sokk := ws.NewSokk()
	f,err := os.Create("logfile.txt")
	if err != nil {
		sokk.OnError("Error creating logfile for errors(ironic)",err)
	}
	defer f.Close()
	
	
	msgFile, errMs := os.Create("messages.txt")
	if errMs != nil{
		sokk.OnError("Error creating logfile for messages", errMs)
	}
	defer msgFile.Close()
	
	var errBuff =  bufio.NewWriter(f)
	var msgBuff = bufio.NewWriter(msgFile)
	//onClose client is already removed from the list.
	sokk.OnClose = func(c net.Conn){
		fmt.Println("OnClose!")
		c.Close()
	}
	sokk.OnConnection = func(c net.Conn){
		fmt.Println("NEW CONNECTIONS!")
		sokk.Clients = append(sokk.Clients,c) // add the user into the accepted client list
	}
	sokk.OnError = func(w string, e error){ // custom handle error
		fmt.Println(w, " ", e.Error())
		errBuff.WriteString(time.Now().String())
		errBuff.WriteString(w)
		errBuff.WriteString(e.Error())
		f.Sync()// Sync commits the current contents of the file to stable storage. Typically, this means flushing the file system's in-memory copy of recently written data to disk.
		errBuff.Flush()
		os.Exit(1)
	}
	sokk.OnMessage = func(b ws.SokkMsg){
		fmt.Println(string(b.Payload[:b.PlLen]))// prints the data
		msgBuff.WriteString(time.Now().String() + " ")
		msgBuff.Write(b.Payload[:b.PlLen])
		msgBuff.WriteString("\n")
		msgFile.Sync() // Sync commits the current contents of the file to stable storage. Typically, this means flushing the file system's in-memory copy of recently written data to disk.
		msgBuff.Flush()
		sokk.Send(&b) // sends to all Clients which exists in the sockets array of connections
	}
	//if you are going to have a websocket and httpserver running together, one of them needs to get "go"ing (goroutine)
	go sokk.Start("127.0.0.1", "3001") // localhost:3000
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe("localhost:3000", nil)
}