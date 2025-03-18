package create_test

import (
	"TCP-Duckdb/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestCreateDb tests the creation of a new database.
func TestCreateDb(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup
	err = utils.CreateDB(conn, "db1")
	assert.Nil(t, err) // Ensure no error during database creation
	err = utils.CleanUpDb("db1")
	assert.Nil(t, err) // Ensure no error during cleanup
}

// TestAlreadyExists tests the behavior when trying to create a database that already exists.
func TestCreateDbAlreadyExists(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup
	err = utils.CreateDB(conn, "db1")
	assert.Nil(t, err) // Ensure no error during first database creation
	err = utils.CreateDB(conn, "db1")
	assert.NotNil(t, err) // Ensure an error is returned when creating a database that already exists
	err = utils.CleanUpDb("db1")
	assert.Nil(t, err) // Ensure no error during cleanup
}

// TestCreateDbAndConnect tests the creation of a new database and connecting to it.
func TestCreateDbAndConnect(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup
	err = utils.CreateDB(conn, "db1")
	assert.Nil(t, err) // Ensure no error during database creation
	err = utils.ConnectDb("db1", conn)
	assert.Nil(t, err) // Ensure no error during database connection
	err = utils.CleanUpDb("db1")
	assert.Nil(t, err) // Ensure no error during cleanup
}
// TestCreateUser tests the creation of a new user.
func TestCreateUser(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup
	err = utils.CreateUser(conn, "ragnar", "ragnar")
	assert.Nil(t, err) // Ensure no error during user creation
	utils.CleanUpUsers() // Clean up users after test
}

// TestCreateUserAlreadyExists tests the behavior when trying to create a user that already exists.
func TestCreateUserAlreadyExists(t *testing.T) {
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup
	err = utils.CreateUser(conn, "ragnar", "ragnar")
	assert.Nil(t, err) // Ensure no error during first user creation
	err = utils.CreateUser(conn, "ragnar", "ragnar")
	assert.NotNil(t, err) // Ensure an error is returned when creating a user that already exists
	utils.CleanUpUsers() // Clean up users after test
}