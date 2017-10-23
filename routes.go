package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}

func setupRoutes(router *mux.Router) *mux.Router {
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/error", ErrorHandler)

	router.HandleFunc("/forecasts", ForecastsHandler)
	router.HandleFunc("/forecasts/{canton}", ForecastsHandler)
	router.HandleFunc("/forecasts/{canton}/{city}", ForecastsHandler)

	return router
}
