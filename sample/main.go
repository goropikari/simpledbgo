package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/goropikari/simpledbgo/driver/embedded"
)

func main() {
	db, err := sql.Open("simpledb", "")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create table T1(A int, B varchar(9))")
	if err != nil {
		log.Fatal(err)
	}

	n := 3
	for i := 0; i < n; i++ {
		cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
		db.Exec(cmd)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := tx.QueryContext(context.Background(), "select A, B from T1")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var a int
		var b string
		err := rows.Scan(&a, &b)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v\n", a, b)
	}
}
