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

func startHTTPServer() *http.Server {
	handler := negroni.Classic()

	router := newRouter()
	setupRoutes(router)
	handler.UseHandler(router)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	addr := ":" + port
	server := &http.Server{Addr: addr, Handler: handler}

	go func() {
		log.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}
