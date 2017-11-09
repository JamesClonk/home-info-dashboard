package api

import (
	"net/http"
	"strconv"

	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/web"
	"github.com/gorilla/mux"
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

func GetSensorType(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var sensorTypeId int64
			if len(id) > 0 {
				sensorTypeId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			sensorType, err := wdb.GetSensorTypeById(int(sensorTypeId))
			if err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, sensorType)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
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

func DeleteSensorType(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var sensorTypeId int64
			if len(id) > 0 {
				sensorTypeId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			if err := wdb.DeleteSensorType(int(sensorTypeId)); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
