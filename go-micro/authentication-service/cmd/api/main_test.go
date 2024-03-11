package main

import (
	"authentication/data"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestConnectDB(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("host", "localhost")
	os.Setenv("port", "5432")
	os.Setenv("user", "testuser")
	os.Setenv("password", "testpassword")
	os.Setenv("dbname", "testdb")

	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Fatal("Expected non-nil DB connection, got nil")
	}

	// Perform additional tests if needed
}

func TestRoutes(t *testing.T) {
	// Create a dummy DB connection
	db, err := sql.Open("postgres", "dummy connection string")
	if err != nil {
		t.Fatalf("Error opening dummy DB connection: %v", err)
	}
	defer db.Close()

	app := Config{
		DB:     db,
		Models: data.New(db),
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.routes().ServeHTTP) // Use ServeHTTP method directly

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Additional tests for handler response if needed
}
