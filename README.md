# Sokker
![Alt text](/logo/sokker2.jpg?raw=true "Sokker logo")



**FIRST GO PROJECT! BEWARE!**
### Prerequisites
- Git
- Go 1.8 and working go environment
  - **GOROOT** to the installation directory ex. /usr/local/go
  - **Path** to the bin subdirectory of GOROOT $GOROOT/bin (/usr/local/go/bin)
  - **GOPATH** to a workspace. A workspace consists of a root directory and a src subdirectory. Ex. /home/myname/GoglandProjects
  
   *Make sure your GOPATH looks something like this* 
     ```
      GOPATH
      └───src
      	└───websocket
      │   │   	encoding.go
      │   │   	handshake.go
      │   │		websocketserver.go
      │
      │
    ```
  
### Installation

**Check your environment with "go env"**
```
go get github.com/havlan/Sokker
go test github.com/havlan/Sokker/test -v
```
If you are in your gopath/src directory
```
cd github.com/havlan/Sokker/examples && go run Example_http.go 
```
Now localhost:3000 hopefully shows a chat.




### Examples
The file Example_http.go

```golang
package main

import (
	"fmt"
	"net"
	ws "github.com/havlan/Sokker/src" // alias and then import/path/to/correct/package. 
	"net/http"      // you can skip the alias bit. Then the way to use it is sokk.MethodName
)


func main() {
	sokk := ws.NewSokk()
	
	//onClose client is already removed from the list.
	sokk.OnClose = func(c net.Conn){
		fmt.Println("OnClose!")
	}
	sokk.OnConnection = func(c net.Conn){
		fmt.Println("NEW CONNECTIONS!")
		sokk.Clients = append(sokk.Clients,c) // add the user into the accepted client list
	}
	sokk.OnError = func(w string, e error){ // custom handle error
		fmt.Println(w, " ", e.Error())
		panic(e)
		
	}
	sokk.OnMessage = func(b ws.SokkMsg){
		fmt.Println(string(b.Payload[:b.PlLen])) // prints the data
		sokk.Send(&b) // sends to all Clients which exists in the sockets array of connections
		
	}
	//handle http on main thread, socket gets new goroutine
	go sokk.Start("127.0.0.1", "3001") // localhost:3001
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe("localhost:3000", nil)
}

```
All the other examples and the basic usage of this library is based on the functions OnClose, OnConnection, OnError and OnMessage.
Add a logger to log both messages and errors (Example_log.go)

```golang
f,err := os.Create("logfile.txt")
	if err != nil {
		sokk.OnError("Error creating logfile for errors(ironic)",err)
	}
	defer f.Close()
	var errBuff =  bufio.NewWriter(f)
//the onError method
sokk.OnError = func(w string, e error){ // custom handle error
		fmt.Println(w, " ", e.Error())
		errBuff.WriteString(time.Now().String())
		errBuff.WriteString(w)
		errBuff.WriteString(e.Error())
		f.Sync()
		errBuff.Flush()
		os.Exit(1)
	}
//setup message logger
msgFile, errMs := os.Create("messages.txt")
if errMs != nil{
	sokk.OnError("Error creating logfile for messages", errMs)
}
defer msgFile.Close()
sokk.OnMessage = func(b ws.SokkMsg){
		fmt.Println(string(b.Payload[:b.PlLen]))// prints the data
		msgBuff.WriteString(time.Now().String() + " ")
		msgBuff.Write(b.Payload[:b.PlLen])
		msgBuff.WriteString("\n")
		msgFile.Sync() // sync file writing
		msgBuff.Flush()
		sokk.Send(&b) // sends to all Clients which exists in the sockets array of connections
	}

```


### Architecture

Why go?
- Go uses GoRoutines which is a lightweight thread of execution. Less overhead, when blocking, the runtime moves other coroutines on the same operating system thread to a different.  
- Go uses channels. Channels are the pipes that connect concurrent goroutines. Send into one and extract in the other.
- To learn something new!  
- **Thoughts on Go throughout the project?** Go can be used for many things, i'm not experienced in go, but i feel like i can do exactly everything i can do in go in another language. But go with channels can be amazingly good, it can make parallel programming "easy". And by that i mean that golang's thread model (goroutines) is quite easy and nice to use, both the normal way and the lambda way. But i still feel like i'm not in control, that is most likely since this is our first go project.



### Inspiration and external links
- http://stackoverflow.com/questions/18368130/how-to-parse-and-validate-a-websocket-frame-in-java
- https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers
- http://stackoverflow.com/questions/11815894/how-to-read-write-arbitrary-bits-in-c-c

So the whole thought about this library is that a goroutine listens to a connection. Which receives the data and sends it onwards. Another gorutine writes to the client(s). 
