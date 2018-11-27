package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	unmasks := []string{
		fmt.Sprintf("%s.%s", "staff", "first_name"),
		fmt.Sprintf("%s.%s", "film_text", "description"),
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

			// skip masking check
			s := fmt.Sprintf("%s.%s", table, columnName)
			skipMasking := func() bool {
				for _, v := range unmasks {
					if s == v {
						return true
					}
				}
				return false
			}()
			if skipMasking {
				log.Println("skip masking", s)
				continue
			}

			if strings.Contains(dataType, "char") || strings.Contains(dataType, "text") {
				// update string field with the value masked
				// TODO: fieldがnameを含んでたら a*******, phoneを含んでたら *****abc***** みたいなことできるようにする
				_, _ = db.Exec(fmt.Sprintf("UPDATE %s SET %s = CONCAT(LEFT(%s, 1), '*******')", table, columnName, columnName))
			}
		}
	}
}
