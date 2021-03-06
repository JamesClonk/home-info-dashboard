package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	sensorError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "homeinfo_dashboard_sensor_errors_total",
		Help: "Total number of Home-Info Dashboard sensor errors.",
	})
)

func GetSensorValues(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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
					sensorError.Inc()
					return
				}
			}
			if len(limit) > 0 {
				valueLimit, err = strconv.ParseInt(limit, 10, 64)
				if err != nil {
					Error(rw, err)
					sensorError.Inc()
					return
				}
			} else {
				valueLimit = 100
			}

			values, err := hdb.GetSensorValues(int(sensorId), int(valueLimit))
			if err != nil {
				Error(rw, err)
				sensorError.Inc()
				return
			}

			web.Render().JSON(rw, http.StatusOK, values)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func AddSensorValue(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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
					sensorError.Inc()
					return
				}
			}

			req.ParseForm()
			value, err = strconv.ParseInt(req.Form.Get("value"), 10, 64)
			if err != nil {
				Error(rw, err)
				sensorError.Inc()
				return
			}

			if err := hdb.InsertSensorValue(int(sensorId), int(value), time.Now()); err != nil {
				Error(rw, err)
				sensorError.Inc()
				return
			}

			web.Render().JSON(rw, http.StatusCreated, value)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteSensorValues(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
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
					sensorError.Inc()
					return
				}
			}

			if err := hdb.DeleteSensorValues(int(sensorId)); err != nil {
				Error(rw, err)
				sensorError.Inc()
				return
			}

			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func Housekeeping(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if err := hdb.Housekeeping(); err != nil {
			Error(rw, err)
			sensorError.Inc()
			return
		}
		web.Render().JSON(rw, http.StatusNoContent, nil)
	}
}
