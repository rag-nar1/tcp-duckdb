package grant_test

import (
	"TCP-Duckdb/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBasicGrantOverDb(t *testing.T) {
	dbname := "db1"
	username := "ragnar"
	password := "ragnar"
	defer utils.CleanUp(dbname)

	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during creating user

	err = utils.Grant(conn, username, dbname, "read")
	assert.Nil(t, err) // Ensure no error during granting user access
	
	conn2 := utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err) // Ensure no error during login

	err = utils.ConnectDb(conn2, dbname)
	assert.Nil(t, err) // Ensure no error during connection
}

func TestGrantInvalidPermission(t *testing.T) {
	dbname := "db3"
	username := "doe"
	password := "doe"
	defer utils.CleanUp(dbname)

	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during creating user

	err = utils.Grant(conn, username, dbname, "invalid_permission")
	assert.NotNil(t, err) // Ensure error during granting invalid permission
}

func TestGrantWithoutDb(t *testing.T) {
	username := "alice"
	password := "alice"
	defer utils.CleanUp("")
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during creating user

	err = utils.Grant(conn, username, "non_existent_db", "read")
	assert.NotNil(t, err) // Ensure error during granting access to non-existent db
}

func TestGrantMultiplePermissions(t *testing.T) {
	dbname := "db4"
	username := "bob"
	password := "bob"
	defer utils.CleanUp(dbname)
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during creating user

	err = utils.Grant(conn, username, dbname, "read")
	assert.Nil(t, err) // Ensure no error during granting read access

	err = utils.Grant(conn, username, dbname, "write")
	assert.Nil(t, err) // Ensure no error during granting write access

	conn2 := utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err) // Ensure no error during login

	err = utils.ConnectDb(conn2, dbname)
	assert.Nil(t, err) // Ensure no error during connection
}

func TestGrantOverTable(t *testing.T) {
	dbname := "db4"
	username := "bob"
	password := "bob"
	tablename := "t1"
	defer utils.CleanUp(dbname)
	conn := utils.StartUp()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during creating user

	err = utils.Grant(conn, username, dbname, "read")
	assert.Nil(t, err) // Ensure no error during granting user access

	err = utils.ConnectDb(conn, dbname)
	assert.Nil(t, err) // Ensure no error during connection

	err = utils.CreateTable(conn, tablename)
	assert.Nil(t, err) // Ensure no error during Creating table

	conn2 := utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err) // Ensure no error during login

	// todo: check the implementation of the grant command

}
