package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/web"
	"github.com/gorilla/mux"
)

func GetSensors(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		sensors, err := wdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusOK, sensors)
	}
}

func GetSensor(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var sensorId int64
			if len(id) > 0 {
				sensorId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			sensor, err := wdb.GetSensor(int(sensorId))
			if err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, sensor)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func GetSensorValues(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]
		limit := vars["limit"]

		if len(id) > 0 {
			var sensorId, valueLimit int64

			if len(id) > 0 {
				sensorId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}
			if len(limit) > 0 {
				valueLimit, err = strconv.ParseInt(limit, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			} else {
				valueLimit = 100
			}

			values, err := wdb.GetSensorValues(int(sensorId), int(valueLimit))
			if err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, values)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func AddSensor(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sensor := &weatherdb.Sensor{
			Name:        req.Form.Get("name"),
			Type:        req.Form.Get("type"),
			TypeId:      req.Form.Get("type_id"),
			Description: req.Form.Get("description"),
		}

		if err := wdb.InsertSensor(sensor); err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusCreated, *sensor)
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

func AddSensorValue(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var sensorId, value int64
			if len(id) > 0 {
				sensorId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			req.ParseForm()
			value, err = strconv.ParseInt(req.Form.Get("value"), 10, 64)
			if err != nil {
				Error(rw, err)
				return
			}

			if err := wdb.InsertSensorValue(int(sensorId), int(value), time.Now()); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusCreated, value)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
