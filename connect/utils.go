package connect

import (
	global "TCP-Duckdb/utils"
	"bufio"
)

func Write(w *bufio.Writer, msg []byte) error {
	return  global.Write(w, msg)
}