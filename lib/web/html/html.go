package html

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/forecasts"
	"github.com/anyandrea/weather_app/lib/web"
)

func Index(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Weather App",
			Active: "home",
		}

		windows, err := wdb.GetWindowStates()
		if err != nil {
			Error(rw, err)
			return
		}

		page.Content = struct {
			Windows []*weatherdb.Window
		}{
			windows,
		}
		web.Render().HTML(rw, http.StatusOK, "index", page)
	}
}

func NotFound(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title: "Weather App - Not Found",
	}
	web.Render().HTML(rw, http.StatusNotFound, "not_found", page)
}

func ErrorHandler(rw http.ResponseWriter, req *http.Request) {
	Error(rw, fmt.Errorf("Internal Server Error"))
}
func Error(rw http.ResponseWriter, err error) {
	page := &Page{
		Title:   "Weather App - Error",
		Content: err,
	}
	web.Render().HTML(rw, http.StatusInternalServerError, "error", page)
}

func Sensors(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Weather App - Sensors",
			Active: "sensors",
		}

		sensors, err := wdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}

		data := make(map[int][]*weatherdb.SensorData, 0)
		for _, sensor := range sensors {
			d, err := wdb.GetSensorData(sensor.Id, 10)
			if err != nil {
				Error(rw, err)
				return
			}
			data[sensor.Id] = d
		}

		page.Content = struct {
			Sensors    []*weatherdb.Sensor
			SensorData map[int][]*weatherdb.SensorData
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
		Title:  fmt.Sprintf("Weather App - Forecasts - %s", city),
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
