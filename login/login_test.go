package login_test

import (
	"testing"

	"github.com/rag-nar1/tcp-duckdb/utils"
	"github.com/stretchr/testify/assert"
)

// TestLoginAsDuck tests the login as admin functionality
func TestLoginAsDuck(t *testing.T) {
	conn := utils.StartUp()         // Start up the connection
	err := utils.LoginAsAdmin(conn) // Attempt to login as admin
	assert.Nil(t, err)              // Assert that there is no error
}

// TestLogin tests the login functionality with valid credentials
func TestLogin(t *testing.T) {
	t.Cleanup(utils.CleanUp)        // Ensure cleanup after test
	conn := utils.StartUp()         // Start up the connection
	err := utils.LoginAsAdmin(conn) // Attempt to login as admin
	assert.Nil(t, err)              // Assert that there is no error

	username := "ragnar"
	password := "ragnar"
	err = utils.CreateUser(conn, username, password) // Create a new user
	assert.Nil(t, err)                               // Assert that there is no error

	conn2 := utils.Connection()                  // Create a new connection
	err = utils.Login(conn2, username, password) // Attempt to login with the new user
	assert.Nil(t, err)                           // Assert that there is no error
}

// TestLoginWithInvalidCredentials tests the login functionality with invalid credentials
func TestLoginWithInvalidCredentials(t *testing.T) {
	conn := utils.StartUp()         // Start up the connection
	err := utils.LoginAsAdmin(conn) // Attempt to login as admin
	assert.Nil(t, err)              // Assert that there is no error

	username := "invalid_user"
	password := "invalid_pass"
	conn2 := utils.Connection()                  // Create a new connection
	err = utils.Login(conn2, username, password) // Attempt to login with invalid credentials
	assert.NotNil(t, err)                        // Assert that there is an error
}

// TestLoginWithEmptyCredentials tests the login functionality with empty credentials
func TestLoginWithEmptyCredentials(t *testing.T) {
	conn := utils.StartUp()         // Start up the connection
	err := utils.LoginAsAdmin(conn) // Attempt to login as admin
	assert.Nil(t, err)              // Assert that there is no error

	username := ""
	password := ""
	conn2 := utils.Connection()                  // Create a new connection
	err = utils.Login(conn2, username, password) // Attempt to login with empty credentials
	assert.NotNil(t, err)                        // Assert that there is an error
}
