package create

import (
	global "TCP-Duckdb/utils"
	"bufio"
)

func Write(w *bufio.Writer, msg []byte) error {
	return  global.Write(w, msg)
}

func Error(w *bufio.Writer, msg []byte) error {
	msg = append([]byte("Error:"), msg...)
	return global.Write(w, msg)
}