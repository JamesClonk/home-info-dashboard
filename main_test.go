package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

var (
	n            *negroni.Negroni
	testUser     = "test"
	testPassword = "test12345"
)

func init() {
	os.Setenv("PORT", "8080")

	os.Setenv("DEFAULT_CANTON", "Bern")
	os.Setenv("DEFAULT_CITY", "Bern")

	os.Setenv("WEATHERAPI_USERNAME", testUser)
	os.Setenv("WEATHERAPI_PASSWORD", testPassword)

	os.Setenv("WEATHERDB_TYPE", "sqlite")
	os.Setenv("WEATHERDB_URI", "sqlite3://_fixtures/temp.db")
	os.Setenv("WEATHERDB_TYPE", "sqlite")

	if err := exec.Command("cp", "-f", "_fixtures/test.db", "_fixtures/temp.db").Run(); err != nil {
		log.Fatal(err)
	}

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
	// ----------------------- Unauthorized -----------------------
	response := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/sensor_type", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	body := response.Body.String()
	assert.Contains(t, response.Body.String(), `Unauthorized`)

	// ----------------------- CREATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/sensor_type", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form := url.Values{}
	form.Add("type", "pressure")
	form.Add("unit", "psi")
	form.Add("description", "atmospheric pressure")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusCreated, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "type": "pressure",
  "unit": "psi",
  "description": "atmospheric pressure"
}`)

	// ----------------------- READ -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor_type", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
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

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor_type/5", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "type": "pressure",
  "unit": "psi",
  "description": "atmospheric pressure"
}`)

	// ----------------------- UPDATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("PUT", "/sensor_type/5", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form = url.Values{}
	form.Add("type", "druck")
	form.Add("unit", "bar")
	form.Add("description", "atmosphärischer druck")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "type": "druck",
  "unit": "bar",
  "description": "atmosphärischer druck"
}`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor_type/5", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "type": "druck",
  "unit": "bar",
  "description": "atmosphärischer druck"
}`)

	// ----------------------- DELETE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/sensor_type/5", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNoContent, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `null`)

	// is sensor_type gone?
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor_type/5", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `"sql: no rows in result set"`)
}

func Test_Main_Sensors(t *testing.T) {
	// ----------------------- Unauthorized -----------------------
	response := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/sensor", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	body := response.Body.String()
	assert.Contains(t, response.Body.String(), `Unauthorized`)

	// ----------------------- CREATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/sensor", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form := url.Values{}
	form.Add("name", "Badezimmer")
	form.Add("type", "temperature")
	form.Add("description", "Badezimmer Temperatur")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusCreated, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "name": "Badezimmer",
  "type": "temperature",
  "type_id": "3",
  "unit": "celsius",
  "description": "Badezimmer Temperatur"
}`)

	// ----------------------- READ -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `
    "name": "temperature #1",
    "type": "temperature",
    "type_id": "3",
    "unit": "celsius",
    "description": "Shows temperature"`)
	assert.Contains(t, body, `
    "id": 5,
    "name": "Badezimmer",
    "type": "temperature",
    "type_id": "3",
    "unit": "celsius",
    "description": "Badezimmer Temperatur"`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor/1", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 1,
  "name": "roof window #1",
  "type": "window_state",
  "type_id": "1",
  "unit": "closed",
  "description": "Shows open/closed state of roof window"
}`)

	// ----------------------- UPDATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("PUT", "/sensor/5", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form = url.Values{}
	form.Add("name", "Wohnzimmer")
	form.Add("type", "humidity")
	form.Add("description", "Wohnzimmer Luftfeuchtigkeit")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "name": "Wohnzimmer",
  "type": "humidity",
  "type_id": "4",
  "unit": "percentage",
  "description": "Wohnzimmer Luftfeuchtigkeit"
}`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor/5", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "name": "Wohnzimmer",
  "type": "humidity",
  "type_id": "4",
  "unit": "percentage",
  "description": "Wohnzimmer Luftfeuchtigkeit"
}`)

	// ----------------------- DELETE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/sensor/5", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNoContent, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `null`)

	// is sensor gone?
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor/5", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `"sql: no rows in result set"`)
}

func Test_Main_SensorValues(t *testing.T) {
	// ----------------------- Unauthorized -----------------------
	response := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/sensor/3/value", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	body := response.Body.String()
	assert.Contains(t, response.Body.String(), `Unauthorized`)

	// ----------------------- CREATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/sensor/3/value", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form := url.Values{}
	form.Add("value", "123456789")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusCreated, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `123456789`)

	// ----------------------- READ -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor/3/values", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `"value": 123456789`)
	assert.Contains(t, body, `
  {
    "timestamp": "1974-10-01T07:02:15+01:00",
    "value": 15
  },
  {
    "timestamp": "1973-07-30T15:47:18+01:00",
    "value": 41
  }`)

	// ----------------------- DELETE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/sensor/3/values", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNoContent, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `null`)

	// is everything gone?
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/sensor/3/values", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `[]`)
}
