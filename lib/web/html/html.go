package html

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/config"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/forecasts"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func Index(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Info",
			Active: "home",
		}
		web.Render().HTML(rw, http.StatusOK, "index", page)
	}
}

func NotFound(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title: "Home Info - Not Found",
	}
	web.Render().HTML(rw, http.StatusNotFound, "not_found", page)
}

func ErrorHandler(rw http.ResponseWriter, req *http.Request) {
	Error(rw, fmt.Errorf("Internal Server Error"))
}
func Error(rw http.ResponseWriter, err error) {
	page := &Page{
		Title:   "Home Info - Error",
		Content: err,
	}
	web.Render().HTML(rw, http.StatusInternalServerError, "error", page)
}

func Graphs(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Info - Graphs",
			Active: "graphs",
		}

		sensors, err := hdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		var weeklyLabels, hourlyLabels []string
		weeklyTemperature := make(map[database.Sensor][]*database.SensorValue)
		hourlyTemperature := make(map[database.Sensor][]*database.SensorValue)
		weeklyHumidity := make(map[database.Sensor][]*database.SensorValue)
		hourlyHumidity := make(map[database.Sensor][]*database.SensorValue)
		for _, sensor := range sensors {
			switch sensor.Type {
			case "temperature":
				values, err := hdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(rw, err)
					return
				}
				hourlyTemperature[*sensor] = values

				if len(hourlyLabels) == 0 && sensor.Id == config.Get().Room.TemperatureSensorID {
					// collect labels
					for _, value := range values {
						hourlyLabels = append(hourlyLabels, value.Timestamp.Format("02.01. - 15:04"))
					}
				}

				values, err = hdb.GetDailyAverages(sensor.Id, 28)
				if err != nil {
					Error(rw, err)
					return
				}
				weeklyTemperature[*sensor] = values

				if len(weeklyLabels) == 0 && sensor.Id == config.Get().Room.TemperatureSensorID {
					// collect labels
					for _, value := range values {
						weeklyLabels = append(weeklyLabels, value.Timestamp.Format("02.01.2006"))
					}
				}
			case "humidity":
				values, err := hdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(rw, err)
					return
				}
				hourlyHumidity[*sensor] = values

				values, err = hdb.GetDailyAverages(sensor.Id, 48)
				if err != nil {
					Error(rw, err)
					return
				}
				weeklyHumidity[*sensor] = values
			}
		}

		page.Content = struct {
			HourlyTemperature map[database.Sensor][]*database.SensorValue
			HourlyHumidity    map[database.Sensor][]*database.SensorValue
			HourlyLabels      []string
			WeeklyTemperature map[database.Sensor][]*database.SensorValue
			WeeklyHumidity    map[database.Sensor][]*database.SensorValue
			WeeklyLabels      []string
		}{
			hourlyTemperature,
			hourlyHumidity,
			hourlyLabels,
			weeklyTemperature,
			weeklyHumidity,
			weeklyLabels,
		}

		web.Render().HTML(rw, http.StatusOK, "graphs", page)
	}
}

func Sensors(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Info - Sensors",
			Active: "sensors",
		}

		sensors, err := hdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		data := make(map[int][]*database.SensorData, 0)
		for _, sensor := range sensors {
			d, err := hdb.GetSensorData(sensor.Id, 10)
			if err != nil {
				Error(rw, err)
				return
			}
			data[sensor.Id] = d
		}

		page.Content = struct {
			Sensors    []*database.Sensor
			SensorData map[int][]*database.SensorData
		}{
			sensors,
			data,
		}
		web.Render().HTML(rw, http.StatusOK, "sensors", page)
	}
}

func Forecasts(rw http.ResponseWriter, req *http.Request) {
	canton, city := web.GetLocation(req)
	page := &Page{
		Title:  fmt.Sprintf("Home Info - Forecasts - %s", city),
		Active: "forecasts",
	}

	forecast, err := forecasts.Get(canton, city)
	if err != nil {
		page.Content = struct {
			Canton string
			City   string
			Error  error
		}{
			canton,
			city,
			err,
		}
		web.Render().HTML(rw, http.StatusNotFound, "forecast_error", page)
		return
	}

	page.Content = struct {
		Canton           string
		City             string
		Forecast         forecasts.WeatherForecast
		Today            time.Time
		Tomorrow         time.Time
		DayAfterTomorrow time.Time
	}{
		canton,
		city,
		forecast,
		time.Now(),
		time.Now().AddDate(0, 0, 1),
		time.Now().AddDate(0, 0, 2),
	}
	web.Render().HTML(rw, http.StatusOK, "forecasts", page)
}
