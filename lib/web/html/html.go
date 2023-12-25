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
			Title:  "Home Automation",
			Active: "home",
		}
		_ = web.Render().HTML(rw, http.StatusOK, "index", page)
	}
}

func NotFound(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title: "Home Automation - Not Found",
	}
	_ = web.Render().HTML(rw, http.StatusNotFound, "not_found", page)
}

func ErrorHandler(rw http.ResponseWriter, req *http.Request) {
	Error(rw, fmt.Errorf("Internal Server Error"))
}
func Error(rw http.ResponseWriter, err error) {
	page := &Page{
		Title:   "Home Automation - Error",
		Content: err,
	}
	_ = web.Render().HTML(rw, http.StatusInternalServerError, "error", page)
}

func Graphs(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Automation - Graphs",
			Active: "graphs",
		}

		sensors, err := hdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		var weeklyLabels, hourlyLabels []string
		weeklyTemperature := make(map[string][]*database.SensorValue)
		hourlyTemperature := make(map[string][]*database.SensorValue)
		weeklyHumidity := make(map[string][]*database.SensorValue)
		hourlyHumidity := make(map[string][]*database.SensorValue)
		weeklyMoisture := make(map[string][]*database.SensorValue)
		hourlyMoisture := make(map[string][]*database.SensorValue)
		weeklyCo2 := make(map[string][]*database.SensorValue)
		hourlyCo2 := make(map[string][]*database.SensorValue)
		for _, sensor := range sensors {
			switch sensor.SensorType.Type {
			case "temperature":
				values, err := hdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(rw, err)
					return
				}
				if len(values) == 0 { // dont show empty sensors
					continue
				}
				hourlyTemperature[sensor.Name] = values

				if len(hourlyLabels) == 0 && sensor.Id == config.Get().LivingRoom.TemperatureSensorID {
					// collect labels
					for _, value := range values {
						hourlyLabels = append(hourlyLabels, value.Timestamp.Format("02.01.2006 - 15:04"))
					}
				}

				values, err = hdb.GetDailyAverages(sensor.Id, 28)
				if err != nil {
					Error(rw, err)
					return
				}
				weeklyTemperature[sensor.Name] = values

				if len(weeklyLabels) == 0 && sensor.Id == config.Get().LivingRoom.TemperatureSensorID {
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
				if len(values) == 0 { // dont show empty sensors
					continue
				}
				hourlyHumidity[sensor.Name] = values

				values, err = hdb.GetDailyAverages(sensor.Id, 28)
				if err != nil {
					Error(rw, err)
					return
				}
				weeklyHumidity[sensor.Name] = values
			case "soil":
				values, err := hdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(rw, err)
					return
				}
				if len(values) == 0 { // dont show empty sensors
					continue
				}
				hourlyMoisture[sensor.Name] = values

				values, err = hdb.GetDailyAverages(sensor.Id, 28)
				if err != nil {
					Error(rw, err)
					return
				}
				weeklyMoisture[sensor.Name] = values
			case "co2":
				values, err := hdb.GetHourlyAverages(sensor.Id, 48)
				if err != nil {
					Error(rw, err)
					return
				}
				if len(values) == 0 { // dont show empty sensors
					continue
				}
				hourlyCo2[sensor.Name] = values

				values, err = hdb.GetDailyAverages(sensor.Id, 28)
				if err != nil {
					Error(rw, err)
					return
				}
				weeklyCo2[sensor.Name] = values
			}
		}

		page.Content = struct {
			HourlyTemperature map[string][]*database.SensorValue
			HourlyHumidity    map[string][]*database.SensorValue
			HourlyMoisture    map[string][]*database.SensorValue
			HourlyCo2         map[string][]*database.SensorValue
			HourlyLabels      []string
			WeeklyTemperature map[string][]*database.SensorValue
			WeeklyHumidity    map[string][]*database.SensorValue
			WeeklyMoisture    map[string][]*database.SensorValue
			WeeklyCo2         map[string][]*database.SensorValue
			WeeklyLabels      []string
		}{
			hourlyTemperature,
			hourlyHumidity,
			hourlyMoisture,
			hourlyCo2,
			hourlyLabels,
			weeklyTemperature,
			weeklyHumidity,
			weeklyMoisture,
			weeklyCo2,
			weeklyLabels,
		}

		_ = web.Render().HTML(rw, http.StatusOK, "graphs", page)
	}
}

func Sensors(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Automation - Sensors",
			Active: "sensors",
		}

		sensorTypes, err := hdb.GetSensorTypes()
		if err != nil {
			Error(rw, err)
			return
		}

		sensors, err := hdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		data := make(map[int][]*database.SensorData)
		for _, sensor := range sensors {
			d, err := hdb.GetSensorData(sensor.Id, 15)
			if err != nil {
				Error(rw, err)
				return
			}
			data[sensor.Id] = d
		}

		page.Content = struct {
			SensorTypes []*database.SensorType
			Sensors     []*database.Sensor
			SensorData  map[int][]*database.SensorData
		}{
			sensorTypes,
			sensors,
			data,
		}
		_ = web.Render().HTML(rw, http.StatusOK, "sensors", page)
	}
}

func Alerts(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Automation - Alerts",
			Active: "alerts",
		}

		alerts, err := hdb.GetAllAlerts()
		if err != nil {
			Error(rw, err)
			return
		}

		page.Content = struct {
			Alerts []*database.Alert
		}{
			alerts,
		}
		_ = web.Render().HTML(rw, http.StatusOK, "alerts", page)
	}
}

func Forecasts(rw http.ResponseWriter, req *http.Request) {
	lat, lon, alt := web.GetLocation(req)
	page := &Page{
		Title:  fmt.Sprintf("Home Automation - Forecasts - %s/%s", lat, lon),
		Active: "forecasts",
	}

	forecast, err := forecasts.Get(lat, lon, alt)
	if err != nil {
		page.Content = struct {
			Latitude  string
			Longitude string
			Altitude  string
			Error     error
		}{
			lat,
			lon,
			alt,
			err,
		}
		_ = web.Render().HTML(rw, http.StatusNotFound, "forecast_error", page)
		return
	}

	page.Content = struct {
		Latitude         string
		Longitude        string
		Altitude         string
		Forecast         forecasts.WeatherForecast
		Today            time.Time
		Tomorrow         time.Time
		DayAfterTomorrow time.Time
	}{
		lat,
		lon,
		alt,
		forecast,
		time.Now().Local(),
		time.Now().Local().AddDate(0, 0, 1),
		time.Now().Local().AddDate(0, 0, 2),
	}
	_ = web.Render().HTML(rw, http.StatusOK, "forecasts", page)
}
