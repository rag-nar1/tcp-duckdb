package response

import (
	"bufio"
)

func Success(w *bufio.Writer) error {
	if _, err := w.Write([]byte(SuccessMsg)); err != nil {
        return err
    }
    return w.Flush()
}

func WriteData(w *bufio.Writer, data []byte) error {
	if _, err := w.Write(data); err != nil {
        return err
    }
    return w.Flush()
}