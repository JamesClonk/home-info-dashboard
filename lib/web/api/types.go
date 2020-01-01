package api

import (
	"net/http"
	"strconv"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
	"github.com/gorilla/mux"
)

func GetSensorTypes(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		sensorTypes, err := hdb.GetSensorTypes()
		if err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusOK, sensorTypes)
	}
}

func GetSensorType(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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

			sensorType, err := hdb.GetSensorTypeById(int(sensorTypeId))
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

func AddSensorType(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sensorType := &database.SensorType{
			Type:        req.Form.Get("type"),
			Unit:        req.Form.Get("unit"),
			Description: req.Form.Get("description"),
		}

		if err := hdb.InsertSensorType(sensorType); err != nil {
			Error(rw, err)
			return
		}

		web.Render().JSON(rw, http.StatusCreated, *sensorType)
	}
}

func UpdateSensorType(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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

			req.ParseForm()
			sensorType := &database.SensorType{
				Id:          int(sensorTypeId),
				Type:        req.Form.Get("type"),
				Unit:        req.Form.Get("unit"),
				Description: req.Form.Get("description"),
			}

			if err := hdb.UpdateSensorType(sensorType); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, *sensorType)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteSensorType(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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

			if err := hdb.DeleteSensorType(int(sensorTypeId)); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
