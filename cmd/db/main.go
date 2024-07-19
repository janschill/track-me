package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/janschill/track-me/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

var (
    dbPath  string
    operation string
)

func init() {
    flag.StringVar(&dbPath, "dbpath", "./data/trips.db", "Path to the database file.")
    flag.StringVar(&operation, "operation", "", "Database operation to perform: create, reset, destroy.")
}

func main() {
    flag.Parse()

    if dbPath == "" || operation == "" {
        fmt.Println("Usage: go run main.go -dbpath=<path-to-db> -operation=<operation>")
        os.Exit(1)
    }

    switch operation {
    case "create":
				db.CreateTables(dbPath)
    case "reset":
        db.ResetDB(dbPath)
    case "destroy":
        db.DestroyDB(dbPath)
    case "seed":
        db.Seed(dbPath)
    default:
        fmt.Println("Invalid operation. Available operations: init, create, setup, reset.")
    }
}
