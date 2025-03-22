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
	err = utils.CreateDB(conn, "db1")
	assert.Nil(t, err) 
	err = utils.ConnectDb(conn, "db1")
	assert.Nil(t, err)
}

func TestConnectFial(t *testing.T) {
	conn := utils.Connection()
	err := utils.LoginAsAdmin(conn); 
	assert.Nil(t, err)
	err = utils.ConnectDb(conn, "doesn't_exist");
	assert.NotNil(t, err)
}