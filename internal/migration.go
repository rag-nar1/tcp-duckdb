package internal

import (
	"TCP-Duckdb/utils"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
type column struct {
	name string
	dataType string
}

type table struct {
	name string
	columns []column
}

func (t *table) Add(name, dataType string) {
	t.columns = append(t.columns, column{name: name, dataType: dataType})
}


func (t table) GenereteSql() string {
	query := "CREATE TABLE IF NOT EXISTS %s (%s);"
	var columns string = ""
	for _, col := range t.columns {
		columns += col.name + " " + utils.DbTypeMap(strings.ToUpper(col.dataType)) + ","
	} 
	columns = columns[:len(columns) - 1]
	return fmt.Sprintf(query, t.name, columns)
}

func Migrate(DBID int, connStr string, stmt *sql.Stmt, postgres, duck, server *sql.Tx) (error) {
	var tables map[string]*table = make(map[string]*table)

	rows, err := postgres.Query("SELECT table_name, column_name, data_type FROM information_schema.columns WHERE table_schema = 'public' and table_name not like 'pg%';")
	if err != nil {
		return err
	}
	defer rows.Close()

	data := make([]string, 3)
	dataptr := make([]interface{}, 3);
	for i := range data {
		dataptr[i] = &data[i]
	}

	for rows.Next() {
		rows.Scan(dataptr...)
		_,ok := tables[data[0]]
		if !ok {
			tables[data[0]] = &table{name: data[0]}
		}
		tables[data[0]].Add(data[1],data[2])
	}

	// create tables in our server
	for _, table := range tables {
		_, err = server.Stmt(stmt).Exec(table.name, DBID)
		if err != nil {
			return err
		}
	}

	// connect to postgresql database to get all data
	_,err = duck.Exec(fmt.Sprintf("ATTACH '%s' AS postgres_db (TYPE postgres);", connStr))
	if err != nil {
		return err
	}

	for _, table := range tables {
		postgrestable := "postgres_db." + table.name 
		stmt := fmt.Sprintf("CREATE TABLE %s AS FROM %s;", table.name ,postgrestable)
		_, err := duck.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func Audit(postgres *sql.Tx) error {
	sqlfile, err := os.ReadFile(filepath.Join(os.Getenv("DBdir"), "server/audit.sql"))
	if err != nil {
		return err
	}

	_, err = postgres.Exec(string(sqlfile))
	if err != nil {
		return err
	}
	return nil
}