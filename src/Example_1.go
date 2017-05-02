package main

import (
	ws "Sokker/src/sokk"
	"fmt"
	"net"
	"os"
)

func main() {
	sokk := ws.NewSokk()
	
	//onClose client is already removed from the list.
	sokk.OnClose = func(){
		fmt.Println("OnClose!")
	}
	sokk.OnConnection = func(c net.Conn){
		fmt.Println("NEW CONNECTIONS!")
		sokk.Clients = append(sokk.Clients,c) // add the user into the accepted client list
	}
	sokk.OnError = func(w string, e error){ // custom handle error
		fmt.Println(w, " ", e.Error())
		os.Exit(1)
		
	}
	sokk.OnMessage = func(b ws.SokkMsg){
		fmt.Println(string(b.Payload[:b.PlLen])) // prints the data
		sokk.Send(&b) // sends to all Clients which exists in the sockets array of connections
		
	}
	sokk.Start("127.0.0.1", "3001") // localhost:3000
	//http.Handle("/", http.FileServer(http.Dir("../static")))
	//http.ListenAndServe("localhost:3000", nil)
}