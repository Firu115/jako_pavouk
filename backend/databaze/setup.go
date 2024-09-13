package databaze

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// inicializuje připojení k databázi
func DBConnect() {
	var err error
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_JMENO"))
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Databaze se pokazila", err)
	}

	log.Printf("Připojeno k db %q", os.Getenv("DB_JMENO"))
}
