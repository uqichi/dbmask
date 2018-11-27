package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3313)/sakila")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	res, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Panic(err)
	}

	var table string
	for res.Next() {
		if err = res.Scan(&table); err != nil {
			log.Println(err)
			continue
		}

		rows, err := db.Query(
			"SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?",
			"sakila", table)
		if err != nil {
			log.Println(err)
			continue
		}

		var columnName, dataType string
		for rows.Next() {
			if err = rows.Scan(&columnName, &dataType); err != nil {
				log.Println(err)
				continue
			}
			if strings.Contains(dataType, "char") || strings.Contains(dataType, "text") {
				// update string field with the value masked
				_, _ = db.Exec(fmt.Sprintf("UPDATE %s SET %s = CONCAT(LEFT(%s, 1), '*****')", table, columnName, columnName))
			}
		}
	}
}
