package api

import (
	"net/http"
	"strconv"

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

			sensor, err := wdb.GetSensorById(int(sensorId))
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

func UpdateSensor(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
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

			req.ParseForm()
			sensor := &weatherdb.Sensor{
				Id:          int(sensorId),
				Name:        req.Form.Get("name"),
				Type:        req.Form.Get("type"),
				TypeId:      req.Form.Get("type_id"),
				Description: req.Form.Get("description"),
			}

			if err := wdb.UpdateSensor(sensor); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, *sensor)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteSensor(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
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

			if err := wdb.DeleteSensor(int(sensorId)); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
