package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/uqichi/dbmask/mask"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	unmasks := []string{
		fmt.Sprintf("%s.%s", "staff", "first_name"),
		fmt.Sprintf("%s.%s", "actor", "last_name"),
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3313)/sakila")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	err = mask.Mask(context.Background(), db, unmasks)
	if err != nil {
		panic(err)
	}
}
