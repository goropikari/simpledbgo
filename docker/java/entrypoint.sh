#!/bin/bash

cd /app/SimpleDB_3.4

find ./simpledb -name *.java | xargs javac
javac simpleclient/SimpleIJ.java

java simpledb.server.StartServer &
sleep 3
java simpleclient.SimpleIJ
