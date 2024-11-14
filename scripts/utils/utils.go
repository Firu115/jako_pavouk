package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func DBConnect() {
	var err error
	err = godotenv.Load("../../backend/.env")
	if err != nil {
		log.Fatal("Nenašel jsem soubor .env v /backend.")
	}
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_JMENO"))
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Databaze se pokazila", err)
	}

	log.Printf("Připojeno k db %q", os.Getenv("DB_JMENO"))
}

var uvozovka string = `"`
var pomlcka string = `-`
var tabulka = map[rune]string{
	'„': uvozovka,
	'”': uvozovka,
	'“': uvozovka,
	'‟': uvozovka,
	'❞': uvozovka,
	'❝': uvozovka,
	'🙷': uvozovka,
	'🙶': uvozovka,
	'＂': uvozovka,
	'‚': uvozovka,
	'’': uvozovka,
	'‘': uvozovka,
	'‛': uvozovka,
	'❛': uvozovka,
	'❜': uvozovka,
	'⹂': uvozovka,
	'…': "...",
	'‑': pomlcka,
	'‐': pomlcka,
	'־': pomlcka,
	'‒': pomlcka,
	'–': pomlcka,
	'—': pomlcka,
	'―': pomlcka,
	'⸺': pomlcka,
	'⸻': pomlcka,
	'﹘': pomlcka,
	'﹣': pomlcka,
	'－': pomlcka,
	'֊': pomlcka,
	'᠆': pomlcka,
	' ': " ",
	'ä': "a", // nemecke
	'ö': "o",
	'ü': "u",
	'ć': "c", // polske
	'ą': "a",
	'ę': "e",
	'ł': "l",
	'ń': "n",
	'ó': "o",
	'ź': "z",
	'ż': "z",
	'ĺ': "l", //slovenske
	'ľ': "l",
	'ô': "o",
	'ŕ': "r",
	'è': "e", // ani nevim
	'ã': "a",
	'ë': "e",
	'æ': "ae",
}
var tabulka2 map[string]string = map[string]string{
	"…": "...",
}
var mezery = regexp.MustCompile("( {2,})")

func UpravaTextu(txt string) string {
	var v strings.Builder
	for _, ch := range txt {
		if tabulka[ch] != "" {
			v.WriteString(tabulka[ch])
		} else {
			v.WriteRune(ch)
		}
	}

	for o, n := range tabulka2 {
		txt = strings.ReplaceAll(txt, o, n)
	}

	vysledek := v.String()
	vysledek = mezery.ReplaceAllString(vysledek, " ") // odstrani vice mezer

	return strings.TrimSpace(vysledek)
}
