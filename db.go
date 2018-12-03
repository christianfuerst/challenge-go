package main

import (
	"database/sql"
	"log"
	"time"

	owm "github.com/briandowns/openweathermap"
)

// Initialisierung der Datenbank, falls noch nicht vorhanden
func prepareDb() {
	database, err := sql.Open("sqlite3", "./db/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS weather (city TEXT, time INT, temp FLOAT, humidity FLOAT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	database.Close()
}

// Wetter-Daten von openweathermap.org holen und in DB speichern
func getWeatherData(city string, apiKey string) {
	weatherData, err := owm.NewCurrent("C", "DE", apiKey)
	if err != nil {
		log.Fatalln(err)
	}

	weatherData.CurrentByName(city)

	database, err := sql.Open("sqlite3", "./db/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := database.Prepare("INSERT INTO weather (city, time, temp, humidity) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(city, time.Now(), weatherData.Main.Temp, weatherData.Main.Humidity)

	database.Close()
}

// Anfrage zu Wetter-Daten behandeln, in eine SQL-Abfrage überführen und Ergebnis zurück liefern
func queryWeatherDb(c string, t time.Time, m string) (w Weather, err error) {
	var temp float32
	var humidity float32
	var timeEnd time.Time

	database, err := sql.Open("sqlite3", "./db/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	if m == "byDay" {
		timeEnd = t.Add(time.Hour * 24)
	}

	if m == "byMonth" {
		timeEnd = endOfMonth(t)
	}

	row := database.QueryRow("select avg(temp) as temp, avg(humidity) as humidity from weather where city = ? and time between ? and ?", c, t, timeEnd)

	err = row.Scan(&temp, &humidity)
	if err != nil {
		log.Println("sql: Request returned no data.")
	}

	database.Close()

	w.City = c
	w.AvgTemp = temp
	w.AvgHum = humidity

	return w, err
}
