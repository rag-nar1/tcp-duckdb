package main

import (
	"fmt"
	"strings"
	"database/sql"
	"log"
	"net"
)

func checkFetal (err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func checkLog (err error) {
	if err != nil {
		log.Println(err)
	}
}

/**
 * Function to create a tcp listener
 * @pram -> address string : address to listen to 
 * @return -> net.Listener , error
*/
func CreateListener(address string) (net.Listener , error) { 
	listner , err := net.Listen("tcp", address)
	if err != nil {
		return nil , err
	}
	return listner , nil
}



func addSingleQuete(argv string) (string) {
	resualt := "'" + argv + "'"
	return resualt
}

func ProccessQuery(queryResult *sql.Rows) (string , error) {
	columns , err := queryResult.Columns()
	columnsTypes , err := queryResult.ColumnTypes()
	if err != nil {
		return "" , err
	}

	var respones string = ""
	for _ , column := range columns {
		respones += column + " , "
	}
	respones = strings.TrimSuffix(respones , " , ")
	respones += "\n"

	for queryResult.Next() {
		rowData := make([]interface{} , len(columns))
		rowDataPtr := make([]interface{} , len(columns))
		for i := range rowData {
			rowDataPtr[i] = &rowData[i]
		}

		err = queryResult.Scan(rowDataPtr...)

		values := make(map[string]interface{})
		for i , column := range columns {
			values[column] = rowData[i]
		}

		for i , column := range columns {
			// interprate the value based on the type
			switch columnsTypes[i].DatabaseTypeName() {
				case "INT":
					respones += fmt.Sprintf("%d" , values[column].(int))
				case "BIGINT":
					respones += fmt.Sprintf("%d" , values[column].(int64))
				case "FLOAT":
					respones += fmt.Sprintf("%f" , values[column].(float64))
				case "VARCHAR":
					respones += fmt.Sprintf("%s" , values[column].(string))
				case "BOOL":
					respones += fmt.Sprintf("%t" , values[column].(bool))
				default:
					respones += fmt.Sprintf("%v" , values[column])
			}
			respones += " , "
		}
		respones = strings.TrimSuffix(respones , " , ")
		respones += "\n"
	}
	return respones , nil
}