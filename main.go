package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/anyandrea/weather_app/lib/database"
	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/env"
	"github.com/anyandrea/weather_app/lib/forecasts"
	"github.com/anyandrea/weather_app/lib/util"
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
			wdb.DeleteSensorValues(sensor.Id)
			wdb.GenerateSensorValues(sensor.Id, 50)
		}
	}

	// background jobs
	spawnHousekeeping(wdb)
	spawnForecastCollection(wdb)

	return wdb
}

func spawnHousekeeping(wdb weatherdb.WeatherDB) {
	go func(wdb weatherdb.WeatherDB) {
		for {
			// retention policy of 33 days and minimum 50'000 values
			if err := wdb.Housekeeping(33, 50000); err != nil {
				log.Println("Database housekeeping failed")
				log.Fatal(err)
			}
			time.Sleep(12 * time.Hour)
		}
	}(wdb)
}

func spawnForecastCollection(wdb weatherdb.WeatherDB) {
	go func(wdb weatherdb.WeatherDB) {
		sensor, err := wdb.GetSensorByName("forecast_temperature")
		if err != nil {
			log.Println("Could not select sensor [forecast_temperature]")
			log.Fatal(err)
		}

		for {
			canton, city := util.GetDefaultLocation("", "")
			forecast, err := forecasts.Get(canton, city)
			if err != nil {
				log.Println("Weather forecast collection failed")
				log.Fatal(err)
			}

			if len(forecast.Forecast.Tabular.Time) > 0 {
				value, err := strconv.ParseInt(forecast.Forecast.Tabular.Time[0].Temperature.Value, 10, 64)
				if err != nil {
					log.Println("Could not read temperature value for [forecast_temperature]")
					log.Fatal(err)
				}

				if err := wdb.InsertSensorValue(sensor.Id, int(value), time.Now()); err != nil {
					log.Println("Could not insert temperature value for [forecast_temperature]")
					log.Fatal(err)
				}
				log.Printf("Weather forecast [forecast_temperature:%v] for [%s/%s] stored to database\n", value, canton, city)
			}

			time.Sleep(1 * time.Hour)
		}
	}(wdb)
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
