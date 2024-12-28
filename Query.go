package main

import (
	"errors"
	"fmt"
	"strings"
	"context"
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

