package calculations

import (
	"database/sql"
	"fmt"
	"testify-tutorial/stocks"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "stocks"
)


func TestCalculate(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	setupDatabase(t, db)
	seedTestTable(t, db)

	pp := stocks.NewPriceProvider(db)
	calculator := NewPriceIncreaseCalculator(pp)

	actual, err := calculator.PriceIncrease()

	if err != nil {
		t.Logf("err must be nul; but was %s", err.Error())
	}

	if actual != 100.0 {
		t.Logf("price increase must be 100.0, but was %f", actual)
	}

	tearDownDatabase(t, db)

}

func setupDatabase(t *testing.T, db *sql.DB) {
	t.Log("setting up database")
	_, err := db.Exec(`CREATE DATABASE stocks_test`)
	if err != nil {
		t.Logf("unable to create database %s", err.Error())
		tearDownDatabase(t, db)
		t.FailNow()
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp TIMESTAMPTZ PRIMARY KEY,
		price INT NOT NULL
	)`)

	if err != nil {
		t.Logf("unable to create table %s", err.Error())
		tearDownDatabase(t, db)
		t.FailNow()
	}

}


func seedTestTable(t *testing.T, db *sql.DB) {
	t.Log("seeding test table")

	for i := 2; i <= 3; i++ {
		_, err := db.Exec("INSERT INTO stockprices (timestamp, price) VALUES ($1,$2)", time.Now().Add(time.Duration(i)*time.Minute), float64(i))
		if err != nil {
			t.Logf("unable to seed table")
		}
	}
}

func tearDownDatabase(t *testing.T, db *sql.DB) {
	t.Log("tearing down database")

	_, err := db.Exec(`DROP TABLE stockprices`)
	if err != nil {
		t.Logf("unable to drop table")
	}

	_, err = db.Exec(`DROP DATABASE stocks_test`)
	if err != nil {
		t.Logf("unable to drop database")
	}

	err = db.Close()
	if err != nil {
		t.Logf("unable to close database")
	}
}
