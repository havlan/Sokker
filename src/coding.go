package main

import "fmt"


type sokkMsg struct{
	fin bool
	opCode int
	plLen int
	payload []byte
}



// TODO MAKE ENCODE CREATE A sokkMSG STRUCT!
func encode(message string) (result []byte) {
	rawBytes := []byte(message) // convert to byte array
	var idxData int

	length := byte(len(rawBytes))
	if len(rawBytes) <= 125 { //one byte to store data length
		result = make([]byte, len(rawBytes)+2)
		result[1] = length
		idxData = 2
	} else if len(rawBytes) >= 126 && len(rawBytes) <= 65535 { //two bytes to store data length
		result = make([]byte, len(rawBytes)+4)
		result[1] = 126 //extra storage needed
		result[2] = byte(len(rawBytes) >> 8) //& 255
		result[3] = (length) //& 255
		idxData = 4
	} else {
		result = make([]byte, len(rawBytes)+10)
		result[1] = 127
		result[2] = byte(len(rawBytes) >> 56) //& 255
		result[3] = byte(len(rawBytes) >> 48) //& 255
		result[4] = byte(len(rawBytes) >> 40) //& 255
		result[5] = byte(len(rawBytes) >> 32) //& 255
		result[6] = byte(len(rawBytes) >> 24) //& 255
		result[7] = byte(len(rawBytes) >> 16) //& 255
		result[8] = byte(len(rawBytes) >>  8) //& 255
		result[9] = byte(len(rawBytes)) & 255
		idxData = 10
	}
	result[0] = 129 //only text is supported
	for i, b := range rawBytes { // put raw data at the correct index
		result[idxData+i] = b
	}
	
	return
}

/*
trim byte array
 */
func trimByteArr(raw []byte){
	for i := range raw {
		if raw[i] == 0 {
			fmt.Print(i)
		}
	}
}

/*
	 Decode algorithm:
	 LENGTH: read bits >=9 # <=15 if val <= 125 == length
	 if val >= 126 read next 2 bytes (16 bits) as uint
	 if val >= 127 read next 8 bytes (64 bits) as uint (MSB must be 0)
 */


func decode(rawBytes []byte) string {
	var idxMask int
	if rawBytes[1]-128 == 126 {
		idxMask = 4
	} else if rawBytes[1]-128 == 127 {
		idxMask = 10
	} else {
		idxMask = 2
	}
	var length = idxMask+4
	masks := rawBytes[idxMask : idxMask+4]
	for i := idxMask+4; i<=len(rawBytes); i++{
		if rawBytes[i]==0 {
			var control_ = true
			for k := 1; k<=10 ;k++{ //  check if the next 10 bytes are 0. TODO get a better algorithm for a checker
				if rawBytes[i+k] != 0{
					control_ = false
					break
				}
			}
			if control_{
				length = i
				break
			}
		}
	}
	data := rawBytes[idxMask+4 : length]
	decoded := make([]byte, len(rawBytes)-idxMask+4)

	for i, b := range data {
		decoded[i] = b ^ masks[i%4]
	}
	return string(decoded)
}


//checks a upcode for
//checks a upcode for
func opcode(rawBytes []byte) int{
	opcodeInt := 0;
	opcodeS := fmt.Sprintf("%08b", rawBytes[0])
	opcodeS = opcodeS[4:len(opcodeS)]
	//opcodes for 0,1,2,8,9,10
	if opcodeS == "0000"{// continue
		opcodeInt = 0
	}else if opcodeS == "0001" { // text
		opcodeInt = 1
	}else if opcodeS == "0010" { // binary
		opcodeInt = 2
	}else if opcodeS == "1000" {// Close 0x8
		opcodeInt = 8
	}else if opcodeS == "1001" {// Ping 0x9
		opcodeInt = 9
	}else if opcodeS  == "1010" { // Pong 0xA
		opcodeInt = 10
	}

	return opcodeInt
}