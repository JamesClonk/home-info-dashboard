package api

import (
	"net/http"
	"strconv"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
	"github.com/gorilla/mux"
)

func GetSensors(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		sensors, err := hdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusOK, sensors)
	}
}

func GetSensor(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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

			sensor, err := hdb.GetSensorById(int(sensorId))
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

func AddSensor(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sensorTypeId, err := strconv.Atoi(req.Form.Get("type_id"))
		if err != nil {
			sensorTypeId = 0
		}
		sensor := &database.Sensor{
			Name:        req.Form.Get("name"),
			Description: req.Form.Get("description"),
			SensorType: database.SensorType{
				Id:   sensorTypeId,
				Type: req.Form.Get("type"),
			},
		}

		if err := hdb.InsertSensor(sensor); err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusCreated, *sensor)
	}
}

func UpdateSensor(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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
			sensorTypeId, err := strconv.Atoi(req.Form.Get("type_id"))
			if err != nil {
				sensorTypeId = 0
			}
			sensor := &database.Sensor{
				Id:          int(sensorId),
				Name:        req.Form.Get("name"),
				Description: req.Form.Get("description"),
				SensorType: database.SensorType{
					Id:   sensorTypeId,
					Type: req.Form.Get("type"),
				},
			}

			if err := hdb.UpdateSensor(sensor); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, *sensor)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteSensor(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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

			if err := hdb.DeleteSensor(int(sensorId)); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
