package login_test

import (
	"TCP-Duckdb/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoginAsDuck(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err)
}

func TestLogin(t *testing.T) {
	defer utils.CleanUp("")
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err)
	username := "ragnar"
	password := "ragnar"
	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err)
	conn2 := utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err)
}
func TestLoginWithInvalidCredentials(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err)
	username := "invalid_user"
	password := "invalid_pass"
	conn2 := utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.NotNil(t, err)
}

func TestLoginWithEmptyCredentials(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err)
	username := ""
	password := ""
	conn2 := utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.NotNil(t, err)
}

