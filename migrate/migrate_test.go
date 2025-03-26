package migrate_test

import (
	"database/sql"
	"encoding/json"

	"github.com/rag-nar1/TCP-Duckdb/utils"

	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func Insert(connStr string) error {
	pq, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer pq.Close()

	_, err = pq.Exec("insert into t1(id) values(777);")
	if err != nil {
		return err
	}
	return nil
}

func Insert100(connStr string) error {
	pq, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer pq.Close()
	for i := 4; i <= 100; i++ {
		_, err = pq.Exec("insert into t1(id) values($1);", i)
		if err != nil {
			return err
		}
	}

	return nil
}

type TestRow struct {
	Id int `json:"id"`
}

func TestBasicMigrate(t *testing.T) {
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

	err = Insert(connStr)
	assert.Nil(t, err) // Ensure no error during linking database

	err = utils.Migrate(conn, dbname)
	assert.Nil(t, err)

	err = utils.ConnectDb(conn, dbname)
	assert.Nil(t, err) // Ensure no error during connecting to the database

	data, err := utils.QueryData(conn, "select * from t1 where id == 777;")
	assert.Nil(t, err)

	res := make([]TestRow, 1)
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, 777, res[0].Id)
}

func TestBigMigration(t *testing.T) {
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

	err = Insert100(connStr)
	assert.Nil(t, err) // Ensure no error during linking database

	err = utils.ConnectDb(conn, dbname)
	assert.Nil(t, err) // Ensure no error during connecting to the database

	data, err := utils.QueryData(conn, "select max(id) as id from t1;")
	assert.Nil(t, err)

	res := make([]TestRow, 1)
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 3, res[0].Id)
	conn.Close()

	conn = utils.Connection()
	err = utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.Migrate(conn, dbname)
	assert.Nil(t, err)

	err = utils.ConnectDb(conn, dbname)
	assert.Nil(t, err) // Ensure no error during connecting to the database

	data, err = utils.QueryData(conn, "select max(id) as id from t1;")
	assert.Nil(t, err)

	if err := json.Unmarshal([]byte(data), &res); err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, 100, res[0].Id)
}
