package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

var (
	n *negroni.Negroni
)

func init() {
	os.Setenv("PORT", "8080")
	os.Setenv("DEFAULT_CANTON", "Bern")
	os.Setenv("DEFAULT_CITY", "Bern")
	os.Setenv("WEATHERDB_TYPE", "sqlite")
	os.Setenv("WEATHERDB_URI", "sqlite3://_fixtures/test.db")
	os.Setenv("WEATHERDB_TYPE", "sqlite")
	// TODO: copy & overwrite each time: _fixtures/fixtures.db to _fixtures/test.db
	db := setupDatabase()
	n = setupNegroni(db)

}

func Test_Main_404(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:8080/something", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>Weather App - Not Found</title>`)
	assert.Contains(t, body, `<section class="hero is-medium is-warning">`)
	assert.Contains(t, body, `<h1 class="title">404 - Not Found</h1>`)
	assert.Contains(t, body, `<h2 class="subtitle">This is not the page you are looking for ...</h2>`)
}

func Test_Main_500(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:8080/error", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>Weather App - Error</title>`)
	assert.Contains(t, body, `<section class="hero is-medium is-danger">`)
	assert.Contains(t, body, `<h1 class="title">Error</h1>`)
	assert.Contains(t, body, `<h2 class="subtitle">Internal Server Error</h2>`)
}

func Test_Main_Index(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>Weather App</title>`)
	assert.Contains(t, body, `<h1 class="title">Weather App</h1>`)
	assert.Contains(t, body, `<img src="/images/smart_temperature.png">`)
}

func Test_Main_Forecasts(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/forecasts", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>Weather App - Forecasts - Bern</title>`)
	assert.Contains(t, body, `<h1 class="title">Berne</h1>`)
	assert.Contains(t, body, `<h2 class="subtitle">Switzerland</h2>`)
	assert.Contains(t, body, `<p>549m</p>`)
	assert.Contains(t, body, `<a href="https://www.google.ch/maps/place/46.94809%C2%B0+7.44744%C2%B0" target="_blank" rel="noopener noreferrer">46.94809°/7.44744°</a>`)
	assert.Contains(t, body, `<p>Weather forecast from Yr, delivered by the Norwegian Meteorological Institute and the NRK<br/><a href="http://www.yr.no/place/Switzerland/Bern/Berne/">http://www.yr.no/place/Switzerland/Bern/Berne/</a></p>`)
}

func Test_Main_SensorTypes(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sensor_type", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `
    "type": "temperature",
    "unit": "celsius",
    "description": "Shows temperature"`)
	assert.Contains(t, body, `
  {
    "id": 4,
    "type": "humidity",
    "unit": "percentage",
    "description": "Shows air humidity"
  }`)

	// TODO: test create/update/delete
}

func Test_Main_Sensors(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sensor", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `
    "name": "temperature #1",
    "type": "temperature",
    "type_id": "3",
    "unit": "celsius",
    "description": "Shows temperature"`)
	assert.Contains(t, body, `
  {
    "id": 1,
    "name": "roof window #1",
    "type": "window_state",
    "type_id": "1",
    "unit": "closed",
    "description": "Shows open/closed state of roof window"
  }`)

	// TODO: test create/update/delete
}

func Test_Main_SensorValues(t *testing.T) {
	// response := httptest.NewRecorder()
	// req, err := http.NewRequest("GET", "/sensor/??/value", nil)
	// if err != nil {
	// 	t.Error(err)
	// }

	// TODO: test create/read/delete
}
