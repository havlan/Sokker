package test

import(
	"testing"
	"net"
	ws "Sokker/src"
)

func TestCloseAllConnections(t *testing.T){
	_,err := net.Dial("tcp","localhost:3000")
	
	if(err == nil){
		t.Fatalf("Expected %s but got %s", nil,err)
	}
}
func TestNewStructClientsLen(t *testing.T){
	s := ws.NewSokk()
	exp := 0
	if len(s.Clients) != exp{
		t.Fatalf("Expected %s but got %s", exp,len(s.Clients))
	}
}
func TestNewSokkStruct(t *testing.T) {
	s := ws.NewSokk()
	if s == nil {
		t.Fatalf("Expected %s but got %s", s,nil)
	}
}