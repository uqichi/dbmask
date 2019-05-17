package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/uqichi/dbmask/mask"
)

var (
	nomask = flag.String("n", "", "nomask table fields")
)

func init() {
	flag.Parse()
}

func main() {
	nomasks := strings.Split(*nomask, ",")

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3313)/sakila")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = mask.Mask(context.Background(), db, nomasks)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
