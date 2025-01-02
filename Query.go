package main

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"context"
	"database/sql"

)

func EXEC (db *SQLitedb , Query string) (error) {
	ctx := context.Background()
	_ , err := db.db.ExecContext(ctx , Query)
	return err
}

func INSERT (db *SQLitedb , Table string , Values ...string) (error) {
	Row := strings.Join(Values , " , ")
	db.latch.Lock()
	defer db.latch.Unlock()

	var ctx context.Context = context.Background()
	_ , err := db.db.ExecContext(ctx , fmt.Sprintf("INSERT INTO %s VALUES(%s);" , Table , Row))
	if err != nil {
		return err
	}

	return  nil
}

func UPDATE (db *SQLitedb , Table string , Condition  string , argc int , argv ...string) (error) {
	if argc % 2 != 0 || argc != len(argv) {
		return errors.New("bad args")
	}

	var Row string
	for i , arg := range argv {
		if i % 2 == 0 {
			Row = Row + " , " + arg + " = "
		}else{
			Row = Row + arg
		}
	}
	Row = strings.TrimPrefix(Row , " , ")

	db.latch.Lock()
	defer db.latch.Unlock()
	var ctx context.Context = context.Background()
	_ , err := db.db.ExecContext(ctx , fmt.Sprintf("UPDATE %s SET %s WHERE %s;" , Table , Row , Condition))
	if err != nil {
		return err
	}

	return  nil
}

func DELETE (db *SQLitedb , Table string , Condition  string ) (error) {
	db.latch.Lock()
	defer db.latch.Unlock()

	var ctx context.Context = context.Background()
	_ , err := db.db.ExecContext(ctx , fmt.Sprintf("DELETE FROM %s WHERE %s;" , Table , Condition))
	return err
}

func SELECT [T any] (db *SQLitedb , Table string , Columns string , Condition string) (*T , error) { // works for single row query
	var result *T = new(T)
	var ctx context.Context = context.Background()
	err := db.db.QueryRowContext(ctx , fmt.Sprintf("SELECT %s FROM %s WHERE %s;" , Columns , Table , Condition)).Scan(result)

	if err != nil {
		return nil , err
	}

	return result , nil
}

func CreateAdmin (db *SQLitedb , NewAdmin string , Password string) (error) {
	err := INSERT(sqlitedb , "admins" , strconv.Itoa(NextID["admins"]) ,NewAdmin , Password)
	if err != nil {
		return err
	}
	NextID["admins"] ++
	return nil
}

func CreateUser(db *SQLitedb , NewUser string , Password string) (error) {
	err := INSERT(sqlitedb , "users" , strconv.Itoa(NextID["users"]) , NewUser , Password)
	if err != nil {
		return err
	}
	NextID["users"] ++
	return nil
}

func CreateDB (db *SQLitedb , NewDB string , SuperUser string) (error) {
	DataPtr , err := SELECT[int](sqlitedb , "users" , "userid" , "username like " + SuperUser)
	if err != nil {
		return err
	}

	UID := *DataPtr
	DBID := NextID["database"]
	_ , err = sql.Open("duckdb" , dbpath + strconv.Itoa(UID) + "_" + NewDB) // creates the db file
	if err != nil {
		return err
	}

	err = INSERT(sqlitedb , "database" , strconv.Itoa(DBID) , NewDB , strconv.Itoa(UID))

	if err != nil {
		return err
	}

	err = INSERT(sqlitedb , "useraccess" ,  strconv.Itoa(DBID) , strconv.Itoa(UID))
	if err != nil {
		return err
	}
	
	
	NextID["database"] ++
	return nil

}

func CreateTable(db *SQLitedb , TableName string , DBName string) (error) {
	DataPtr , err := SELECT[int](sqlitedb , "database" , "dbid" , "dbname like " + DBName)

	if err != nil {
		return err
	}

	DBID := *DataPtr
	TableID := NextID["tables"]
	err = INSERT(sqlitedb , "tables" , strconv.Itoa(TableID) , addSingleQuete(TableName) , strconv.Itoa(DBID))

	if err != nil {
		return err
	}

	NextID["tables"] ++
	return nil
}