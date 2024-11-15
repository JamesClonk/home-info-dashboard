package router

import (
	"crypto/subtle"
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/util"
	"github.com/JamesClonk/home-info-dashboard/lib/web/api"
	"github.com/JamesClonk/home-info-dashboard/lib/web/html"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(hdb database.HomeInfoDB) *mux.Router {
	router := mux.NewRouter()
	setupRoutes(hdb, router)
	return router
}

func setupRoutes(hdb database.HomeInfoDB, router *mux.Router) *mux.Router {
	// HTML
	router.NotFoundHandler = http.HandlerFunc(html.NotFound)

	router.HandleFunc("/", html.Index(hdb))
	router.HandleFunc("/error", html.ErrorHandler)

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/dashboard", html.Dashboard(hdb))
	router.HandleFunc("/graphs", html.Graphs(hdb))
	router.HandleFunc("/sensor_data", html.Sensors(hdb))
	router.HandleFunc("/alert_data", html.Alerts(hdb))

	router.HandleFunc("/forecasts", html.Forecasts)
	router.HandleFunc("/forecasts/{canton}", html.Forecasts)
	router.HandleFunc("/forecasts/{canton}/{city}", html.Forecasts)

	router.HandleFunc("/fitness", html.Fitness(hdb)).Methods("GET")

	// secure HTML
	router.HandleFunc("/fitness", basicAuth(html.Fitness(hdb))).Methods("POST")
	router.HandleFunc("/message_queue", basicAuth(html.Messages(hdb)))

	// API
	router.HandleFunc("/sensor_type", api.GetSensorTypes(hdb)).Methods("GET")
	router.HandleFunc("/sensor_types", api.GetSensorTypes(hdb)).Methods("GET")
	router.HandleFunc("/sensor_type/{id}", api.GetSensorType(hdb)).Methods("GET")

	router.HandleFunc("/sensor", api.GetSensors(hdb)).Methods("GET")
	router.HandleFunc("/sensors", api.GetSensors(hdb)).Methods("GET")
	router.HandleFunc("/sensor/{id}", api.GetSensor(hdb)).Methods("GET")

	router.HandleFunc("/sensor/{id}/values", api.GetSensorValues(hdb)).Methods("GET")
	router.HandleFunc("/sensor/{id}/values/{limit}", api.GetSensorValues(hdb)).Methods("GET")

	router.HandleFunc("/alert", api.GetAlerts(hdb)).Methods("GET")
	router.HandleFunc("/alerts", api.GetAlerts(hdb)).Methods("GET")
	router.HandleFunc("/alert/{id}", api.GetAlert(hdb)).Methods("GET")

	router.HandleFunc("/queues", api.GetQueues(hdb)).Methods("GET")
	router.HandleFunc("/queue/{queue}", api.GetMessagesByQueue(hdb)).Methods("GET")
	router.HandleFunc("/messages", api.GetAllMessages(hdb)).Methods("GET")
	router.HandleFunc("/messages/{queue}", api.GetMessagesByQueue(hdb)).Methods("GET")
	router.HandleFunc("/message/{id}", api.GetMessage(hdb)).Methods("GET")

	// secured API
	router.HandleFunc("/sensor_type", basicAuth(api.AddSensorType(hdb))).Methods("POST")
	router.HandleFunc("/sensor_type/{id}", basicAuth(api.UpdateSensorType(hdb))).Methods("PUT")
	router.HandleFunc("/sensor_type/{id}", basicAuth(api.DeleteSensorType(hdb))).Methods("DELETE")

	router.HandleFunc("/sensor", basicAuth(api.AddSensor(hdb))).Methods("POST")
	router.HandleFunc("/sensor/{id}", basicAuth(api.UpdateSensor(hdb))).Methods("PUT")
	router.HandleFunc("/sensor/{id}", basicAuth(api.DeleteSensor(hdb))).Methods("DELETE")

	router.HandleFunc("/sensor/{id}/value", basicAuth(api.AddSensorValue(hdb))).Methods("POST")
	router.HandleFunc("/sensor/{id}/values", basicAuth(api.DeleteSensorValues(hdb))).Methods("DELETE")

	router.HandleFunc("/queue/{queue}", basicAuth(api.DeleteQueue(hdb))).Methods("DELETE")
	router.HandleFunc("/message", basicAuth(api.AddMessage(hdb))).Methods("POST")
	router.HandleFunc("/message/{id}", basicAuth(api.DeleteMessage(hdb))).Methods("DELETE")

	router.HandleFunc("/alert", basicAuth(api.AddAlert(hdb))).Methods("POST")
	router.HandleFunc("/alert/{id}", basicAuth(api.UpdateAlert(hdb))).Methods("PUT")
	router.HandleFunc("/alert/{id}", basicAuth(api.DeleteAlert(hdb))).Methods("DELETE")

	router.HandleFunc("/housekeeping", basicAuth(api.Housekeeping(hdb))).Methods("POST")

	router.HandleFunc("/telebot", basicAuth(api.TelebotStatus())).Methods("GET") // needs to be secured too, exposes Telebot config information
	router.HandleFunc("/telebot", basicAuth(api.TelebotInit())).Methods("PUT")
	router.HandleFunc("/telebot/message", basicAuth(api.TelebotMessage())).Methods("POST")

	router.HandleFunc("/slack", basicAuth(api.SlackStatus())).Methods("GET") // needs to be secured too, exposes Slack config information
	router.HandleFunc("/slack/message", basicAuth(api.SlackMessage())).Methods("POST")

	router.HandleFunc("/teams", basicAuth(api.TeamsStatus())).Methods("GET") // needs to be secured too, exposes Teams config information
	router.HandleFunc("/teams/message", basicAuth(api.TeamsMessage())).Methods("POST")

	return router
}

func basicAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("WWW-Authenticate", `Basic realm="Home-Info"`)

		user, pass, _ := req.BasicAuth()
		username, password := util.GetUserAndPassword()
		if len(username) < 8 && len(password) < 8 {
			http.Error(rw, "Unauthorized", 401)
			return
		}

		if subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			http.Error(rw, "Unauthorized.", 401)
			return
		}
		fn(rw, req)
	}
}
