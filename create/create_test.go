package create_test

import (
	"testing"

	"github.com/rag-nar1/TCP-Duckdb/utils"
	"github.com/stretchr/testify/assert"
)

// TestCreateDb tests the creation of a new database.
func TestCreateDb(t *testing.T) {
	t.Cleanup(utils.CleanUp)          // Ensure cleanup after test
	conn := utils.StartUp()           // Start up the connection
	err := utils.LoginAsAdmin(conn)   // Log in as admin
	assert.Nil(t, err)                // Ensure no error during login
	err = utils.CreateDB(conn, "db1") // Create a new database named "db1"
	assert.Nil(t, err)                // Ensure no error during database creation
}

// TestCreateDbAlreadyExists tests the behavior when trying to create a database that already exists.
func TestCreateDbAlreadyExists(t *testing.T) {
	t.Cleanup(utils.CleanUp)          // Ensure cleanup after test
	conn := utils.StartUp()           // Start up the connection
	err := utils.LoginAsAdmin(conn)   // Log in as admin
	assert.Nil(t, err)                // Ensure no error during login
	err = utils.CreateDB(conn, "db1") // Create a new database named "db1"
	assert.Nil(t, err)                // Ensure no error during first database creation
	err = utils.CreateDB(conn, "db1") // Try to create the same database again
	assert.NotNil(t, err)             // Ensure an error is returned when creating a database that already exists
}

// TestCreateDbAndConnect tests the creation of a new database and connecting to it.
func TestCreateDbAndConnect(t *testing.T) {
	t.Cleanup(utils.CleanUp)           // Ensure cleanup after test
	conn := utils.StartUp()            // Start up the connection
	err := utils.LoginAsAdmin(conn)    // Log in as admin
	assert.Nil(t, err)                 // Ensure no error during login
	err = utils.CreateDB(conn, "db1")  // Create a new database named "db1"
	assert.Nil(t, err)                 // Ensure no error during database creation
	err = utils.ConnectDb(conn, "db1") // Connect to the newly created database
	assert.Nil(t, err)                 // Ensure no error during database connection
}

// TestCreateUser tests the creation of a new user.
func TestCreateUser(t *testing.T) {
	t.Cleanup(utils.CleanUp)                         // Ensure cleanup after test
	conn := utils.StartUp()                          // Start up the connection
	err := utils.LoginAsAdmin(conn)                  // Log in as admin
	assert.Nil(t, err)                               // Ensure no error during login
	err = utils.CreateUser(conn, "ragnar", "ragnar") // Create a new user with username and password "ragnar"
	assert.Nil(t, err)                               // Ensure no error during user creation
}

// TestCreateUserAlreadyExists tests the behavior when trying to create a user that already exists.
func TestCreateUserAlreadyExists(t *testing.T) {
	t.Cleanup(utils.CleanUp)                         // Ensure cleanup after test
	conn := utils.StartUp()                          // Start up the connection
	err := utils.LoginAsAdmin(conn)                  // Log in as admin
	assert.Nil(t, err)                               // Ensure no error during login
	err = utils.CreateUser(conn, "ragnar", "ragnar") // Create a new user with username and password "ragnar"
	assert.Nil(t, err)                               // Ensure no error during first user creation
	err = utils.CreateUser(conn, "ragnar", "ragnar") // Try to create the same user again
	assert.NotNil(t, err)                            // Ensure an error is returned when creating a user that already exists
}
