# Packet capture postgres wire protocol

```bash
$ docker compose build
$ docker compose up

$ docker compose exec client bash
# psql -h db -Upostgres
postgres=# BEGIN;
BEGIN

$ docker compose exec client bash
# ngrep -x -q '.' 'host db'
```


```sql
postgres=# BEGIN;
BEGIN

postgres=*# select * from hoge;
 id
----
(0 rows)

postgres=*# ROLLBACK;
ROLLBACK
postgres=#
```

```
T 172.18.0.2:54100 -> 172.18.0.3:5432 [AP] #47
  51 00 00 00 0b 42 45 47    49 4e 3b 00                Q....BEGIN;.

T 172.18.0.3:5432 -> 172.18.0.2:54100 [AP] #49
  43 00 00 00 0a 42 45 47    49 4e 00 5a 00 00 00 05    C....BEGIN.Z....
  54                                                    T



T 172.18.0.2:54100 -> 172.18.0.3:5432 [AP] #51
  51 00 00 00 18 73 65 6c    65 63 74 20 2a 20 66 72    Q....select * fr
  6f 6d 20 68 6f 67 65 3b    00                         om hoge;.

T 172.18.0.3:5432 -> 172.18.0.2:54100 [AP] #53
  54 00 00 00 1b 00 01 69    64 00 00 00 40 00 00 01    T......id...@...
  00 00 00 17 00 04 ff ff    ff ff 00 00 43 00 00 00    ............C...
  0d 53 45 4c 45 43 54 20    30 00 5a 00 00 00 05 54    .SELECT 0.Z....T



T 172.18.0.2:54100 -> 172.18.0.3:5432 [AP] #55
  51 00 00 00 0e 52 4f 4c    4c 42 41 43 4b 3b 00       Q....ROLLBACK;.

T 172.18.0.3:5432 -> 172.18.0.2:54100 [AP] #57
  43 00 00 00 0d 52 4f 4c    4c 42 41 43 4b 00 5a 00    C....ROLLBACK.Z.
  00 00 05 49                                           ...I
```
