package api

import (
	"net/http"

	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/web"
)

func GetSensorTypes(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		sensorTypes, err := wdb.GetSensorTypes()
		if err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusOK, sensorTypes)
	}
}

func AddSensorType(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sensorType := &weatherdb.SensorType{
			Type:        req.Form.Get("type"),
			Unit:        req.Form.Get("unit"),
			Description: req.Form.Get("description"),
		}

		if err := wdb.InsertSensorType(sensorType); err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusCreated, *sensorType)
	}
}
