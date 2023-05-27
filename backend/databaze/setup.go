package databaze

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func DBConnect() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error .env", err)
	}
	connStr := fmt.Sprintf("postgresql://%s:%s@%s", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"))
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Databaze se pokazila", err)
	}
}
