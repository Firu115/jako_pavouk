package databaze

import (
	"log"

	"github.com/lib/pq"
)

type Lekce struct {
	ID      uint   `json:"id"`
	Pismena string `json:"pismena"`
	// Skupina uint        nepouzivame ale je tam
}

type Cviceni struct {
	ID  uint   `json:"id"`
	Typ string `json:"typ"`
}

type Uzivatel struct {
	ID               uint          `json:"id"`
	Email            string        `json:"email"`
	Jmeno            string        `json:"jmeno"`
	Heslo            string        `json:"heslo"`
	Preklepy         pq.Int32Array `json:"preklepy"` // takhle se dela array z sql wtf
	Rychlosti        pq.Int32Array `json:"rychlost"`
	DokonceneCviceni pq.Int32Array `json:"dokoncenecviceni"`
}

type Slovnik struct {
	ID    uint   `json:"id"`
	Slovo string `json:"slovo"`
}

func GetLekce() ([][]Lekce, error) {
	rows, err := DB.Query("SELECT * FROM lekce;")
	if err != nil {
		return [][]Lekce{}, nil
	}
	defer rows.Close()

	var lekce [][]Lekce
	var skupina []Lekce

	var cisloSkupiny uint = 1
	for rows.Next() {
		jednaLekce := Lekce{}
		var skup uint
		err := rows.Scan(&jednaLekce.ID, &jednaLekce.Pismena, &skup)
		if err != nil {
			log.Fatal(err)
		}
		if cisloSkupiny == skup {
			skupina = append(skupina, jednaLekce)
		} else if cisloSkupiny < skup {
			lekce = append(lekce, skupina)
			skupina = []Lekce{jednaLekce}
			cisloSkupiny += 1
		} else {
			lekce[skup] = append(lekce[skup], jednaLekce)
		}
	}
	lekce = append(lekce, skupina)

	return lekce, nil
}

func GetDokonceneLekce(uzivID uint) ([]int32, error) {
	var doko pq.Int32Array
	err := DB.QueryRow("SELECT dokoncenecviceni FROM uzivatele WHERE id = $1;", uzivID).Scan(&doko)
	if err != nil {
		return []int32{}, err
	}
	return doko, nil //ted mam dokonceny cviceni potrebuju zjistit jestli to je vsechno v ty lekce TODO
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
	err := DB.QueryRow("SELECT id, email, jmeno, heslo, preklepy, rychlosti, dokoncenecviceni FROM uzivatele WHERE id = $1;", id).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo, &uziv.Preklepy, &uziv.Rychlosti, &uziv.DokonceneCviceni)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv = Uzivatel{}
	err := DB.QueryRow("SELECT id, email, jmeno, heslo, preklepy, rychlosti, dokoncenecviceni FROM uzivatele WHERE email = $1;", email).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo, &uziv.Preklepy, &uziv.Rychlosti, &uziv.DokonceneCviceni)
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
