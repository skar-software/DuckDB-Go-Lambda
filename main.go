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
	db.Exec("CREATE SECRET (TYPE S3, ENDPOINT 'r2.dev')")

	row := db.QueryRow(`select count(*) from 's3://pub-<redacted>/student-data.csv'`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return strconv.Itoa(count), nil
}

func main() {
	lambda.Start(hello)
}
