package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3313)/sakila")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// list all tables
	res, _ := db.Query("SHOW TABLES")
	var table string
	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)

		// each table
		rows, err := db.Query(
			"SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?",
			"sakila", table)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var columnName, dataType string
		for rows.Next() {
			if err = rows.Scan(&columnName, &dataType); err != nil {
				fmt.Println(err)
				continue
			}
			if strings.Contains(dataType, "char") || strings.Contains(dataType, "text") {
				// update string field with the value masked
				_, _ = db.Exec(fmt.Sprintf("UPDATE %s SET %s = CONCAT(LEFT(%s, 1), '*****')", table, columnName, columnName))
			}
		}

		//break
	}
}
