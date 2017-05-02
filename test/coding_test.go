package test

import (
	"testing"
	ws "Sokker/src"
)


func TestMagicStr(t *testing.T){
	act := ws.Magic_str("dGhlIHNhbXBsZSBub25jZQ==" + ws.Magic_server_key)
	exp := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="
	if(act != exp){
		t.Fatalf("Expected %s but got %s", exp, act)
	}
}

func TestDecodeLength_1(t *testing.T){
	pl := []byte{129,131,53,218,212,244,84,169,176} // asd
	st := ws.Decode(pl)
	exp := 3
	act := st.PlLen
	
	if(exp != int(act)){
		t.Fatalf("Expected %s but got %s", exp, act)
	}
}
func TestDecodeLength_2(t *testing.T){
	pl := []byte{129,140,71,247,101,246,47,146,13,147,47,158, 13, 159, 47, 152, 13, 153} // hehehihihoho
	st := ws.Decode(pl)
	exp := 12
	if exp != int(st.PlLen){
		t.Fatalf("Expected %s but got %s", exp, int(st.PlLen))
	}
}

func TestDecodeByteContain(t *testing.T)  {
	pl := []byte{129,131,53,218,212,244,84,169,176} // asd
	st := ws.Decode(pl)
	exp := []byte{97,115,100}
	
	for i := range exp {
		if exp[i] != st.Payload[i] {
			t.Fatalf("Expected %s but got %s", exp[i], st.Payload[i])
		}
	}
}

func TestDecodeFinContain(t *testing.T){
	pl := []byte{129,131,53,218,212,244,84,169,176} // asd
	st := ws.Decode(pl)
	exp := true
	if st.Fin != exp {
		t.Fatalf("Expected %s but got %s", exp,st.Fin)
	}
}
func TestDecodeOpCode(t *testing.T){
	pl := []byte{129,131,53,218,212,244,84,169,176} // asd
	st := ws.Decode(pl)
	exp := 1
	if st.OpCode != exp{
		t.Fatalf("Expected %s but got %s", exp,st.OpCode)
	}
}
func TestEncodeLength(t *testing.T){
	pl := []byte{129,131,53,218,212,244,84,169,176} // asd
	st := ws.Decode(pl)
	act := ws.Encode(st)[2:] // last three bytes represents data
	exp := []byte{97,115,100} // asd
	for i := range act{
		if(act[i] != exp[i]){
			t.Fatalf("Expected %s but got %s", exp,act)
		}
	}
}
