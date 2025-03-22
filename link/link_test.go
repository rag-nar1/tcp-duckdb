package link_test

import (
	"TCP-Duckdb/utils"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func TestBasicLink(t *testing.T) {
	dbname := "db2"
	connStr := "postgresql://postgres:1242003@localhost:5432/testdb"
	t.Cleanup(utils.CleanUp)

	conn := utils.StartUp()
	defer conn.Close()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.Link(conn, dbname, connStr)
	assert.Nil(t, err) // Ensure no error during linking db

	err = utils.ConnectDb(conn, dbname);
	assert.Nil(t, err)

	err = utils.Query(conn, "select * from t1;")
	assert.Nil(t, err) // Ensure no error during linking db
}

func TestDoubleLink(t *testing.T) {
	dbname := "db2"
	connStr := "postgresql://postgres:1242003@localhost:5432/testdb"
	connStr2 := "postgresql://postgres:1242003@localhost:5432/testdb2"
	t.Cleanup(utils.CleanUp)

	conn := utils.StartUp()
	defer conn.Close()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.Link(conn, dbname, connStr)
	assert.Nil(t, err) // Ensure no error during linking db

	err = utils.Link(conn, dbname, connStr2)
	log.Println(err)
	assert.NotNil(t, err) // Ensure no error during linking db
}

func TestNormLink(t *testing.T) {
	dbname := "db2"
	connStr := "postgresql://postgres:1242003@localhost:5432/testdb"
	t.Cleanup(utils.CleanUp)

	conn := utils.StartUp()
	defer conn.Close()
	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.CreateDB(conn, dbname)
	assert.Nil(t, err) // Ensure no error during creating db

	err = utils.CreateUser(conn, "ragnar", "ragnar")
	assert.Nil(t, err)
	conn = utils.StartUp()
	defer conn.Close()
	err = utils.Login(conn,"ragnar", "ragnar")
	assert.Nil(t, err) // Ensure no error during startup

	err = utils.Link(conn, dbname, connStr)
	assert.NotNil(t, err) // Ensure no error during linking db
}

func TestClean(t *testing.T) {
	utils.CleanUp()
}