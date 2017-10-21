package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}

func setupRoutes(router *mux.Router) *mux.Router {
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/forecasts", ForecastsHandler)

	return router
}
