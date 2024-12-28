package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func EXEC (db *sql.DB , Query string) (error) {
	_ , err := db.Exec(Query)
	return err
}

func INSERT (db *sql.DB , Table string , Values ...string) (error) {
	Row := strings.Join(Values , " , ")
	_ , err := db.Exec("INSERT INTO %s VALUES(%s);" , Table , Row)
	if err != nil {
		return err
	}

	return  nil
}

func UPDATE (db *sql.DB , Table string , Condition  string , argc int , argv ...string) (error) {
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

	_ , err := db.Exec("UPDATE %s SET %s WHERE %s;" , Table , Row , Condition)
	if err != nil {
		return err
	}

	return  nil
}

func DELETE (db *sql.DB , Table string , Condition  string ) (error) {
	_ , err := db.Exec("DELETE FROM %s WHERE %s;" , Table , Condition)
	return err
}

func SELECT [T any] (db *sql.DB , Table string , Columns string , Condition string) (*T , error) { // works for single row query

	RawData , err := db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s;" , Columns , Table , Condition))
	if err != nil {
		return nil , err
	}
	var DataPtr *T = new(T)	
	RawData.Next()
	err = RawData.Scan(DataPtr)
	if err != nil {
		return nil , err
	}
	return DataPtr , nil
}

