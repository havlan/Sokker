# Sokker
![Alt text](/logo/sokker2.jpg?raw=true "Sokker logo")


- [X] Client-Server communication
- [X] Client1-Server-Client2 communicaiton
- [ ] Sequence Diagram
- [ ] GIF to show example
- [X] ExampleCode1
- [ ] ExampleCode2
- [ ] Continuation Frame (no pri)
- [ ] Ping
- [ ] Pong
- [X] OP-Codes
- [X] WS-Frame
- [ ] Accept different data

**FIRST GO PROJECT! BEWARE!**

```golang

```

### Example usage
**No out of the box functionality until we understand golang's interfaces**
But it's still more than usable!

Example #1 (full code under examples)
```golang
func main() {
    sokk := NewSokk()
    go sokk.Start("127.0.0.1", "3001") // localhost:3000
    http.Handle("/", http.FileServer(http.Dir("../static")))
    http.ListenAndServe("localhost:3000", nil)
}
```
That was the entire main method!

### Responses

### Architecture

Why go?
- Go uses GoRoutines which is a lightweight thread of execution. Less overhead, when blocking, the runtime moves other coroutines on the same operating system thread to a different.  
- Go uses channels. Channels are the pipes that connect concurrent goroutines. Send into one and extract in the other.
- To learn something new!  
- **Thoughts on Go throughout the project?** Go can be used for many things, i'm not experienced in go, but i feel like i can do exactly everything i can do in go in another language. But go with channels can be amazingly good, it can make parallel programming "easy". And by that i mean that golang's thread model (goroutines) is quite easy and nice to use, both the normal way and the lambda way. But i still feel like i'm not in control, that is most likely since this is our first go project.

So if I want to use this, what do i need to implement?
Lets take the chat example:
- You need a Sokk struct
- You need to start it to listen
- You need to decode the byte array / data received to a websocket frame
- You need to encode the websocket frame to a byte array
- You need to send the data to all your clients (ws.clients)



### Inspiration and external links
- http://stackoverflow.com/questions/18368130/how-to-parse-and-validate-a-websocket-frame-in-java
- https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers
- http://stackoverflow.com/questions/11815894/how-to-read-write-arbitrary-bits-in-c-c
- 

So the whole thought about this library is that a goroutine listens to a connection. Which receives the data and sends it onwards. Another gorutine writes to the client(s). 
