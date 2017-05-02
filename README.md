# Sokker
![Alt text](/logo/sokker2.jpg?raw=true "Sokker logo")



**FIRST GO PROJECT! BEWARE!**
### Prerequisites
- Git
- Go 1.8 and working go environment
  - **GOROOT** to the installation directory ex. /usr/local/go
  - **Path** to the bin subdirectory of GOROOT $GOROOT/bin (/usr/local/go/bin)
  - **GOPATH** to a workspace. A workspace consists of a root directory and a src subdirectory. Ex. /home/myname/GoglandProjects
  - So it looks like this:
  - Goglandprojects (where my gopath is set to)
     - src 
       - github.com
         - havlan
	 	- Sokker
			- src
  
### Installation

**Check your environment with "go env"**
```
go get github.com/havlan/Sokker
```
```
If you are in your gopath/src directory
cd github.com && cd havlan && cd Sokker && cd Examples && go run Example_http.go 
```

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

### Example usage
Example HTTP/WS (full code under examples)
```golang

```
That was the entire main method!

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
