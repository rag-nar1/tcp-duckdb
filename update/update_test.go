package update_test

import (
	"encoding/json"
	"testing"

	"github.com/rag-nar1/tcp-duckdb/utils"
	"github.com/stretchr/testify/assert"
)

type TestRow struct {
	Id int `json:"id"`
}
func TestUpdateUser(t *testing.T) {
	t.Cleanup(utils.CleanUp)

	username := "mohamed"
	password := "mohamed"
	newusername := "fathy"
	newpassword := "fathy" 

	conn := utils.StartUp()
	defer conn.Close()

	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.CreateUser(conn, username, password)
	assert.Nil(t, err)

	conn2 := utils.Connection()
	defer conn2.Close()

	err = utils.Login(conn2, username, password)
	assert.Nil(t, err)
	conn2.Close()
	
	err = utils.UpdateUserName(conn, username, newusername)
	assert.Nil(t, err)

	conn2 = utils.Connection()
	err = utils.Login(conn2, newusername, password)
	assert.Nil(t, err)
	conn2.Close()

	conn2 = utils.Connection()
	err = utils.Login(conn2, username, password)
	assert.NotNil(t, err)
	conn2.Close()

	err = utils.UpdateUserName(conn, username, newusername)
	assert.NotNil(t, err)

	err = utils.UpdateUserPassword(conn, newusername, newpassword)
	assert.Nil(t, err)

	conn2 = utils.Connection()
	err = utils.Login(conn2, newusername, password)
	assert.NotNil(t, err)
	conn2.Close()

	conn2 = utils.Connection()
	err = utils.Login(conn2, newusername, newpassword)
	assert.Nil(t, err)
	conn2.Close()
}


func TestUpdateDatabase(t *testing.T) {
	t.Cleanup(utils.CleanUp)

	dbname := "db1"
	newdbname := "db2"

	conn := utils.StartUp()
	defer conn.Close()

	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.CreateDB(conn,dbname)
	assert.Nil(t, err)

	err = utils.UpdateDatabase(conn, dbname, newdbname)
	assert.Nil(t, err)

	err = utils.ConnectDb(conn, newdbname)
	assert.Nil(t, err)
	conn.Close()

	conn = utils.Connection()
	err = utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.ConnectDb(conn, dbname)
	assert.NotNil(t, err)
}

func TestUpdateDatabase2(t *testing.T) {
	t.Cleanup(utils.CleanUp)

	dbname := "db1"
	newdbname := "db2"

	conn := utils.StartUp()
	defer conn.Close()

	err := utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.CreateDB(conn,dbname)
	assert.Nil(t, err)

	err = utils.UpdateDatabase(conn, dbname, newdbname)
	assert.Nil(t, err)

	err = utils.ConnectDb(conn, newdbname)
	assert.Nil(t, err)
	
	err = utils.Query(conn, "create table t1(id int);")
	assert.Nil(t, err)
	
	err = utils.Query(conn, "insert into t1(id) values(1);")
	assert.Nil(t, err)
	conn.Close()

	conn = utils.Connection()
	err = utils.LoginAsAdmin(conn)
	assert.Nil(t, err) // Ensure no error during login as admin

	err = utils.ConnectDb(conn, newdbname)
	assert.Nil(t, err)

	data, err := utils.QueryData(conn, "select * from t1;")
	assert.Nil(t, err)

	res := make([]TestRow, 1)
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, 1, res[0].Id)
}

