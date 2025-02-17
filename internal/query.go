package internal

import (
	"database/sql"
	"encoding/json"
	"strings"
	"github.com/xwb1989/sqlparser"
)


func SELECT(db *sql.DB, query string) ([]byte, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil ,err
	}

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

func EXEC(db *sql.DB, query string) (int64, int64, error) {
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

	return true, nil
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