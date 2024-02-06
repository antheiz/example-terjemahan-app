package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Dictionary struct {
	EnToId map[string]string `json:"moy-id"`
	IdToEn map[string]string `json:"id-moy"`
}

func main() {
	dictionary := Dictionary{}

	data, err := ioutil.ReadFile("../dictionary.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &dictionary)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		return c.Next()
	})

	app.Post("/translate", func(c *fiber.Ctx) error {
		source := c.FormValue("source")
		target := c.FormValue("target")
		text := c.FormValue("text")

		result, err := Translate(dictionary, source, target, text)
		if err != nil {
			return c.JSON(map[string]string{"error": err.Error()})
		}
		return c.JSON(map[string]string{"translation": result})
	})

	log.Fatal(app.Listen(":8080"))
}

func Translate(dictionary Dictionary, source, target, text string) (string, error) {
	translation := ""

	if source == "moy" && target == "id" {
		words := strings.Split(text, " ")
		for _, word := range words {
			translated, found := dictionary.EnToId[strings.ToLower(word)]
			if found {
				translation += translated + " "
			} else {
				translation += word + " "
			}
		}
	} else if source == "id" && target == "moy" {
		words := strings.Split(text, " ")
		for _, word := range words {
			translated, found := dictionary.IdToEn[strings.ToLower(word)]
			if found {
				translation += translated + " "
			} else {
				translation += word + " "
			}
		}
	} else {
		return "", fmt.Errorf("Mohon maaf, Kata/Bahasa tersebut belum tersedia")
	}

	return strings.TrimSpace(translation), nil
}
