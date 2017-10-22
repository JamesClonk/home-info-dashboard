package main

import (
	"fmt"
	"net/http"
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
		Title: "Weather App - 404",
	}
	r.HTML(rw, http.StatusOK, "not_found", page)
}

func ForecastsHandler(rw http.ResponseWriter, req *http.Request) {
	page := &Page{
		Title:  "Weather App - Forecasts",
		Active: "forecasts",
	}

	vars := mux.Vars(req)
	canton := vars["canton"]
	city := vars["city"]

	req.ParseForm()
	if len(canton) == 0 {
		canton = req.Form.Get("canton")
	}
	if len(city) == 0 {
		city = req.Form.Get("city")
	}
	page.Title += fmt.Sprintf(" - %s", city)

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
