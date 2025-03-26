package connect_test

import (
	"testing"

	"github.com/rag-nar1/TCP-Duckdb/utils"

	"github.com/stretchr/testify/assert"
)

// TestConnectBasic tests the basic connection flow
func TestConnectBasic(t *testing.T) {
	t.Cleanup(utils.CleanUp)
	conn := utils.Connection()         // Establish a new connection
	err := utils.LoginAsAdmin(conn)    // Login as admin
	assert.Nil(t, err)                 // Ensure no error occurred during login
	err = utils.CreateDB(conn, "db1")  // Create a new database named "db1"
	assert.Nil(t, err)                 // Ensure no error occurred during database creation
	err = utils.ConnectDb(conn, "db1") // Connect to the newly created database
	assert.Nil(t, err)                 // Ensure no error occurred during database connection
}

// TestConnectFail tests the connection to a non-existent database
func TestConnectFail(t *testing.T) {
	t.Cleanup(utils.CleanUp)
	conn := utils.Connection()                   // Establish a new connection
	err := utils.LoginAsAdmin(conn)              // Login as admin
	assert.Nil(t, err)                           // Ensure no error occurred during login
	err = utils.ConnectDb(conn, "doesn't_exist") // Attempt to connect to a non-existent database
	assert.NotNil(t, err)                        // Ensure an error occurred during database connection
}
