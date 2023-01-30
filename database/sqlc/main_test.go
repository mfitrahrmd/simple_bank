package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

var queriesTest *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgres://dev:dev@localhost:5432/simple_bank?sslmode=disable")
	if err != nil {
		log.Fatalf("error connecting to database : %v", err)
	}

	queriesTest = New(conn)

	m.Run()
}