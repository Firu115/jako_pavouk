package main

import (
	"bufio"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"scripts/utils"
	"strings"
	"unicode/utf8"

	"github.com/blockloop/scan"
)

var htmlRegex *regexp.Regexp = regexp.MustCompile(`<.*?>`)
var spaceRegex *regexp.Regexp = regexp.MustCompile(`(\r\n|\r|\n)|\t`)
var jmenoSouboru string = "clanky.txt"

type Rss struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"encoded"` // content:encoded <- NAMESPACE:TAG
}

type DobraZprava struct {
	ID            uint         `db:"id"`
	Jmeno         string       `db:"jmeno"`
	Txt           string       `db:"txt"`
	DatumZarazeni sql.NullTime `db:"datum_zarazeni"`
	Delka         int          `db:"delka"`
}

func main() {
	newTexts := parseRss().Channel.Items

	f, err := os.Open(jmenoSouboru)
	if err != nil {
		f, err = os.Create(jmenoSouboru)
		if err != nil {
			chyba(err)
		}
	}

	scanner := bufio.NewScanner(f)

	var texty map[string]string = make(map[string]string)
	var jmeno string
	for scanner.Scan() {
		if jmeno == "" {
			jmeno = strings.TrimSpace(scanner.Text())
		} else if scanner.Text() == "" {
			texty[jmeno] = strings.TrimSpace(texty[jmeno])
			jmeno = ""
		} else {
			texty[jmeno] += strings.TrimSpace(scanner.Text())
			texty[jmeno] += " "
		}
	}
	f.Close()

	utils.DBConnect()

	var textyZdb []DobraZprava
	rows, err := utils.DB.Query(`SELECT * FROM dobry_zpravy;`)
	if err != nil {
		chyba(err)
	}
	err = scan.Rows(&textyZdb, rows)
	if err != nil {
		chyba(err)
	}

	for _, new := range newTexts {
		_, err := utils.DB.Exec(`INSERT INTO dobry_zpravy (jmeno, txt, datum_zarazeni, delka) VALUES ($1, $2, NULL, $3);`, new.Title, new.Content, utf8.RuneCountInString(new.Content))
		fmt.Println(err)
	}

	/* 	f, err = os.OpenFile(jmenoSouboru, os.O_APPEND|os.O_WRONLY, 0644)
	   	if err != nil {
	   		chyba(err)
	   	}

	   	for _, new := range rss.Channel.Items {
	   		if texty[new.Title] == "" {
	   			if _, err = f.WriteString(fmt.Sprintf("%s\n%s\n\n", new.Title, new.Content)); err != nil {
	   				chyba(err)
	   			}
	   		}
	   	} */
}

func smazatHtml(s string) string {
	s = htmlRegex.ReplaceAllString(s, "")
	s = spaceRegex.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	return s
}

func oriznout(s string) string {
	konec := regexp.MustCompile(`< *hr.*?>|Autor:|(Máte rádi naše dobrý zprávy a chtěli byste nás nějak podpořit?)`).FindStringIndex(s)
	if len(konec) != 0 {
		s = s[:konec[0]]
	}
	return s
}

func getXML() []byte {
	resp, err := http.Get("https://dobryzpravy.cz/feed/")
	if err != nil {
		chyba(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		chyba(errors.New(fmt.Sprintf("Status code: %d", resp.StatusCode)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		chyba(err)
	}

	return data
}

func parseRss() Rss {
	var rss Rss
	err := xml.Unmarshal(getXML(), &rss)
	if err != nil {
		chyba(err)
	}

	for i := 0; i < len(rss.Channel.Items); i++ {
		str := rss.Channel.Items[i].Content
		str = smazatHtml(str)
		str = utils.UpravaTextu(str)
		str = oriznout(str)

		rss.Channel.Items[i].Content = str

		rss.Channel.Items[i].Title = utils.UpravaTextu(rss.Channel.Items[i].Title)
		//fmt.Printf("%s\n\n\n\n", rss.Channel.Items[i].Content)
	}
	return rss
}

func chyba(err error) {
	panic(err)
}
