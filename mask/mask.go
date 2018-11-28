package mask

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Mask(ctx context.Context, db *sql.DB, nomasks []string) error {

	dbNames, err := db.Query("SELECT DATABASE()")
	if err != nil {
		return err
	}

	var dbName string
	if dbNames.Next() {
		if err = dbNames.Scan(&dbName); err != nil {
			return err
		}
	}
	if dbName == "" {
		return errors.New("empty database")
	}

	tableNames, err := db.Query("SHOW FULL TABLES WHERE Table_Type = 'BASE TABLE'")
	if err != nil {
		return err
	}

	var tableName string
	for tableNames.Next() {
		if err = tableNames.Scan(&tableName, trashScanner{}); err != nil {
			log.Println(err)
			continue
		}

		rows, err := db.QueryContext(ctx,
			"SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?",
			dbName, tableName)
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
			skipper := fmt.Sprintf("%s.%s", tableName, columnName)
			skipMasking := func() bool {
				for _, v := range nomasks {
					if skipper == v {
						return true
					}
				}
				return false
			}()
			if skipMasking {
				log.Println("skip masking", skipper)
				continue
			}

			if strings.Contains(dataType, "char") || strings.Contains(dataType, "text") {
				valueLen := 3
				if strings.Contains(columnName, "name") {
					valueLen = 1
				}
				if strings.Contains(columnName, "phone") {
					valueLen = 5
				}
				// update string field with the value masked
				_, err = db.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET %s = CONCAT(LEFT(%s, %d), '******')", tableName, columnName, columnName, valueLen))
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}

	return nil
}

type trashScanner struct{}

func (trashScanner) Scan(interface{}) error {
	return nil
}
