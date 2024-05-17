package main

import (
	"database/sql"
	"fmt"
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

	row := db.QueryRow(`select count(*) from 'https://raw.githubusercontent.com/anonranger/Go-DuckDB-Lambda/main/student-data.csv'`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return strconv.Itoa(count), nil
}

func main() {
	fmt.Println(hello())
	lambda.Start(hello)
}
