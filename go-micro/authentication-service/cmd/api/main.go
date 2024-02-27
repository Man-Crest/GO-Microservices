package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect to DB
	conn, err := ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer conn.Close()

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// func NewDatabase() (*sql.DB, error) {

// 	dsn := os.Getenv("DSN")

// 	// dsn := "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"

// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Println("error at sql.Open")
// 		return nil, err
// 	}

// 	if err = db.Ping(); err != nil {
// 		log.Println("error at Ping()")
// 		return nil, err
// 	}

// 	return db, nil
// }

func ConnectDB() (*sql.DB, error) {
	// Define PostgreSQL connection string
	connStr := "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC"

	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connectivity
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL database")
	return db, nil
}
