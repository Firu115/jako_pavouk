package main

import (
	"backend/databaze"
	"backend/utils"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func setupRouter(app *fiber.App) {
	app.Get("/lekce", getVsechnyLekce)
	app.Get("/lekce/:pismena", getCviceniVLekci)
	app.Get("/cvic/:pismena/:cislo", getCviceni)
	app.Post("/registrace", registrace)
	app.Post("/prihlaseni", prihlaseni)
	app.Get("/ja", prehled)
	app.Get("/test", test)
}

func test(c *fiber.Ctx) error {
	return c.JSON("")
}

func getVsechnyLekce(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}
	lekce, err := databaze.GetLekce()
	if err != nil {
		return err
	}
	if id != 0 {
		return c.Status(http.StatusOK).JSON(fiber.Map{"lekce": lekce, "uziv": id}) //TODO
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"lekce": lekce, "uziv": id})
}

func getCviceniVLekci(c *fiber.Ctx) error {
	_, err := utils.Autentizace(c, false)
	if err != nil {
		return err
	}
	cvic, err := databaze.GetCviceniVLekciByPismena(c.Params("pismena"))
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"cviceni": cvic})
}

func getCviceni(c *fiber.Ctx) error {
	pismena := c.Params("pismena")
	vsechnyCviceni, err := databaze.GetCviceniVLekciByPismena(pismena)
	if err != nil {
		return err
	}
	cislo, err := strconv.Atoi(c.Params("cislo")) // str -> int
	if err != nil {
		return err
	}
	if cislo > len(vsechnyCviceni) {
		return c.Status(http.StatusBadRequest).JSON("Cviceni neexistuje")
	}
	if vsechnyCviceni[cislo-1].Typ == "pismena" {
		var text string = ""
		for i := 0; i < pocetSlov; i++ {
			var slovo string = ""
			for j := 0; j < pocetPismenVeSlovu; j++ {
				slovo += string(pismena[rand.Intn(len(pismena))])
			}
			slovo += " "
			text += slovo
		}
		text = text[0 : len(text)-1] // smazat mezeru na konci
		return c.Status(http.StatusOK).JSON(fiber.Map{"text": text})
	} else {
		return c.Status(http.StatusOK).JSON(fiber.Map{"text": "dsf"})
	}
}

func registrace(c *fiber.Ctx) error {
	// overeni spravnych dat co prijdou
	var body struct {
		Email string `json:"email" binding:"required"`
		Jmeno string `json:"jmeno" binding:"required,min=3"`
		Heslo string `json:"heslo" binding:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	// overeni jestli nahodou uz neexistuje
	_, err := databaze.GetUzivByEmail(body.Email)
	if err != nil {
		hesloHASH, err := utils.HashPassword(body.Heslo)
		if err != nil {
			return err
		}
		id, err := databaze.CreateUziv(body.Email, hesloHASH, body.Jmeno)
		if err != nil {
			return err
		}
		token, err := utils.GenerovatToken(body.Email, id)
		if err != nil {
			return err
		} else {
			return c.Status(http.StatusOK).JSON(fiber.Map{"Token": token})
		}
	} else {
		return c.Status(http.StatusBadRequest).JSON("Uzivatel jiz existuje")
	}
}

func prihlaseni(c *fiber.Ctx) error {
	// validace body dat
	var body struct {
		Email string `json:"email" binding:"required"`
		Heslo string `json:"heslo" binding:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	// validace emailu
	if !utils.ValidFormat(body.Email) {
		return c.Status(http.StatusBadRequest).JSON("Invalidni email")
	}
	// pull z databaze
	uziv, err := databaze.GetUzivByEmail(body.Email)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Uzivatel neexistuje")
	}

	if err := utils.CheckPassword(body.Heslo, uziv.Heslo); err != nil {
		return c.Status(http.StatusUnauthorized).JSON("Spatne heslo")
	} else {
		token, err := utils.GenerovatToken(uziv.Email, uziv.ID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON("Token se pokazil")
		} else {
			return c.Status(http.StatusOK).JSON(fiber.Map{"Token": token})
		}
	}
}

func prehled(c *fiber.Ctx) error {
	id, err := utils.Autentizace(c, true)
	if err != nil {
		return err
	}
	uziv, err := databaze.GetUzivByID(id)
	if err != nil {
		return err
	}
	prumerPreklepu := utils.Prumer(uziv.Preklepy)
	prumerRychlosti := utils.Prumer(uziv.Rychlosti)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"email":           uziv.Email,
		"jmeno":           uziv.Jmeno,
		"prumerPreklepu":  prumerPreklepu,
		"prumerRychlosti": prumerRychlosti,
	})
}
