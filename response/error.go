package response

import (
	"bufio"
	"fmt"
	"strings"
)

func Error(w *bufio.Writer, msg []byte) error {
	msg = append([]byte("Error: "), msg...)
	if _, err := w.Write(msg); err != nil {
        return err
    }
    return w.Flush()
}

func BadRequest(w *bufio.Writer) error {
    return Error(w, []byte(BadRequestMsg))
}

func InternalError(w *bufio.Writer) error {
	return Error(w, []byte(InternalErrorMSG))
}

func UnauthorizedError(w *bufio.Writer) error {
	return Error(w, []byte(UnauthorizedMSG))
}

func DoesNotExist(w *bufio.Writer, prefix string, objNames ...string) error {
	objName := strings.Join(objNames, ",")
	return Error(w, []byte(prefix + " " + objName + " " + DoesNotExistMsg))
}

func DoesNotExistDatabse(w *bufio.Writer, dbname string) error {
	return DoesNotExist(w, "Database", dbname)
}

func DoesNotExistUser(w *bufio.Writer, username string) error {
	return DoesNotExist(w, "User", username)
}

func DoesNotExistTables(w *bufio.Writer, tables ...string) error {
	return DoesNotExist(w, "Tables", tables...)
}


func AccessDenied(w *bufio.Writer, prtObj, prtVal, childObj, childVal string) error {
	return Error(w, []byte(fmt.Sprintf(AccessDeniedMsg, prtObj, prtVal, childObj, childVal)))
}

func AccesDeniedOverDatabase(w *bufio.Writer, username, dbname string) error {
	return AccessDenied(w, "User", username, "Database", dbname)
}

func AccesDeniedOverTables(w *bufio.Writer, username string, tables []byte) error {
	return AccessDenied(w, "User", username, "Tables", string(tables))
}