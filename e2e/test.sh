#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $0); pwd)
echo $SCRIPT_DIR

cd $SCRIPT_DIR/..


# make docker-build
id=$(docker run --rm -p 5432:5432 -d simpledb)
sleep 1

echo "start e2e test"
PSQL="psql --csv -q -h localhost -c "

# fixture
$PSQL "create table hoge (id int, name varchar(10));"
$PSQL "insert into  hoge (id, name) values (1, 'foo');"
$PSQL "insert into  hoge (id, name) values (2, 'bar');"

$PSQL "create table piyo (id2 int, name2 varchar(10));"
$PSQL "insert into  piyo (id2, name2) values (1, 'foo');"
$PSQL "insert into  piyo (id2, name2) values (2, 'bar');"

# valid test 1
ret=$($PSQL "select id, name, id2, name2 from hoge, piyo;")
expected="id,name,id2,name2 1,foo,1,foo 1,foo,2,bar 2,bar,1,foo 2,bar,2,bar"
ret=$(echo -n $ret)
expected=$(echo -n $expected)

if [ "$ret" != "$expected" ]; then
    docker stop $id
    exit 1
fi

# valid test 2
ret=$($PSQL "select * from hoge, piyo;")
expected="id,name,id2,name2 1,foo,1,foo 1,foo,2,bar 2,bar,1,foo 2,bar,2,bar"
ret=$(echo -n $ret)
expected=$(echo -n $expected)

if [ "$ret" != "$expected" ]; then
    docker stop $id
    exit 1
fi

# error test 1
ret=$($PSQL "select id, name, id2, name2 from hoge, hogehoge;" 2>&1)
expected='ERROR: relation "hogehoge" does not exist'
ret=$(echo -n $ret)
expected=$(echo -n $expected)

if [ "$ret" != "$expected" ]; then
    docker stop $id
    exit 1
fi

# error test 2
ret=$($PSQL "select unknown from hoge;" 2>&1)
expected='ERROR: column "unknown" does not exist'
ret=$(echo -n $ret)
expected=$(echo -n $expected)

if [ "$ret" != "$expected" ]; then
    docker stop $id
    exit 1
fi

echo "stop e2e test"
docker stop $id
