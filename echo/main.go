package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

	e := echo.New()

	e.POST("/translate", func(c echo.Context) error {
		source := c.FormValue("source")
		target := c.FormValue("target")
		text := c.FormValue("text")

		result, err := Translate(dictionary, source, target, text)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"translation": result})
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			return next(c)
		}
	})

	e.Logger.Fatal(e.Start(":8080"))
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
