package internal

import (
	"TCP-Duckdb/utils"
	
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

type JSONB map[string]interface{}
// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
    if j == nil {
        return nil, nil
    }
    return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(src interface{}) error {
    if src == nil {
        *j = nil
        return nil
    }

    // Handle different possible source types
    switch v := src.(type) {
    case []byte:
        if len(v) == 0 {
            *j = nil
            return nil
        }
        return json.Unmarshal(v, j)
    case string:
        if v == "" {
            *j = nil
            return nil
        }
        return json.Unmarshal([]byte(v), j)
    default:
        return fmt.Errorf("invalid type for JSONB")
    }
}

func (j JSONB) Get() ([]string, []interface{}) {
	columns := make([]string, len(j))
	values := make([]interface{}, len(j))
	i := 0
	for k, v := range j {
		columns[i] = k
		values[i] = v
		i ++
	} 
	return columns, values
}

type AuditRecord struct {
    EventID          int64      `json:"event_id" db:"event_id"`          // BIGSERIAL PRIMARY KEY
    SchemaName       string     `json:"schema_name" db:"schema_name"`    // TEXT NOT NULL
    TableName        string     `json:"table_name" db:"table_name"`      // TEXT NOT NULL
    TablePK          string     `json:"table_pk" db:"table_pk"`          // TEXT
    TablePKColumn    string     `json:"table_pk_column" db:"table_pk_column"` // TEXT
    ActionTimestamp  time.Time  `json:"action_tstamp" db:"action_tstamp"` // TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    Action           string     `json:"action" db:"action"`              // TEXT NOT NULL CHECK (action IN ('I','D','U','T'))
    OriginalData     JSONB      `json:"original_data" db:"original_data"` // JSONB
    NewData          JSONB      `json:"new_data" db:"new_data"`          // JSONB
    ChangedFields    JSONB      `json:"changed_fields" db:"changed_fields"` // JSONB
    TransactionID    *int64     `json:"transaction_id" db:"transaction_id"` // BIGINT (nullable)
    ApplicationName  *string    `json:"application_name" db:"application_name"` // TEXT (nullable)
    ClientAddr       *net.IP    `json:"client_addr" db:"client_addr"`    // INET (nullable)
    ClientPort       *int32     `json:"client_port" db:"client_port"`    // INTEGER (nullable)
    SessionUserName  *string    `json:"session_user_name" db:"session_user_name"` // TEXT (nullable)
    CurrentUserName  *string    `json:"current_user_name" db:"current_user_name"` // TEXT (nullable)
}


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

func ReadAudit(duck, postgres *sqlx.DB) error {
	var records []AuditRecord
	err := postgres.Select(&records, "select * from audit.logged_actions;")
	if err != nil {
		return err
	}
	var ctx context.Context
	transaction, err := duck.BeginTxx(ctx,&sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	for i := range records {
		switch records[i].Action {
		case "I":
			ApplyInsert(transaction, &records[i])
		case "U":
			ApplyUpdate(transaction, &records[i])
		case "D":
			ApplyDelete(transaction, &records[i])
		case "T":
			ApplyTrancate(transaction, &records[i])
		default:
			return fmt.Errorf("Unsupported Action")
		}
	}

	return nil
}

func GenPlaceHoldersForDuck(num int) string {
	placeholders := make([]string, num)
	for i := range placeholders {
		placeholders[i] = "?"
	}

	res := strings.Join(placeholders, ",")
	return res
}

func GenSetForDuck(keys []string) string {
	placeholders := make([]string, len(keys))
	for i := range placeholders {
		placeholders[i] = keys[i] + " = ?"
	}
	res := strings.Join(placeholders, ",")
	return res
}

func ApplyInsert(db *sqlx.Tx, record *AuditRecord) error {
	query := "INSERT INTO %s(%s) VALUES(%s);"
	keys, valuesInterfaces := record.NewData.Get()
	columns := strings.Join(keys, ",")
	query = fmt.Sprintf(query, record.TableName, columns, GenPlaceHoldersForDuck(len(keys)))
	_, err := db.NamedExec(query,valuesInterfaces)
	return err
}

func ApplyUpdate(db *sqlx.Tx, record *AuditRecord) error {
	query := "UPDATE %s SET %s WHERE %s;"
	keys, valuesInterfaces := record.ChangedFields.Get()
	predicate := fmt.Sprintf("%s = %s", record.TablePKColumn, record.TablePK)
	columns := GenSetForDuck(keys)
	query = fmt.Sprintf(query, record.TableName, columns, predicate)
	_, err := db.NamedExec(query,valuesInterfaces)
	return err
}

func ApplyDelete(db *sqlx.Tx, record *AuditRecord) error {
	query := "DELETE FROM %s WHERE %s;"
	predicate := fmt.Sprintf("%s = %s", record.TablePKColumn, record.TablePK)
	query = fmt.Sprintf(query, record.TableName, predicate)
	_, err := db.Exec(query)
	return err
}

func ApplyTrancate(db *sqlx.Tx, record *AuditRecord) error {
	query := "TRUNCATE %s;"
	query = fmt.Sprintf(query, record.TableName)
	_, err := db.Exec(query)
	return err
}