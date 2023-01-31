package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mfitrahrmd/simple_bank/config"
	"log"
	"testing"
)

var queriesTest *Queries
var dbTest *sql.DB

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatalf("error loading config : %v", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSourceName)
	if err != nil {
		log.Fatalf("error connecting to database : %v", err)
	}

	dbTest = conn

	queriesTest = New(conn)

	m.Run()
}
