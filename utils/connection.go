package utils

import (
	"bufio"
)


func Write(writer *bufio.Writer, data []byte) error {
	data = append(data, '\n')
	if _, err := writer.Write(data); err != nil {
        return err
    }
    return writer.Flush()
}
