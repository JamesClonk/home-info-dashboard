package main

import (
	"net/http"

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

	data, err := getData("http://www.yr.no/place/Switzerland/Bern/Bützberg/forecast.xml")
	if err != nil {
		page.Content = err
		r.HTML(rw, http.StatusInternalServerError, "error", page)
		return
	}

	page.Content = string(data)
	r.HTML(rw, http.StatusOK, "forecasts", page)
}
