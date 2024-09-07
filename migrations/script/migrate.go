package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	var (
		dsn     string
		path    string
		command string
	)
	flag.StringVar(&dsn, "dsn", "", "database connection string")
	flag.StringVar(&path, "path", "", "migration file path")
	flag.StringVar(&command, "command", "", "command to run (up/down)")
	flag.Parse()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(file))
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	log.Println("migration done")
}
