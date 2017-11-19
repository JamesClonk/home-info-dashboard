package router

import (
	"net/http"

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

	router.HandleFunc("/dashboard", html.Dashboard(wdb))
	router.HandleFunc("/graphs", html.Graphs(wdb))
	router.HandleFunc("/sensor_data", html.Sensors(wdb))

	router.HandleFunc("/forecasts", html.Forecasts)
	router.HandleFunc("/forecasts/{canton}", html.Forecasts)
	router.HandleFunc("/forecasts/{canton}/{city}", html.Forecasts)

	// API
	router.HandleFunc("/sensor_type", api.GetSensorTypes(wdb)).Methods("GET")
	router.HandleFunc("/sensor_types", api.GetSensorTypes(wdb)).Methods("GET")
	router.HandleFunc("/sensor_type/{id}", api.GetSensorType(wdb)).Methods("GET")

	router.HandleFunc("/sensor", api.GetSensors(wdb)).Methods("GET")
	router.HandleFunc("/sensors", api.GetSensors(wdb)).Methods("GET")
	router.HandleFunc("/sensor/{id}", api.GetSensor(wdb)).Methods("GET")

	router.HandleFunc("/sensor/{id}/values", api.GetSensorValues(wdb)).Methods("GET")
	router.HandleFunc("/sensor/{id}/values/{limit}", api.GetSensorValues(wdb)).Methods("GET")

	// secured API
	router.HandleFunc("/sensor_type", basicAuth(api.AddSensorType(wdb))).Methods("POST")
	router.HandleFunc("/sensor_type/{id}", basicAuth(api.UpdateSensorType(wdb))).Methods("PUT")
	router.HandleFunc("/sensor_type/{id}", basicAuth(api.DeleteSensorType(wdb))).Methods("DELETE")

	router.HandleFunc("/sensor", basicAuth(api.AddSensor(wdb))).Methods("POST")
	router.HandleFunc("/sensor/{id}", basicAuth(api.UpdateSensor(wdb))).Methods("PUT")
	router.HandleFunc("/sensor/{id}", basicAuth(api.DeleteSensor(wdb))).Methods("DELETE")

	router.HandleFunc("/sensor/{id}/value", basicAuth(api.AddSensorValue(wdb))).Methods("POST")
	router.HandleFunc("/sensor/{id}/values", basicAuth(api.DeleteSensorValues(wdb))).Methods("DELETE")

	router.HandleFunc("/housekeeping", basicAuth(api.Housekeeping(wdb))).Methods("POST")

	return router
}

func basicAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		user, pass, _ := req.BasicAuth()
		username, password := util.GetUserAndPassword()
		if user != username && pass != password {
			http.Error(rw, "Unauthorized.", 401)
			return
		}
		fn(rw, req)
	}
}
