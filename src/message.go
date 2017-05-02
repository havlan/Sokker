package main


type SokkMsg struct {
	fin     bool
	opCode  int
	plLen   uint64
	payload []byte
}