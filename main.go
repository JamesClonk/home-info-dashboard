package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/negroni"
)

func main() {
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

func setupNegroni() *negroni.Negroni {
	n := negroni.Classic()

	r := newRouter()
	setupRoutes(r)
	n.UseHandler(r)

	return n
}

func startHTTPServer() *http.Server {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	addr := ":" + port
	server := &http.Server{Addr: addr, Handler: setupNegroni()}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}
