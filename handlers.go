package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var r *render.Render

type Page struct {
	Title   string
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
		Title: "Weather App",
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
		Title: "Weather App - Forecasts",
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

	forecast, err := GetWeatherForecast(canton, city)
	if err != nil {
		page.Content = err
		r.HTML(rw, http.StatusInternalServerError, "error", page)
		return
	}

	page.Title += fmt.Sprintf(" - %s", city)
	page.Content = forecast
	r.HTML(rw, http.StatusOK, "forecasts", page)
}
