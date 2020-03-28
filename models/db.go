package models

import (
	"fmt"
	"database/sql"
	"log"
	"flag"
	"os"

	// "github.com/gin-gonic/gin"

	// For DATABASE
	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbInstance *sql.DB

func setDb(db *sql.DB) {
	dbInstance = db
}

func RetrieveDB() (*sql.DB) {
	return dbInstance
}


func InitializeDB() {
	user := os.Getenv("db_user")
	dbname := os.Getenv("db_name")
	password := os.Getenv("db_password")
	sslmode := os.Getenv("db_sslmode")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	sqlDBConnectionString := "user=" + user + " dbname=" + dbname + " password=" + password + " sslmode" + sslmode
	db, err := sql.Open("postgres", sqlDBConnectionString)

	if err != nil {
		fmt.Println("Db not inititaed", err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error: Could not establish a connection with the database")
	} 

	
	// postgres://user:pass@localhost:5432/database?sslmode=disable
	sqlDBMigrationString := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname "?sslmode" + sslmode 
	db1, err := sql.Open("postgres", sqlDBMigrationString)

	if err != nil {
		fmt.Println("DB IS NOT INIT.")
	}

	driver, err := postgres.WithInstance(db1, &postgres.Config{})
	if err != nil {
		fmt.Println("ERror", err)
		return
	}

	fmt.Println("driver ", driver)

	migrationDir := flag.String("migration.files", "./migrations", "./migrations")
	fmt.Println("MIGRATION DDIR", migrationDir)

	m, err := migrate.NewWithDatabaseInstance(
        "file:///home/deepen/Documents/go_practises/go-gin-api/migrations",
		"postgres", driver)

	fmt.Println("M ", m)
	if err != nil {
		fmt.Println("ERror for new with db instance", err)
		return
	}	
	
	m.Up()

	// dbInstance = db
	setDb(db)

	// return dbInit, nil
}

