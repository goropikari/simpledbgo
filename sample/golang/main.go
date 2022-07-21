package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=dummy password=dummy dbname=dummy sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec("create table T1(A int, B varchar(9))"); err != nil {
		log.Fatal(err)
	}

	n := 3
	for i := 0; i < n; i++ {
		cmd := fmt.Sprintf("insert into T1(A, B) values (%v, 'rec%v')", i, i)
		if _, err := db.Exec(cmd); err != nil {
			log.Fatal(err)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	rows2, err := tx.QueryContext(context.Background(), "select A, B from T1")
	if err != nil {
		log.Fatal(err)
	}
	for rows2.Next() {
		var a int
		var b string
		if err := rows2.Scan(&a, &b); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v\n", a, b)
	}
	if err := rows2.Err(); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
