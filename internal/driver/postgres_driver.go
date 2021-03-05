package driver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "host.docker.internal"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "togo"
)

type PostgresDB struct {
	SQL *sql.DB
}

var Postgres = &PostgresDB{}

func Connect() *PostgresDB {
	psql := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.Fatal("error opening db", err)
	}

	Postgres.SQL = db
	return Postgres
}
