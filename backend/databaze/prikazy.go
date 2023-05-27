package databaze

import (
	"log"

	"github.com/lib/pq"
)

type Lekce struct {
	ID      uint
	Pismena string
}

type Cviceni struct {
	ID  uint
	Typ string
}

type Uzivatel struct {
	ID        uint
	Email     string
	Jmeno     string
	Heslo     string
	Preklepy  pq.Int32Array // takhle se dela array v sql
	Rychlosti pq.Int32Array
}

type Slovnik struct {
	ID    uint
	Slovo string
}

func GetLekce() ([]Lekce, error) {
	rows, err := DB.Query("SELECT * FROM lekce;")
	if err != nil {
		return []Lekce{}, nil
	}
	defer rows.Close()

	var lekce []Lekce
	for rows.Next() {
		jednaLekce := Lekce{}
		err := rows.Scan(&jednaLekce.ID, &jednaLekce.Pismena)
		if err != nil {
			log.Fatal(err)
		}

		lekce = append(lekce, jednaLekce)
	}

	return lekce, nil
}

func GetLekceIDbyPismena(pismena string) (uint, error) {
	var id uint
	err := DB.QueryRow("SELECT id FROM lekce WHERE pismena = $1;", pismena).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCviceniVLekciByID(lekceID uint) ([]Cviceni, error) {
	rows, err := DB.Query("SELECT id, typ FROM cviceni WHERE lekce_id = $1;", lekceID)
	if err != nil {
		return []Cviceni{}, err
	}
	defer rows.Close()

	var cviceni []Cviceni
	for rows.Next() {
		jednoCviceni := Cviceni{}
		err := rows.Scan(&jednoCviceni.ID, &jednoCviceni.Typ)
		if err != nil {
			log.Fatal(err)
		}

		cviceni = append(cviceni, jednoCviceni)
	}

	return cviceni, nil
}

func GetCviceniVLekciByPismena(pismena string) ([]Cviceni, error) {
	cviceni := []Cviceni{}
	id, err := GetLekceIDbyPismena(pismena)
	if err != nil {
		return cviceni, err
	}
	cviceni, err = GetCviceniVLekciByID(id)
	return cviceni, err
}

func GetUzivByID(id uint) (Uzivatel, error) {
	var uziv Uzivatel
	err := DB.QueryRow("SELECT id, email, jmeno, heslo, preklepy, rychlosti FROM uzivatele WHERE id = $1;", id).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo, &uziv.Preklepy, &uziv.Rychlosti)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv = Uzivatel{}
	err := DB.QueryRow("SELECT id, email, jmeno, heslo, preklepy, rychlosti FROM uzivatele WHERE email = $1;", email).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo, &uziv.Preklepy, &uziv.Rychlosti)
	return uziv, err
}

func CreateUziv(email string, hesloHash string, jmeno string) (uint, error) {
	_, err := DB.Exec("INSERT INTO uzivatele (email, jmeno, heslo) VALUES ($1, $2, $3)", email, jmeno, hesloHash)
	if err != nil {
		return 0, err
	}
	uziv, err := GetUzivByEmail(email)
	if err != nil {
		return 0, err
	}
	return uziv.ID, nil
}
