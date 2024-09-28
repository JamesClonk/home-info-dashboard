package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/alerting"
	"github.com/JamesClonk/home-info-dashboard/lib/config"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/env"
	"github.com/JamesClonk/home-info-dashboard/lib/forecasts"
	"github.com/JamesClonk/home-info-dashboard/lib/util"
	"github.com/JamesClonk/home-info-dashboard/lib/web/router"
	"github.com/urfave/negroni"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	env.MustGet("AUTH_PASSWORD")
	hdb := setupDatabase()

	// setup SIGINT catcher for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start a http server with negroni
	server := startHTTPServer(hdb)

	// wait for SIGINT
	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = server.Shutdown(ctx)
	log.Println("Server gracefully stopped")
	cancel()
}

func setupDatabase() database.HomeInfoDB {
	// setup weather database
	adapter := database.NewAdapter()
	if err := adapter.RunMigrations("lib/database/migrations"); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			log.Println("Could not run database migrations")
			log.Fatal(err)
		}
	}
	hdb := database.NewHomeInfoDB(adapter)

	// background jobs
	spawnAlertMonitor(hdb)
	spawnHousekeeping(hdb)
	spawnForecastCollection(hdb)

	return hdb
}

func spawnAlertMonitor(hdb database.HomeInfoDB) {
	alerting.Init(hdb)

	go func() {
		time.Sleep(5 * time.Second)
		alerting.RegisterAlerts()
	}()
}

func spawnHousekeeping(hdb database.HomeInfoDB) {
	go func(hdb database.HomeInfoDB) {
		time.Sleep(time.Duration(rand.Intn(60)) * time.Minute)
		time.Sleep(24 * time.Hour) // initial waiting period

		for {
			if err := hdb.Housekeeping(); err != nil {
				log.Println("Database housekeeping failed")
				log.Fatal(err)
			}
			time.Sleep(time.Duration(rand.Intn(30)) * time.Minute)
			time.Sleep(12 * time.Hour)
		}
	}(hdb)
}

func spawnForecastCollection(hdb database.HomeInfoDB) {
	go func(hdb database.HomeInfoDB) {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Minute)
		time.Sleep(2 * time.Minute) // initial waiting period

		tempSensorId := config.Get().Forecast.TemperatureSensorID
		humSensorId := config.Get().Forecast.HumiditySensorID
		windSensorId := config.Get().Forecast.WindSpeedSensorID
		for {
			lat, lon, alt := util.GetDefaultLocation("", "", "")
			forecast, err := forecasts.Get(lat, lon, alt)
			if err != nil {
				log.Println("Weather forecast collection failed")
				log.Fatal(err)
			}

			if len(forecast.Properties.Timeseries) > 0 {
				var temp, hum, wind float64

				// try to get an entry for "current" hour
				var foundCurrent bool
				for _, entry := range forecast.Properties.Timeseries {
					now := time.Now()
					if entry.Time.Local().Day() == now.Local().Day() && entry.Time.Local().Hour() == now.Local().Hour() {
						temp = entry.Data.Instant.Details.AirTemperature
						hum = entry.Data.Instant.Details.RelativeHumidity
						wind = entry.Data.Instant.Details.WindSpeed
						foundCurrent = true
						break
					}
				}
				// else fallback to first entry
				if !foundCurrent {
					temp = forecast.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
					hum = forecast.Properties.Timeseries[0].Data.Instant.Details.RelativeHumidity
					wind = forecast.Properties.Timeseries[0].Data.Instant.Details.WindSpeed
				}

				// store temp
				if err := hdb.InsertSensorValue(tempSensorId, int(temp), time.Now()); err != nil {
					log.Println("Could not insert temperature value for forecast")
					log.Fatal(err)
				}
				log.Printf("Weather forecast temperature:%v for [%s/%s/%s] stored to database\n", temp, lat, lon, alt)
				// store humidity
				if err := hdb.InsertSensorValue(humSensorId, int(hum), time.Now()); err != nil {
					log.Println("Could not insert humidity value for forecast")
					log.Fatal(err)
				}
				log.Printf("Weather forecast humidity:%v for [%s/%s/%s] stored to database\n", hum, lat, lon, alt)
				// store wind speed
				if err := hdb.InsertSensorValue(windSensorId, int(wind), time.Now()); err != nil {
					log.Println("Could not insert wind speed value for forecast")
					log.Fatal(err)
				}
				log.Printf("Weather forecast wind speed:%v for [%s/%s/%s] stored to database\n", wind, lat, lon, alt)
			}

			time.Sleep(time.Duration(rand.Intn(5)) * time.Minute)
			time.Sleep(15 * time.Minute)
		}
	}(hdb)
}

func setupNegroni(hdb database.HomeInfoDB) *negroni.Negroni {
	n := negroni.Classic()

	r := router.New(hdb)
	n.UseHandler(r)

	return n
}

func startHTTPServer(hdb database.HomeInfoDB) *http.Server {
	addr := ":" + env.Get("PORT", "8080")
	server := &http.Server{Addr: addr, Handler: setupNegroni(hdb)}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}
