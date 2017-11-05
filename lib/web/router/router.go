package router

import (
	"net/http"

	auth "github.com/abbot/go-http-auth"
	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/util"
	"github.com/anyandrea/weather_app/lib/web/api"
	"github.com/anyandrea/weather_app/lib/web/html"
	"github.com/gorilla/mux"
)

func New(wdb weatherdb.WeatherDB) *mux.Router {
	router := mux.NewRouter()
	setupRoutes(wdb, router)
	return router
}

func setupRoutes(wdb weatherdb.WeatherDB, router *mux.Router) *mux.Router {
	// HTML
	router.NotFoundHandler = http.HandlerFunc(html.NotFound)

	router.HandleFunc("/", html.Index(wdb))
	router.HandleFunc("/error", html.ErrorHandler)

	router.HandleFunc("/sensors", html.Sensors(wdb))

	router.HandleFunc("/forecasts", html.Forecasts)
	router.HandleFunc("/forecasts/{canton}", html.Forecasts)
	router.HandleFunc("/forecasts/{canton}/{city}", html.Forecasts)

	// API
	router.HandleFunc("/sensor", api.Sensors(wdb)).Methods("GET")
	router.HandleFunc("/sensor/all", api.Sensors(wdb)).Methods("GET")
	router.HandleFunc("/sensor/list", api.Sensors(wdb)).Methods("GET")
	router.HandleFunc("/sensor/{id}", api.Sensor(wdb)).Methods("GET")

	router.HandleFunc("/sensor/{id}/values", api.SensorValues(wdb)).Methods("GET")
	router.HandleFunc("/sensor/{id}/values/{limit}", api.SensorValues(wdb)).Methods("GET")

	// secured API
	authenticator := auth.NewBasicAuthenticator("weatherapp", secret)
	router.HandleFunc("/sensor/{id}/value", auth.JustCheck(authenticator, api.AddSensorValue(wdb))).Methods("POST")

	return router
}

func secret(user, realm string) string {
	username, password := util.GetUserAndPassword()
	if user == username {
		return password
	}
	return ""
}
