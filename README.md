# simpledbgo

[![coverage](https://raw.githubusercontent.com/goropikari/simpledbgo/gh-pages/coverage.svg)](https://goropikari.github.io/simpledbgo/coverage/)

This is Go implementation of [SimpleDB](http://cs.bc.edu/~sciore/simpledb/) by [Database Design and Implementation](https://link.springer.com/book/10.1007/978-3-030-33836-7).

## Run SimpleDB

### Server mode
```bash
$ make docker-build
$ make docker-run

# open another terminal
$ psql -h localhost
psql (14.3, server 0.0.0)
Type "help" for help.

arch=> create table foo (id int, name varchar(10));
OK
arch=> insert into foo (id, name) values (1, 'dog');
OK
arch=> insert into foo (id, name) values (2, 'cat');
OK
arch=> select name, id from foo;
 name | id
------+----
 dog  | 1
 cat  | 2
(2 rows)
```

### Embedded mode

```go
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
```

output
```
0 rec0
1 rec1
2 rec2
```


## Implementation Progress

| Book Chapter | Feature                                    | Implemented        |
|--------------|--------------------------------------------|--------------------|
| 3            | File Manager                               | :heavy_check_mark: |
| 4            | Log Manager                                | :heavy_check_mark: |
| 4            | Buffer Manager                             | :heavy_check_mark: |
| 5            | Recovery Manager                           | :heavy_check_mark: |
| 5            | Concurrency Manager                        | :heavy_check_mark: |
| 5            | Transaction                                | :heavy_check_mark: |
| 6            | Record Pages                               | :heavy_check_mark: |
| 6            | Table Scans                                | :heavy_check_mark: |
| 7            | Metadata Manager                           | :heavy_check_mark: |
| 8            | Select Scans, Project Scans, Product Scans | :heavy_check_mark: |
| 9            | Parser                                     | :heavy_check_mark: |
| 10           | Planner                                    | :heavy_check_mark: |
| 11           | JDBC Interface                             | :x:                |
| 12           | Static Hash Indexes                        | :heavy_check_mark: |
| 12           | Btree Indexes                              | :x:                |
| 13           | Materialization and Sorting                | :x:                |
| 14           | MultiBuffer Sorting/Product                | :x:                |
| 15           | Query Optimization                         | :x:                |

Instead of JDBC interface, I implemented a Go SQL driver interface and Postgres wire protocol.


## Appendix
### Original SimpleDB setup

```bash
wget http://www.cs.bc.edu/%7Esciore/simpledb/SimpleDB_3.4.zip
unzip SimpleDB_3.4.zip
sed -i -e "1i package simpleclient;" SimpleDB_3.4/simpleclient/SimpleIJ.java
docker build -t simpledb -f ./docker/java/Dockerfile ./docker/java
docker run --rm -it -v $(pwd)/SimpleDB_3.4:/app/SimpleDB_3.4 simpledb


recovering existing database
transaction 1 committed
database server ready
Connect>
jdbc:simpledb:foobar
creating new database
transaction 1 committed

SQL> create table baz (id int, name varchar(10));
transaction 2 committed
0 records processed

SQL> insert into baz (id, name) values (123, 'mike');
transaction 3 committed
1 records processed

SQL> insert into baz (id, name) values (456, 'joe');
transaction 4 committed
1 records processed

SQL> select id, name from baz;
     id       name
------------------
    123       mike
    456        joe
transaction 5 committed
```
