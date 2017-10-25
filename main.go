package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/anyandrea/weather_app/lib/database"
	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/env"
	"github.com/urfave/negroni"
)

var (
	wdb weatherdb.WeatherDB
)

func main() {
	setupDatabase()

	// setup SIGINT catcher for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start a http server with negroni
	server := startHTTPServer()

	// wait for SIGINT
	<-stop
	log.Println("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}

func setupDatabase() {
	// setup weather database
	adapter := database.NewAdapter()
	adapter.RunMigrations("lib/database/migrations")
	wdb = weatherdb.NewWeatherDB(adapter)
}

func setupNegroni() *negroni.Negroni {
	n := negroni.Classic()

	r := newRouter()
	setupRoutes(r)
	n.UseHandler(r)

	return n
}

func startHTTPServer() *http.Server {
	addr := ":" + env.Get("PORT", "8080")
	server := &http.Server{Addr: addr, Handler: setupNegroni()}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}
