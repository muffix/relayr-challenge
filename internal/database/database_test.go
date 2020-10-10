package database

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// testDatabasePath points to an in-memory database for testing
const testDatabasePath = "file::memory:?mode=memory&cache=shared"

func setupTestDatabase(t *testing.T, testData []Offer) Offers {
	db, err := InitSQLiteDatabase(testDatabasePath)
	if err != nil {
		t.Fatalf("Expected no error creating the database, got %s", err.Error())
	}

	err = insertDummyData(db, testData)

	if err != nil {
		t.Fatalf("Expected no error inserting test data, got %s", err.Error())
	}

	return db
}

func insertDummyData(db Offers, testData []Offer) error {
	return db.InsertMultiple(testData)
}

func TestDatabase(t *testing.T) {
	testData := []Offer{
		{"Vogon Poetry", "Better not haves", "Hitchhiker Essentials", 1},
		{"Babelfish", "Must Haves", "Hitchhiker Essentials", 1},
		{"Towel", "Must Haves", "Hitchhiker Essentials, just more expensive", 43},
		{"Towel", "Must Haves", "Hitchhiker Essentials, just more expensive", 44},
		{"21 is only half the Truth", "Books", "Hitchhiker Essentials", 2},
		{"Towel", "Must Haves", "Hitchhiker Essentials", 42},
	}

	db := setupTestDatabase(t, testData)
	defer db.Close()

	offers, err := db.Get("Towel", "Must Haves")
	if err != nil {
		t.Fatalf("Expected no error retrieving offers, got %v", err)
	}

	if len(offers) != 2 {
		t.Fatalf("Expected exactly two offers, got %d", len(offers))
	}

	expectedOffers := []Offer{
		{"Towel", "Must Haves", "Hitchhiker Essentials", 42},
		{"Towel", "Must Haves", "Hitchhiker Essentials, just more expensive", 44},
	}

	if !reflect.DeepEqual(offers, expectedOffers) {
		t.Fatalf("Expected the retrieved offers to be the same as the ones we put in")
	}
}

func TestOffersSQLiteDatabase_Insert(t *testing.T) {
	db, err := InitSQLiteDatabase(testDatabasePath)
	if err != nil {
		t.Fatalf("Expected no error creating the database, got %s", err.Error())
	}

	err = db.Insert("mock", "mock", "mock", 0)
	if err != nil {
		t.Fatalf("Expected no error inserting an offer, got %s", err.Error())
	}

	database := (*sql.DB)(db)
	rows, err := database.Query("SELECT * FROM offers")
	if err != nil {
		t.Fatalf("Expected no error when querying offer, got %v", err)
	}

	var rowsCount int
	got := Offer{}

	for rows.Next() {
		err = rows.Scan(&got.Product, &got.Category, &got.Supplier, &got.Price)
		rowsCount++
	}

	if rowsCount != 1 {
		t.Fatalf("Expected 1 row back, got %d", rowsCount)
	}

	if err != nil {
		t.Fatalf("Expected no error when reconstructing, got %v", err)
	}

	want := Offer{"mock", "mock", "mock", 0}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected the same offer back, but got %v", got)
	}
}

func TestDatabase_withoutData(t *testing.T) {
	db := setupTestDatabase(t, []Offer{})
	defer db.Close()

	offers, err := db.Get("Towel", "Must Haves")
	if err != nil {
		t.Fatalf("Expected no error retrieving offers, got %v", err)
	}

	if len(offers) != 0 {
		t.Fatalf("Expected to find no offers, got %d", len(offers))
	}
}
