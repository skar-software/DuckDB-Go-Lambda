package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/marcboeker/go-duckdb"
)

func hello() (string, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return "", err
	}
	defer db.Close()

	db.Exec("INSTALL httpfs")
	db.Exec("LOAD httpfs")

	// Get count of records in csv
	row := db.QueryRow(`select count(*) from 'https://raw.githubusercontent.com/anonranger/Go-DuckDB-Lambda/main/student-data.csv'`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("No of records in CSV: " + strconv.Itoa(count) + "\n\n")

	// Get a random record
	query := `
        SELECT *
        FROM 'https://raw.githubusercontent.com/anonranger/Go-DuckDB-Lambda/main/student-data.csv'
        ORDER BY RANDOM()
        LIMIT 1;
    `

    // Execute the query
    rows, err := db.Query(query)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        log.Fatal(err)
    }

    // Create a slice of interface{}'s to represent each column, and a second slice to contain pointers to each item in the columns slice
    values := make([]interface{}, len(columns))
    valuePtrs := make([]interface{}, len(columns))
    for i := range values {
        valuePtrs[i] = &values[i]
    }

    // Iterate over the rows
    for rows.Next() {
        // Scan the result into the column pointers
        err = rows.Scan(valuePtrs...)
        if err != nil {
            log.Fatal(err)
        }

		fmt.Println("Random record from CSV\n")
        // Print each column's value
        for i, col := range columns {
            var val interface{}
            // Retrieve the value
            val = *(valuePtrs[i].(*interface{}))

            // Print the value
            fmt.Println(col, val)
        }
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

	return "See logs below for output", nil
}

func main() {
	lambda.Start(hello)
}
