package db

import (
	"database/sql"
	"log"
	"os"
	"simple-bank/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("./../../")
	if err != nil {
		log.Fatal("error loading the config", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
