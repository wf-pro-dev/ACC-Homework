package crud

import (
	"database/sql"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	_ "github.com/lib/pq"
)

func GetDB() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "williamfotso"
		password = "076ifDf9f5ovd8zn4zrO8S5RS0CaT9u2"
		dbname   = "notion_homework"
	)

	// Info Formatting
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// ERROR handler
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil

}

// ***************************************************************
//
// name : GetHandler
//
// params :
//
//	query (string) SQL query to execute
//	db (*sql.DB) SQL db connection object
//
// return :
//
//	The contains of the query response in a list/slice of map
//	Error in case of one , <nil> otherwise
//
// ***************************************************************
func GetHandler(query string, db *sql.DB) ([]map[string]string, error) {

	rows, err := db.Query(query) // Excution of the SELECT query

	if err != nil {
		log.Fatalln(err)
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column names: %w", err)
	}

	// Create a slice to hold all rows
	var results []map[string]string

	// Create slices to store values and scan pointers
	values := make([]string, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Iterate through rows
	for rows.Next() {
		// Scan the row into scanArgs
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Create a map for this row
		row := make(map[string]string)

		// For each column, add column name and value to the map
		for i, colName := range columns {
			val := values[i]

			// Handle NULL values (will be represented as nil)

			// Handle byte slices (common for string data from databases)
			// This converts []byte to string which is usually more useful
			if val != "" {
				row[colName] = val
			}

		}

		// Add the row map to results
		results = append(results, row)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil

} // End getHandler

// ***************************************************************
//
//	name : PostHandler
//	params :
//
//	newItem (map[string]string) data to insert in the row
//	table (string) name of the table to insert the row
//	db (*sql.DB) SQL db connection object
//
//	Return :
//
//	Error in case of one , <nil> otherwise
//
// ***************************************************************
func PostHandler(newItem map[string]string, table string, db *sql.DB) error {

	//fmt.Printf("%v", newItem)

	colomns := slices.Collect(maps.Keys(newItem)) // name of the colums (slice)
	colomnsStr := strings.Join(colomns, ", ")     // name of the colums (string)

	values := []string{} // data to insert (slice)
	for _, col := range colomns {
		val := fmt.Sprintf("'%v'", newItem[col])
		values = append(values, val)
	}
	valuesStr := strings.Join(values, ", ") // data to insert (string)

	// Test if there is data to insert
	if len(newItem) > 0 {

		query := fmt.Sprintf("INSERT into %v (%v) VALUES (%v)", table, colomnsStr, valuesStr)

		// Excution of the INSERT query
		_, err := db.Query(query)
		// Handle error
		if err != nil {
			return fmt.Errorf("error while exec POST query: %w", err)
		}
	}

	return nil
}

// ***************************************************************
//
// name : PutHandler
//
// params :
//
//	newItem (map[string]string) data to insert in the row
//	table (string) name of the table to insert the row
//	db (*sql.DB) SQL db connection object
//
// return :
//
//	Error in case of one , <nil> otherwise
//
// ***************************************************************
func PutHanlder(id int, col, table string, newValue string, db *sql.DB) (err error) {
	query := fmt.Sprintf("UPDATE %v SET %v = '%v' WHERE id='%v'", table, col, newValue, id)
	// Excution of the UPDATE query
	_, err = db.Query(query)
	// Handle error
	if err != nil {
		log.Fatalf("An error occured while executing query: %v", err)
	}

	return err
}

// ***************************************************************
//
// name : DeleteHandler
//
// params :
//
//	table (string) name of the table to delete the row
//	column (string) name of the column to delete the row
//	value (string) value of the column to delete the row
//	db (*sql.DB) SQL db connection object
//
// return :
//
//	True if the operation was succesful, False otherwise
//
// ***************************************************************
func DeleteHandler(table, column, value string, db *sql.DB) error {

	query := fmt.Sprintf("DELETE FROM %v WHERE %v='%v'", table, column, value)

	// Excution of the DELETE query
	_, err := db.Query(query)
	// Handle error
	if err != nil {
		log.Fatalf("An error occured while executing query: %v", err)

	}

	return err
}
