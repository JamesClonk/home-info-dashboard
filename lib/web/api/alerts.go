package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
	"github.com/gorilla/mux"
)

func GetAlerts(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		alerts, err := hdb.GetAllAlerts()
		if err != nil {
			Error(rw, err)
			return
		}
		_ = web.Render().JSON(rw, http.StatusOK, alerts)
	}
}

func GetAlert(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var alertId int64
			if len(id) > 0 {
				alertId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			alert, err := hdb.GetAlertById(int(alertId))
			if err != nil {
				Error(rw, err)
				return
			}

			_ = web.Render().JSON(rw, http.StatusOK, alert)
			return
		}
		_ = web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func AddAlert(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			Error(rw, err)
			return
		}
		sensorId, err := strconv.Atoi(req.Form.Get("sensor_id"))
		if err != nil {
			Error(rw, err)
			return
		}
		if sensorId == 0 {
			Error(rw, fmt.Errorf("sensor_id missing!"))
			return
		}
		silenceDuration, err := strconv.ParseInt(req.Form.Get("silence_duration"), 10, 64)
		if err != nil {
			Error(rw, err)
			return
		}
		alert := &database.Alert{
			Name:            req.Form.Get("name"),
			Description:     req.Form.Get("description"),
			Condition:       req.Form.Get("condition"),
			Execution:       req.Form.Get("execution"),
			SilenceDuration: silenceDuration,
			Active:          sql.NullBool{true, true},
			Sensor: database.Sensor{
				Id: sensorId,
			},
		}

		if err := hdb.InsertAlert(alert); err != nil {
			Error(rw, err)
			return
		}
		_ = web.Render().JSON(rw, http.StatusCreated, *alert)
	}
}

func UpdateAlert(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var alertId int64
			if len(id) > 0 {
				alertId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			_ = req.ParseForm()
			sensorId, err := strconv.Atoi(req.Form.Get("sensor_id"))
			if err != nil {
				Error(rw, err)
				return
			}
			if sensorId == 0 {
				Error(rw, fmt.Errorf("sensor_id missing!"))
				return
			}
			silenceDuration, err := strconv.ParseInt(req.Form.Get("silence_duration"), 10, 64)
			if err != nil {
				Error(rw, err)
				return
			}
			activeStr := req.Form.Get("active")
			if len(activeStr) == 0 {
				activeStr = "true"
			}
			active, err := strconv.ParseBool(activeStr)
			if err != nil {
				Error(rw, err)
				return
			}
			lastAlert, err := time.Parse(time.RFC3339, req.Form.Get("last_alert"))
			if err != nil {
				Error(rw, err)
				return
			}
			alert := &database.Alert{
				Id:              int(alertId),
				Name:            req.Form.Get("name"),
				Description:     req.Form.Get("description"),
				Condition:       req.Form.Get("condition"),
				Execution:       req.Form.Get("execution"),
				LastAlert:       &lastAlert,
				SilenceDuration: silenceDuration,
				Active:          sql.NullBool{active, true},
				Sensor: database.Sensor{
					Id: sensorId,
				},
			}

			if err := hdb.UpdateAlert(alert); err != nil {
				Error(rw, err)
				return
			}
			_ = web.Render().JSON(rw, http.StatusOK, *alert)
			return
		}
		_ = web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteAlert(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var alertId int64
			if len(id) > 0 {
				alertId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			if err := hdb.DeleteAlert(int(alertId)); err != nil {
				Error(rw, err)
				return
			}
			_ = web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}
		_ = web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}
