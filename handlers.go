package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var r *render.Render

type Page struct {
	Title   string
	Active  string
	Content interface{}
}

func init() {
	// setup template rendering
	r = render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		IndentJSON: true,
	})
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title:  "Weather App",
		Active: "home",
	}
	r.HTML(rw, http.StatusOK, "index", page)
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title: "Weather App - Not Found",
	}
	r.HTML(rw, http.StatusNotFound, "not_found", page)
}

func ErrorHandler(rw http.ResponseWriter, req *http.Request) {
	Error(rw, fmt.Errorf("Internal Server Error"))
}
func Error(rw http.ResponseWriter, err error) {
	page := &Page{
		Title:   "Weather App - Error",
		Content: err,
	}
	r.HTML(rw, http.StatusInternalServerError, "error", page)
}

func ForecastsHandler(rw http.ResponseWriter, req *http.Request) {
	canton, city := getLocation(req)
	page := &Page{
		Title:  fmt.Sprintf("Weather App - Forecasts - %s", city),
		Active: "forecasts",
	}

	forecast, err := GetWeatherForecast(canton, city)
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
		r.HTML(rw, http.StatusNotFound, "forecast_error", page)
		return
	}

	page.Content = struct {
		Canton           string
		City             string
		Forecast         *WeatherForecast
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
	r.HTML(rw, http.StatusOK, "forecasts", page)
}

func getLocation(req *http.Request) (string, string) {
	// first, try to read values from gorilla mux
	vars := mux.Vars(req)
	canton := vars["canton"]
	city := vars["city"]

	// then, parse the form and try to read the values from POST data
	req.ParseForm()
	if len(canton) == 0 {
		canton = req.Form.Get("canton")
	}
	if len(city) == 0 {
		city = req.Form.Get("city")
	}

	// now, try to read defaults from ENV
	if len(canton) == 0 {
		canton = os.Getenv("DEFAULT_CANTON")
	}
	if len(city) == 0 {
		city = os.Getenv("DEFAULT_CITY")
	}

	// if we still have no values, set reasonable defaults
	if len(canton) == 0 {
		canton = "Bern"
	}
	if len(city) == 0 {
		city = "Bern"
	}

	return canton, city
}
