package link_test

import (
	"log"
	"testing"

	"github.com/rag-nar1/TCP-Duckdb/utils"

	"github.com/stretchr/testify/assert"
)

// TestBasicLink tests the basic linking of a database
func TestBasicLink(t *testing.T) {
	dbname := "db2"
	connStr := "postgresql://postgres:1242003@localhost:5432/testdb"
	t.Cleanup(utils.CleanUp)

	conn := utils.StartUp()
	defer conn.Close()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.Link(conn, dbname, connStr)
	assert.Nil(t, err) // Ensure no error during linking database

	err = utils.ConnectDb(conn, dbname)
	assert.Nil(t, err) // Ensure no error during connecting to the database

	err = utils.Query(conn, "select * from t1;")
	assert.Nil(t, err) // Ensure no error during querying the database
}

// TestDoubleLink tests linking the same database to two different connections
func TestDoubleLink(t *testing.T) {
	dbname := "db2"
	connStr := "postgresql://postgres:1242003@localhost:5432/testdb"
	connStr2 := "postgresql://postgres:1242003@localhost:5432/testdb2"
	t.Cleanup(utils.CleanUp)

	conn := utils.StartUp()
	defer conn.Close()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.Link(conn, dbname, connStr)
	assert.Nil(t, err) // Ensure no error during linking first connection

	err = utils.Link(conn, dbname, connStr2)
	log.Println(err)
	assert.NotNil(t, err) // Ensure error during linking second connection
}

// TestNormLink tests linking a database with a normal user
func TestNormLink(t *testing.T) {
	dbname := "db2"
	connStr := "postgresql://postgres:1242003@localhost:5432/testdb"
	t.Cleanup(utils.CleanUp)

	conn := utils.StartUp()
	defer conn.Close()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during database creation

	err = utils.CreateUser(conn, "ragnar", "ragnar")
	assert.Nil(t, err) // Ensure no error during user creation

	conn = utils.StartUp()
	defer conn.Close()
	err = utils.Login(conn, "ragnar", "ragnar")
	assert.Nil(t, err) // Ensure no error during login as normal user

	err = utils.Link(conn, dbname, connStr)
	assert.NotNil(t, err) // Ensure error during linking database with normal user
}

// TestClean tests the cleanup utility
func TestClean(t *testing.T) {
	utils.CleanUp()
}
