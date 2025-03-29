package grant_test

import (
	"testing"

	"github.com/rag-nar1/tcp-duckdb/utils"
	"github.com/stretchr/testify/assert"
)

func TestBasicGrantOverDb(t *testing.T) {
	dbname := "db1"
	username := "ragnar"
	password := "ragnar"
	t.Cleanup(utils.CleanUp) // Clean up resources after test

	conn := utils.StartUp() // Start up the database connection
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during admin login

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during user creation

	err = utils.GrantDb(conn, username, dbname, "read")
	assert.Nil(t, err) // Ensure no error during granting read access to the user

	conn2 := utils.Connection() // Create a new connection
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err) // Ensure no error during user login

	err = utils.ConnectDb(conn2, dbname)
	assert.Nil(t, err) // Ensure no error during database connection
}

func TestGrantInvalidPermission(t *testing.T) {
	dbname := "db3"
	username := "doe"
	password := "doe"
	t.Cleanup(utils.CleanUp) // Clean up resources after test

	conn := utils.StartUp() // Start up the database connection
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during admin login

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during user creation

	err = utils.GrantDb(conn, username, dbname, "invalid_permission")
	assert.NotNil(t, err) // Ensure error during granting invalid permission
}

func TestGrantWithoutDb(t *testing.T) {
	username := "alice"
	password := "alice"
	t.Cleanup(utils.CleanUp) // Clean up resources after test

	conn := utils.StartUp() // Start up the database connection
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during admin login

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during user creation

	err = utils.GrantDb(conn, username, "non_existent_db", "read")
	assert.NotNil(t, err) // Ensure error during granting access to non-existent database
}

func TestGrantMultiplePermissions(t *testing.T) {
	dbname := "db4"
	username := "bob"
	password := "bob"
	t.Cleanup(utils.CleanUp) // Clean up resources after test

	conn := utils.StartUp() // Start up the database connection
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during admin login

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during user creation

	err = utils.GrantDb(conn, username, dbname, "read")
	assert.Nil(t, err) // Ensure no error during granting read access

	err = utils.GrantDb(conn, username, dbname, "write")
	assert.Nil(t, err) // Ensure no error during granting write access

	conn2 := utils.Connection() // Create a new connection
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err) // Ensure no error during user login

	err = utils.ConnectDb(conn2, dbname)
	assert.Nil(t, err) // Ensure no error during database connection
}

func TestGrantOverTable(t *testing.T) {
	dbname := "db4"
	username := "bob"
	password := "bob"
	tablename := "t1"
	t.Cleanup(utils.CleanUp) // Clean up resources after test

	conn := utils.StartUp() // Start up the database connection
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during admin login

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err) // Ensure no error during user creation

	err = utils.GrantDb(conn, username, dbname, "read")
	assert.Nil(t, err) // Ensure no error during granting database read access

	err = utils.ConnectDb(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database connection

	err = utils.CreateTable(conn, tablename)
	assert.Nil(t, err) // Ensure no error during table creation

	conn2 := utils.Connection() // Create a new connection
	err = utils.Login(conn2, username, password)
	assert.Nil(t, err) // Ensure no error during user login

	err = utils.ConnectDb(conn2, dbname)
	assert.Nil(t, err) // Ensure no error during database connection

	err = utils.Query(conn2, "select * from t1;")
	assert.NotNil(t, err) // Ensure error during query without table access

	err = utils.Query(conn2, "insert into t1(id, name) values(1, 'ragnar');")
	assert.NotNil(t, err) // Ensure error during insert without table access

	conn = utils.Connection() // Create a new connection
	err = utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during admin login

	err = utils.GrantTable(conn, username, dbname, tablename, "select")
	assert.Nil(t, err) // Ensure no error during granting table select access

	err = utils.Query(conn2, "select * from t1;")
	assert.Nil(t, err) // Ensure no error during query after granting select access

	err = utils.Query(conn2, "insert into t1(id, name) values(1, 'ragnar');")
	assert.NotNil(t, err) // Ensure error during insert without insert access

	err = utils.GrantTable(conn, username, dbname, tablename, "insert")
	assert.NotNil(t, err) // Ensure error during granting insert access without write access

	err = utils.GrantDb(conn, username, dbname, "write")
	assert.Nil(t, err) // Ensure no error during granting database write access

	err = utils.GrantTable(conn, username, dbname, tablename, "insert")
	assert.Nil(t, err) // Ensure no error during granting table insert access

	err = utils.Query(conn2, "insert into t1(id, name) values(1, 'ragnar');")
	assert.Nil(t, err) // Ensure no error during insert after granting insert access
}
