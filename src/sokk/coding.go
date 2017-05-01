package sokk

import (
    "encoding/binary"
)

type SokkMsg struct {
    fin     bool
    opCode  int
    plLen   uint64
    payload []byte
}

func encode(message *SokkMsg) (result []byte) {
    var idxData int
    length := byte(len(message.payload))
    if len(message.payload) <= 125 { //one byte to store data length
	result = make([]byte, len(message.payload)+2)
	result[1] = length
	idxData = 2
    } else if len(message.payload) >= 126 && len(message.payload) <= 65535 { //two bytes to store data length
	result = make([]byte, len(message.payload)+4)
	result[1] = 126                      //extra storage needed
	result[2] = byte(len(message.payload) >> 8) //& 255
	result[3] = (length)                 //& 255
	idxData = 4
    } else {
	result = make([]byte, len(message.payload)+10)
	result[1] = 127
	result[2] = byte(len(message.payload) >> 56) //& 255
	result[3] = byte(len(message.payload) >> 48) //& 255
	result[4] = byte(len(message.payload) >> 40) //& 255
	result[5] = byte(len(message.payload) >> 32) //& 255
	result[6] = byte(len(message.payload) >> 24) //& 255
	result[7] = byte(len(message.payload) >> 16) //& 255
	result[8] = byte(len(message.payload) >> 8)  //& 255
	result[9] = byte(len(message.payload)) & 255
	idxData = 10
    }
    result[0] = 129              //only text is supported
    for i := range message.payload { // put raw data at the correct index
	result[idxData+i] = message.payload[i]
    }

    return
}

/*
 Decode algorithm:
 LENGTH: read bits >=9 # <=15 if val <= 125 == length
 if val >= 126 read next 2 bytes (16 bits) as uint
 if val >= 127 read next 8 bytes (64 bits) as uint (MSB must be 0)
*/

func decode(rawBytes []byte) (result *SokkMsg) {
    var idxMask int
    result = &SokkMsg{
	fin:false,
	opCode:int(0x7F & rawBytes[1]),

    }
    //var mask = ((rawBytes[1] & 0x80) != 0)
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
    result.payload = make([]byte,plLen)
    result.payload = rawBytes[idxMask+4 : int(plLen)+idxMask+4]
    for i := range result.payload {
	result.payload[i] ^= masks[i%4]
    }
    result.plLen = plLen
    return
}
