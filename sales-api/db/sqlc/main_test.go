package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries

const (
	dbDriver = "mysql"
	dbSource = "root:admin@tcp(127.0.0.1:3306)/sales-db?parseTime=true"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
