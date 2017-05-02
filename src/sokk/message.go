package sokk


type SokkMsg struct {
	Fin     bool
	OpCode  int
	PlLen   uint64
	Payload []byte
}