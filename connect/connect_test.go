package connect_test

import (
	// "bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	success = "success"
)

func Connection() *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp","localhost:4000")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp",nil,tcpAddr)
	if err != nil {
		panic(err)
	}
	// Set a deadline for the operation (optional, for timeout)
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func LoginAsAdmin(conn *net.TCPConn) error {
	_, err := conn.Write([]byte("login duck duck"))
	if err != nil {
		return err
	}
	res := Read(conn)
	if res != success {
		return fmt.Errorf("unauth: %s", res)
	}
	return nil
}

func Read(conn *net.TCPConn) string {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(buffer[:n])," \n\t")
} 

func ConnectDb(dbname string, conn *net.TCPConn) error {
	_, err := conn.Write([]byte("connect " + dbname))
	if err != nil {
		return err
	}
	res := Read(conn)
	if res != success {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func TestConnectBasic(t *testing.T) {
	conn := Connection()
	err := LoginAsAdmin(conn);
	assert.Nil(t, err)
	err = ConnectDb("mydb", conn)
	assert.Nil(t, err)
}

func TestConnectFial(t *testing.T) {
	conn := Connection()
	err := LoginAsAdmin(conn); 
	assert.Nil(t, err)
	err = ConnectDb("doesn't_exist", conn);
	assert.NotNil(t, err)
}