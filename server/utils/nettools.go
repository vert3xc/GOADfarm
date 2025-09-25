package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

func RecvUntil(conn net.Conn, delim string) (string, error) {
	reader := bufio.NewReader(conn)
	var buf bytes.Buffer
	delimBytes := []byte(delim)
	matchIdx := 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return buf.String(), err
		}

		buf.WriteByte(b)
		if b == delimBytes[matchIdx] {
			matchIdx++
			if matchIdx == len(delimBytes) {
				return buf.String(), nil
			}
		} else {
			if b == delimBytes[0] {
				matchIdx = 1
			} else {
				matchIdx = 0
			}
		}
	}
}

func SendLine(conn net.Conn, line string) error {
	_, err := fmt.Fprintf(conn, "%s\n", line)
	return err
}