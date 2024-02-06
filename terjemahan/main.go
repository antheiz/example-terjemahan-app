package main 

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	http.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		source := r.FormValue("source")
		target := r.FormValue("target")
		text := r.FormValue("text")
		result, err := Translate(dictionary, source, target, text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `{"translation": "%s"}`, result)
	})

	http.ListenAndServe(":8080", nil)
}

func Translate(dictionary Dictionary, source, target, text string) (string, error) {
	translation := ""
	// var ok bool

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
