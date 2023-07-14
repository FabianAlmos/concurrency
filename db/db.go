package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var dbInstance *sql.DB = nil

func GetDB() *sql.DB {
	if dbInstance == nil {
		connStr := "postgres://dbdata:dbdatapswd@localhost/lesson?sslmode=disable"

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		dbInstance = db
	}
	return dbInstance
}
