# simpledbgo

[![coverage](https://raw.githubusercontent.com/goropikari/simpledbgo/gh-pages/coverage.svg)](https://goropikari.github.io/simpledbgo/coverage/)

This is Go implementation of [SimpleDB](http://cs.bc.edu/~sciore/simpledb/) by [Database Design and Implementation](https://link.springer.com/book/10.1007/978-3-030-33836-7).


## Appendix
### Original SimpleDB setup
```
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
