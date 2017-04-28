package horse

import (
	"bufio"
	"net"
)

type client struct {
	writer *bufio.Writer
	reader *bufio.Reader
	conn net.Conn
}

func netClient(c net.Conn) *client{
	w := bufio.NewWriter(c)
	r := bufio.NewReader(c)
	new_c := &client{
		conn:c,
		writer:w,
		reader:r,
	}
	return new_c
}