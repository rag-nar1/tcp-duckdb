package internal

import (
	"database/sql"
	"encoding/json"
	"strings"
	"github.com/xwb1989/sqlparser"
)

type SQLExecutor interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    Exec(query string, args ...interface{}) (sql.Result, error)
}


func SELECT(db SQLExecutor, query string) ([]byte, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil ,err
	}
	defer rows.Close()

	result, err := ReadRows(rows)
	if err != nil {
		return nil ,err
	}

	dataJson, err := json.Marshal(result)
	if err != nil {
		return nil ,err
	}

	return dataJson, nil
}

func EXEC(db SQLExecutor, query string) (int64, int64, error) {
    result, err := db.Exec(query)
    if err != nil {
		return 0, 0, err
	}

    LastInsertedID, err := result.LastInsertId()
    if err != nil {
		return 0, 0, err
	}

    RowsAffected,err := result.RowsAffected()
    if err != nil {
		return 0, 0, err
	}
    
    return LastInsertedID, RowsAffected, nil
}


func CREATE(db SQLExecutor, server *sql.DB, stmt *sql.Stmt, query string, DBID int) (error){
	tables, err := ExtractTableNames(query)
	if err != nil {
		return err
	}
	
	servertx, err := server.Begin();
	if err != nil {
		return err
	}
	defer servertx.Rollback()
	
	for _, table := range tables {
		_, err := servertx.Stmt(stmt).Exec(table, DBID)
		if err != nil {
			return err
		}
	}
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	err = servertx.Commit()
	if err != nil {
		return err
	}

	return nil
}


func ReadRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil ,err
	}

	columnsSz := len(columns)
	for rows.Next() {
		row := make([]interface{}, columnsSz)
		rowPtr := make([]interface{}, columnsSz)
		for i := range row {
			rowPtr[i] = &row[i]
		}

		rows.Scan(rowPtr...)
		object := make(map[string]interface{})
		for i , columnName := range columns {
			object[columnName] = row[i] 
		}

		result = append(result, object)
	}

	return result, nil
}

func CheckDDLActions(query string) (bool, error) {
	DDL := make([]string, 0)
	// Parse the query into an AST.
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return false, err
	}
	// Walk the AST to catch both top-level and nested statements.
	err = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		var action string
		switch n := node.(type) {
		// Data Definition Language:
		case *sqlparser.DDL:
			// n.Action might be "create", "alter", "drop", etc.
			action = strings.ToUpper(n.Action)
			DDL = append(DDL, action)
		default:
			return true, nil
		}
		return true, nil
	}, stmt)
	if err != nil {
		return false, err
	}

	return (len(DDL) > 0), nil
}

func CheckAccesOverTable(db *sql.DB, stmt *sql.Stmt, query string, UID, DBID int) (bool, error){
	// parse the table and map each to type of acction used to access it
	tables , err := classifySQLTables(query)
	if err != nil {
		return false, err
	}

	var TID int
	var cnt int
	for table , actions := range tables {
		err := 	db.QueryRow("SELECT tableid FROM tables WHERE tablename LIKE ? AND dbid == ?", table, DBID).Scan(&TID)
		if err != nil {
			return false, err
		}

		for _ , action := range actions {
			err = stmt.QueryRow(UID, TID, strings.ToLower(action)).Scan(&cnt)
			if err != nil {
				return false, err
			}
			if cnt == 0 {
				return false, nil
			}
		}
	}

	return true, nil
}

func ExtractTableNames(query string) ([]string, error) {
    stmt, err := sqlparser.Parse(query)
    if err != nil {
        return nil, err
    }

    tableNames := make(map[string]bool)

    // Walk through the AST and collect table names
    sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
        switch n := node.(type) {
        case *sqlparser.TableName:
            if !n.IsEmpty() {
                tableName := n.Name.String()
                tableNames[tableName] = true
            }
        case *sqlparser.DDL:
            // Handle CREATE, DROP, ALTER table statements
            if !n.NewName.IsEmpty() {
                tableName := n.NewName.Name.String()
                tableNames[tableName] = true
            }
            if !n.Table.IsEmpty() {
                tableName := n.Table.Name.String()
                tableNames[tableName] = true
            }
        }
        return true, nil
    }, stmt)

    // Convert map to slice
    result := make([]string, 0, len(tableNames))
    for tableName := range tableNames {
        result = append(result, tableName)
    }

    return result, nil
}

// extractTableNamesFromTableExprs extracts table names from a list of TableExprs.
// It handles AliasedTableExpr, JoinTableExpr, and ParenTableExpr.
func extractTableNamesFromTableExprs(exprs sqlparser.TableExprs) []string {
	var tables []string
	for _, expr := range exprs {
		switch tbl := expr.(type) {
		case *sqlparser.AliasedTableExpr:
			// If the expression is a simple table name.
			switch expr := tbl.Expr.(type) {
			case sqlparser.TableName:
				// Use sqlparser.String to get the fully qualified name.
				tables = append(tables, sqlparser.String(expr))
			}
		case *sqlparser.JoinTableExpr:
			tables = append(tables, extractTableNamesFromTableExprs(sqlparser.TableExprs{tbl.LeftExpr})...)
			tables = append(tables, extractTableNamesFromTableExprs(sqlparser.TableExprs{tbl.RightExpr})...)
		case *sqlparser.ParenTableExpr:
			tables = append(tables, extractTableNamesFromTableExprs(tbl.Exprs)...)
		}
	}
	return tables
}

// classifySQLTables parses the given SQL query and returns a map that maps table names
// to a list of actions (like SELECT, INSERT, CREATE, etc.) that are performed on them.
func classifySQLTables(query string) (map[string][]string, error) {
	// Map of table name to a slice of actions.
	result := make(map[string][]string)

	// Parse the query into an AST.
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	// Walk the AST to catch both top-level and nested statements.
	err = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		var action string
		var tables []string

		switch n := node.(type) {
		// Data Query Language: SELECT statements.
		case *sqlparser.Select:
			action = "SELECT"
			tables = extractTableNamesFromTableExprs(n.From)
		// Data Manipulation Language:
		case *sqlparser.Insert:
			action = "INSERT"
			tables = []string{sqlparser.String(n.Table)}
		case *sqlparser.Update:
			action = "UPDATE"
			tables = extractTableNamesFromTableExprs(n.TableExprs)
		case *sqlparser.Delete:
			action = "DELETE"
			tables = extractTableNamesFromTableExprs(n.TableExprs)
		default:
			return true, nil
		}

		// Append the action to each table found.
		for _, t := range tables {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			// Append action to the list for the table.
			result[t] = append(result[t], action)
		}

		return true, nil
	}, stmt)
	if err != nil {
		return nil, err
	}

	return result, nil
}