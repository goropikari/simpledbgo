import psycopg

with psycopg.connect("host=127.0.0.1 port=5432") as conn:
    with conn.cursor() as cur:
        cur.execute("create table piyo (id int, name varchar(10))")
        cur.execute("INSERT INTO piyo (id, name) VALUES (100, 'taro')")
        cur.execute("SELECT * FROM piyo")
        cur.fetchone()

        for record in cur:
            print(record)

        conn.commit()
