package connect_test

import (
	"TCP-Duckdb/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestConnectBasic(t *testing.T) {
	conn := utils.Connection()
	err := utils.LoginAsAdmin(conn);
	assert.Nil(t, err)
	err = utils.ConnectDb("mydb", conn)
	assert.Nil(t, err)
}

func TestConnectFial(t *testing.T) {
	conn := utils.Connection()
	err := utils.LoginAsAdmin(conn); 
	assert.Nil(t, err)
	err = utils.ConnectDb("doesn't_exist", conn);
	assert.NotNil(t, err)
}