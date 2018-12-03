package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"

	_ "github.com/mattn/go-sqlite3"
)

// Configuration repräsentiert die Konfigurations-Datei config.json
type Configuration struct {
	APIKey        string
	Cities        []string
	ErrorLogFile  string
	AccessLogFile string
	Port          string
}

// Weather repräsentiert die Rückgabewerte des /weather Endpunktes der API
type Weather struct {
	City    string  `json:"city"`
	AvgTemp float32 `json:"avg_temp"`
	AvgHum  float32 `json:"avg_hum"`
}

func main() {

	// Lade Konfigurations-Datei
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	// Log-Dateien initialisieren
	errorFile, err := os.OpenFile(configuration.ErrorLogFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(errorFile)
	defer errorFile.Close()

	accessFile, err := os.OpenFile(configuration.AccessLogFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer accessFile.Close()

	// DB-Struktur erstellen, falls nötig
	prepareDb()

	// API-Endpunkt /weather konfigurieren
	r := mux.NewRouter()
	r.HandleFunc("/weather", WeatherEndPoint).Methods("GET")
	loggedRouter := handlers.LoggingHandler(accessFile, r)

	// http Server im Hintergrund starten
	go func() {
		if err := http.ListenAndServe(":"+configuration.Port, loggedRouter); err != nil {
			log.Fatal(err)
		}
	}()

	// Wetterdaten von openweathermap.org einmal pro Stunde abfragen und in DB speichern
	t := time.NewTicker(time.Hour)
	for {
		for _, element := range configuration.Cities {
			getWeatherData(element, configuration.APIKey)
		}
		<-t.C
	}

}
