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

	"github.com/JamesClonk/home-info-dashboard/lib/config"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/env"
	"github.com/JamesClonk/home-info-dashboard/lib/forecasts"
	"github.com/JamesClonk/home-info-dashboard/lib/util"
	"github.com/JamesClonk/home-info-dashboard/lib/web/router"
	"github.com/urfave/negroni"
)

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
	spawnHousekeeping(hdb)
	spawnForecastCollection(hdb)

	return hdb
}

func spawnHousekeeping(hdb database.HomeInfoDB) {
	go func(hdb database.HomeInfoDB) {
		time.Sleep(24 * time.Hour) // initial waiting period

		for {
			// retention policy of 33 days and minimum 50'000 values
			if err := hdb.Housekeeping(33, 50000); err != nil {
				log.Println("Database housekeeping failed")
				log.Fatal(err)
			}
			time.Sleep(12 * time.Hour)
		}
	}(hdb)
}

func spawnForecastCollection(hdb database.HomeInfoDB) {
	go func(hdb database.HomeInfoDB) {
		time.Sleep(2 * time.Minute) // initial waiting period

		sensorId := config.Get().Forecast.TemperatureSensorID
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
					log.Println("Could not read temperature value for forecast temperature")
					log.Fatal(err)
				}

				if err := hdb.InsertSensorValue(sensorId, int(value), time.Now()); err != nil {
					log.Println("Could not insert temperature value for forecast temperature")
					log.Fatal(err)
				}
				log.Printf("Weather forecast temperature:%v for %s/%s stored to database\n", value, canton, city)
			}

			time.Sleep(1 * time.Hour)
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
