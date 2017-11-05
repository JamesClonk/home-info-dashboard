package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/anyandrea/weather_app/lib/database"
	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/env"
	"github.com/anyandrea/weather_app/lib/web/router"
	"github.com/urfave/negroni"
)

func main() {
	env.MustGet("WEATHERAPI_PASSWORD")
	wdb := setupDatabase()

	// setup SIGINT catcher for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start a http server with negroni
	server := startHTTPServer(wdb)

	// wait for SIGINT
	<-stop
	log.Println("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}

func setupDatabase() weatherdb.WeatherDB {
	// setup weather database
	adapter := database.NewAdapter()
	if err := adapter.RunMigrations("lib/database/migrations"); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			log.Println("Could not run database migrations")
			log.Fatal(err)
		}
	}
	wdb := weatherdb.NewWeatherDB(adapter)

	// generate fake data
	if env.Get("WEATHERDB_MOCK_DATA", "false") == "true" {
		sensors, err := wdb.GetSensors()
		if err != nil {
			log.Println("Could not get sensors from database")
			log.Fatal(err)
		}
		for _, sensor := range sensors {
			wdb.DropSensorValues(sensor.Id)
			wdb.GenerateSensorValues(sensor.Id, 50)
		}
	}

	return wdb
}

func setupNegroni(wdb weatherdb.WeatherDB) *negroni.Negroni {
	n := negroni.Classic()

	r := router.New(wdb)
	n.UseHandler(r)

	return n
}

func startHTTPServer(wdb weatherdb.WeatherDB) *http.Server {
	addr := ":" + env.Get("PORT", "8080")
	server := &http.Server{Addr: addr, Handler: setupNegroni(wdb)}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}
