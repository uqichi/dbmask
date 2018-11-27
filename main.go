package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/somali")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	res, _ := db.Query("SHOW TABLES")

	var table string

	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)
	}
}
