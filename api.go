package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tkanos/gonfig"
)

// WeatherEndPoint beinhaltet die Programm-Logik zur Auslieferung der angefragten Daten
func WeatherEndPoint(w http.ResponseWriter, r *http.Request) {
	var queryBy string
	var t time.Time

	defer r.Body.Close()

	// Query Parameter aus der URL extrahieren
	queryValues := r.URL.Query()

	// Lade Konfigurations-Datei
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	// Prüfung ob der Parameter 'city' einen gültigen Wert enthält
	if !stringInSlice(queryValues.Get("city"), configuration.Cities) {
		http.Error(w, "400 Bad Request: Parameter 'city' is missing or invalid.", http.StatusBadRequest)
		return
	}

	// Prüfung ob Parameter weder 'day' noch 'month' übergeben wurden
	if queryValues.Get("day") == "" && queryValues.Get("month") == "" {
		http.Error(w, "400 Bad Request: At least 'day' or 'month' parameter has to be used.", http.StatusBadRequest)
		return
	}

	// Prüfung ob versucht wurde Parameter 'day' und 'month' gleichzeitig zu übergeben
	if queryValues.Get("day") != "" && queryValues.Get("month") != "" {
		http.Error(w, "400 Bad Request: Either 'day' or 'month' parameter can be used.", http.StatusBadRequest)
		return
	}

	// Validierung des übergebenen Wertes von Parameter 'day'
	if queryValues.Get("day") != "" {
		queryBy = "byDay"

		layout := "2006-01-02"
		t, err = time.Parse(layout, queryValues.Get("day"))
		if err != nil {
			http.Error(w, "400 Bad Request: Parameter 'day' could not be processed.", http.StatusBadRequest)
			return
		}
	}

	// Validierung des übergebenen Wertes von Parameter 'month'
	if queryValues.Get("month") != "" {
		queryBy = "byMonth"

		layout := "2006-01"
		t, err = time.Parse(layout, queryValues.Get("month"))
		if err != nil {
			http.Error(w, "400 Bad Request: Parameter 'month' could not be processed.", http.StatusBadRequest)
			return
		}
	}

	// Gewünschte Daten aus der DB extrahieren
	msgData, err := queryWeatherDb(queryValues.Get("city"), t, queryBy)
	if err != nil {
		http.Error(w, "404 Not Found: Request returned no data.", http.StatusNotFound)
		return
	}

	// SQL Ergebnis in ein JSON Objekt überführen und anschließend ausliefern
	msg, err := json.Marshal(msgData)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, string(msg))
}
