package databaze

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
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Jmeno string `json:"jmeno"`
	Heslo string `json:"heslo"`
}

type Slovnik struct {
	ID    uint   `json:"id"`
	Slovo string `json:"slovo"`
}

type Dokoncene struct {
	ID        uint    `json:"id"`
	UzivID    uint    `json:"uziv_id"`
	CviceniID uint    `json:"cviceni_id"`
	CPM       float32 `json:"cpm"`
	Preklepy  uint    `json:"preklepy"`
}

func GetLekce() ([][]Lekce, error) {
	var lekce [][]Lekce = [][]Lekce{}

	rows, err := DB.Query(`SELECT * FROM lekce;`)
	if err != nil {
		return lekce, nil
	}
	defer rows.Close()

	var skupina []Lekce

	var cisloSkupiny uint = 1
	var jednaLekce Lekce
	for rows.Next() {
		jednaLekce = Lekce{}
		var skup uint

		err := rows.Scan(&jednaLekce.ID, &jednaLekce.Pismena, &skup)
		if err != nil {
			return lekce, err
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
	var vysledek []int32 = []int32{}
	// zjistim kolik ma kazda lekce cviceni
	rows, err := DB.Query(`SELECT id FROM lekce l WHERE 0 = (SELECT COUNT(*) FROM cviceni c WHERE c.lekce_id = l.id) - (SELECT COUNT(*) FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id AND d.uziv_id = $1 AND l.id = c.lekce_id);`, uzivID)
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return vysledek, err
		}
		vysledek = append(vysledek, int32(id))
	}

	return vysledek, nil
}

func GetDokonceneCvicVLekci(uzivID uint, lekceID uint) ([]int32, error) {
	var cviceni_ids []int32 = []int32{}

	rows, err := DB.Query(`SELECT d.cviceni_id FROM dokoncene d JOIN cviceni c ON d.cviceni_id = c.id WHERE lekce_id = $1 AND uziv_id = $2;`, lekceID, uzivID)
	if err != nil {
		return cviceni_ids, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return cviceni_ids, err
		}
		cviceni_ids = append(cviceni_ids, id)
	}

	return cviceni_ids, nil
}

func GetLekceIDbyPismena(pismena string) (uint, error) {
	var id uint
	err := DB.QueryRow(`SELECT id FROM lekce WHERE pismena = $1;`, pismena).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCviceniVLekciByID(lekceID uint) ([]Cviceni, error) {
	var cviceni []Cviceni

	rows, err := DB.Query(`SELECT id, typ FROM cviceni WHERE lekce_id = $1;`, lekceID)
	if err != nil {
		return cviceni, err
	}
	defer rows.Close()

	for rows.Next() {
		jednoCviceni := Cviceni{}
		err := rows.Scan(&jednoCviceni.ID, &jednoCviceni.Typ)
		if err != nil {
			return cviceni, err
		}

		cviceni = append(cviceni, jednoCviceni)
	}

	return cviceni, nil
}

func GetCviceniVLekciByPismena(pismena string) ([]Cviceni, error) {
	var cviceni []Cviceni
	id, err := GetLekceIDbyPismena(pismena)
	if err != nil {
		return cviceni, err
	}
	cviceni, err = GetCviceniVLekciByID(id)
	return cviceni, err
}

func GetUzivByID(id uint) (Uzivatel, error) {
	var uziv Uzivatel
	err := DB.QueryRow(`SELECT id, email, jmeno, heslo FROM uzivatel WHERE id = $1;`, id).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo)
	return uziv, err
}

func GetUzivByEmail(email string) (Uzivatel, error) {
	var uziv = Uzivatel{}
	err := DB.QueryRow(`SELECT id, email, jmeno, heslo FROM uzivatel WHERE email = $1;`, email).Scan(&uziv.ID, &uziv.Email, &uziv.Jmeno, &uziv.Heslo)
	return uziv, err
}

func GetPreklepyACPM(uzivID uint) ([]float32, []float32, error) {
	var preklepy []float32
	var cpm []float32

	var posledni int = 10
	rows, err := DB.Query(`SELECT preklepy, cpm FROM dokoncene WHERE uziv_id = $1 ORDER BY id DESC LIMIT $2;`, uzivID, posledni)
	if err != nil {
		return preklepy, cpm, err
	}
	defer rows.Close()

	for rows.Next() {
		var preklep float32
		var cpmko float32
		err := rows.Scan(&preklep, &cpmko)
		if err != nil {
			return preklepy, cpm, err
		}

		preklepy = append(preklepy, float32(preklep))
		cpm = append(cpm, cpmko)
	}

	return preklepy, cpm, nil
}

func DokonceneProcento(uzivID uint) (float32, error) {
	var pocet int32
	err := DB.QueryRow(`SELECT COUNT(*) FROM dokoncene WHERE uziv_id = $1;`, uzivID).Scan(&pocet)
	if err != nil {
		return 0, err
	}

	var pocet2 int32
	err = DB.QueryRow(`SELECT COUNT(*) FROM cviceni;`).Scan(&pocet2)
	if err != nil {
		return 0, err
	}
	return float32(pocet) / float32(pocet2) * 100, nil
}

func CreateUziv(email string, hesloHash string, jmeno string) (uint, error) {
	_, err := DB.Exec(`INSERT INTO uzivatel (email, jmeno, heslo) VALUES ($1, $2, $3)`, email, jmeno, hesloHash)
	if err != nil {
		return 0, err
	}
	uziv, err := GetUzivByEmail(email)
	if err != nil {
		return 0, err
	}
	return uziv.ID, nil
}

func PridatDokonceneCvic(cvicID uint, uzivID uint, cpm float32, preklepy int) error {
	_, err := DB.Exec(`INSERT INTO dokoncene (uziv_id, cviceni_id, cpm, preklepy) VALUES ($1, $2, $3, $4);`, uzivID, cvicID, cpm, preklepy)
	return err
}

func OdebratDokonceneCvic(cvicID uint, uzivID uint) error {
	_, err := DB.Exec(`DELETE FROM dokoncene WHERE uziv_id = $1 AND cviceni_id = $2;`, uzivID, cvicID)
	return err
}

func GetSlovaProLekci(lekceID uint, vsechny bool) ([]string, error) {
	var vysledek []string
	rows, err := DB.Query(`SELECT slovo FROM slovnik WHERE lekce_id = $1;`, lekceID)
	if err != nil {
		return vysledek, err
	}
	defer rows.Close()

	var pismenaJedny string
	for rows.Next() {
		pismenaJedny = ""
		err := rows.Scan(&pismenaJedny)
		if err != nil {
			return vysledek, err
		}

		vysledek = append(vysledek, pismenaJedny)
	}

	return vysledek, nil
}