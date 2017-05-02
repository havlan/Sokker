package sokk

import (
	"encoding/binary"
)


func Encode(message *SokkMsg) (result []byte) {
	var idxData int
	length := byte(len(message.Payload))
	if len(message.Payload) <= 125 { //one byte to store data length
		result = make([]byte, len(message.Payload)+2)
		result[1] = length
		idxData = 2
	} else if len(message.Payload) >= 126 && len(message.Payload) <= 65535 { //two bytes to store data length
		result = make([]byte, len(message.Payload)+4)
		result[1] = 126                             //extra storage needed
		result[2] = byte(len(message.Payload) >> 8) //& 255
		result[3] = (length)                        //& 255
		idxData = 4
	} else {
		result = make([]byte, len(message.Payload)+10)
		result[1] = 127
		result[2] = byte(len(message.Payload) >> 56)
		result[3] = byte(len(message.Payload) >> 48)
		result[4] = byte(len(message.Payload) >> 40)
		result[5] = byte(len(message.Payload) >> 32)
		result[6] = byte(len(message.Payload) >> 24)
		result[7] = byte(len(message.Payload) >> 16)
		result[8] = byte(len(message.Payload) >> 8)
		result[9] = byte(len(message.Payload)) & 255
		idxData = 10
	}
	result[0] = 129              //only text is supported
	for i := range message.Payload { // put raw data at the correct index
		result[idxData+i] = message.Payload[i]
	}
	
	return
}

/*
	https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers
	 Decode algorithm:
	 LENGTH: read bits >=9 # <=15 if val <= 125 == length
	 if val >= 126 read next 2 bytes (16 bits) as uint
	 if val >= 127 read next 8 bytes (64 bits) as uint (MSB must be 0)
*/

func Decode(rawBytes []byte) (result *SokkMsg) {
	var idxMask int
	result = &SokkMsg{
		Fin:(rawBytes[0] & 0x80) != 0,
		OpCode: int(0x7F & rawBytes[0]),
	}
	var plLen = uint64(0x7F & rawBytes[1])
	if plLen == 126 {
		idxMask = 4
		plLen = uint64(binary.LittleEndian.Uint16(rawBytes[2:4])) // short
	} else if plLen == 127 {
		idxMask = 10
		plLen = binary.LittleEndian.Uint64(rawBytes[2:10]) // long
	} else {
		idxMask = 2
	}
	masks := rawBytes[idxMask : idxMask+4]
	result.Payload = make([]byte,plLen)
	result.Payload = rawBytes[idxMask+4 : int(plLen)+idxMask+4]
	for i := range result.Payload {
		result.Payload[i] ^= masks[i%4]
	}
	result.PlLen = plLen
	return
}

