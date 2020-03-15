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
	os.Setenv("PORT", "9090")

	os.Setenv("DEFAULT_CANTON", "Bern")
	os.Setenv("DEFAULT_CITY", "Bern")

	os.Setenv("AUTH_USERNAME", testUser)
	os.Setenv("AUTH_PASSWORD", testPassword)

	os.Setenv("HOME_INFO_DB_TYPE", "sqlite")
	os.Setenv("HOME_INFO_DB_URI", "sqlite3://_fixtures/temp.db")

	os.Setenv("CONFIG_LIVING_ROOM_TEMPERATURE_SENSOR_ID", "3")
	os.Setenv("CONFIG_LIVING_ROOM_HUMIDITY_SENSOR_ID", "4")
	os.Setenv("CONFIG_BEDROOM_TEMPERATURE_SENSOR_ID", "6")
	os.Setenv("CONFIG_BEDROOM_HUMIDITY_SENSOR_ID", "7")
	os.Setenv("CONFIG_HOME_OFFICE_TEMPERATURE_SENSOR_ID", "8")
	os.Setenv("CONFIG_HOME_OFFICE_HUMIDITY_SENSOR_ID", "9")
	os.Setenv("CONFIG_FORECAST_TEMPERATURE_SENSOR_ID", "5")

	os.Setenv("TELEBOT_TOKEN", "deadbeef-beefdead")
	os.Setenv("TELEBOT_USERNAME", "JamesClonk")

	if err := exec.Command("cp", "-f", "_fixtures/test.db", "_fixtures/temp.db").Run(); err != nil {
		log.Fatal(err)
	}

	db := setupDatabase()
	n = setupNegroni(db)
}

func Test_Main_404(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:9090/something", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>Home Automation - Not Found</title>`)
	assert.Contains(t, body, `<section class="hero is-medium is-warning">`)
	assert.Contains(t, body, `<h1 class="title">404 - Not Found</h1>`)
	assert.Contains(t, body, `<h2 class="subtitle">This is not the page you are looking for ...</h2>`)
}

func Test_Main_500(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:9090/error", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>Home Automation - Error</title>`)
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
	assert.Contains(t, body, `<title>Home Automation</title>`)
	assert.Contains(t, body, `<h1 class="title">Home Automation</h1>`)
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
	assert.Contains(t, body, `<title>Home Automation - Forecasts - Bern</title>`)
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
	assert.Contains(t, body, `Unauthorized`)

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
	form.Add("symbol", "psi")
	form.Add("description", "atmospheric pressure")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusCreated, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "type": "pressure",
  "unit": "psi",
  "symbol": "psi",
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
    "symbol": "°",
    "description": "Shows temperature"`)
	assert.Contains(t, body, `
  {
    "id": 4,
    "type": "humidity",
    "unit": "percentage",
    "symbol": "%",
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
  "symbol": "psi",
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
	form.Add("symbol", "bar")
	form.Add("description", "atmosphärischer druck")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 5,
  "type": "druck",
  "unit": "bar",
  "symbol": "bar",
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
  "symbol": "bar",
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
	assert.Contains(t, body, `Unauthorized`)

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
  "id": 6,
  "name": "Badezimmer",
  "sensor_type": {
    "id": 3,
    "type": "temperature",
    "unit": "celsius",
    "symbol": "°",
    "description": "Shows temperature"
  },
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
    "sensor_type": {
      "id": 3,
      "type": "temperature",
      "unit": "celsius",
      "symbol": "°",
      "description": "Shows temperature"
    },
    "description": "Shows temperature"`)
	assert.Contains(t, body, `
    "id": 6,
    "name": "Badezimmer",
    "sensor_type": {
      "id": 3,
      "type": "temperature",
      "unit": "celsius",
      "symbol": "°",
      "description": "Shows temperature"
    },
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
  "sensor_type": {
    "id": 1,
    "type": "window_state",
    "unit": "closed",
    "symbol": "¬",
    "description": "Shows open/closed state of windows"
  },
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
  "sensor_type": {
    "id": 4,
    "type": "humidity",
    "unit": "percentage",
    "symbol": "%",
    "description": "Shows air humidity"
  },
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
  "sensor_type": {
    "id": 4,
    "type": "humidity",
    "unit": "percentage",
    "symbol": "%",
    "description": "Shows air humidity"
  },
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
	assert.Contains(t, body, `Unauthorized`)

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
    "timestamp": "1981-05-16T05:20:46Z",
    "value": 98
  },
  {
    "timestamp": "1980-12-15T17:17:16Z",
    "value": 16
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

func Test_Main_Alerts(t *testing.T) {
	// ----------------------- Unauthorized -----------------------
	response := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/alert", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `Unauthorized`)

	// ----------------------- CREATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/alert", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form := url.Values{}
	form.Add("sensor_id", "3")
	form.Add("name", "Alerteroni")
	form.Add("description", "Alert for fun!")
	form.Add("condition", "> 20")
	form.Add("execution", "15 * * * *")
	form.Add("silence_duration", "300")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusCreated, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 21,
  "sensor": {
    "id": 3,
    "name": "temperature #1",
    "sensor_type": {
      "id": 3,
      "type": "temperature",
      "unit": "celsius",
      "symbol": "°",
      "description": "Shows temperature"
    },
    "description": "Shows temperature"
  },
  "name": "Alerteroni",
  "description": "Alert for fun!",
  "alert_condition": "\u003e 20",
  "execution_schedule": "15 * * * *",`)
	assert.Contains(t, body, `"silence_duration_in_minutes": 300`)

	// ----------------------- READ -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/alert", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `
    "id": 5,
    "sensor": {
      "id": 4,
      "name": "humidity #1",
      "sensor_type": {
        "id": 4,
        "type": "humidity",
        "unit": "percentage",
        "symbol": "%",
        "description": "Shows air humidity"
      },
      "description": "Shows air humidity"
    },
    "name": "bedroom too cold",
    "description": "Alerts if bedroom temperature gets too cold",
    "alert_condition": "3;\u003c;13",
    "execution_schedule": "*/12 * * * *",`)
	assert.Contains(t, body, `"silence_duration_in_minutes": 720`)
	assert.Contains(t, body, `
    "id": 3,
    "sensor": {
      "id": 2,
      "name": "roof window #2",
      "sensor_type": {
        "id": 1,
        "type": "window_state",
        "unit": "closed",
        "symbol": "¬",
        "description": "Shows open/closed state of windows"
      },
      "description": "Shows open/closed state of roof window"
    },
    "name": "living room low humidity",
    "description": "Alerts if living room humidity gets too low",
    "alert_condition": "5;\u003c;30",
    "execution_schedule": "*/10 * * * *",`)
	assert.Contains(t, body, `"silence_duration_in_minutes": 300`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/alert/2", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 2,
  "sensor": {
    "id": 1,
    "name": "roof window #1",
    "sensor_type": {
      "id": 1,
      "type": "window_state",
      "unit": "closed",
      "symbol": "¬",
      "description": "Shows open/closed state of windows"
    },
    "description": "Shows open/closed state of roof window"
  },
  "name": "living room too hot",
  "description": "Alerts if living room temperature gets too hot",
  "alert_condition": "3;\u003e;30",
  "execution_schedule": "*/6 * * * *",`)
	assert.Contains(t, body, `"silence_duration_in_minutes": 300`)

	// ----------------------- UPDATE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("PUT", "/alert/2", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	form = url.Values{}
	form.Add("sensor_id", "1")
	form.Add("name", "Alerterei")
	form.Add("description", "Bad alert")
	form.Add("condition", "< 20")
	form.Add("execution", "30 * * * *")
	form.Add("silence_duration", "600")
	form.Add("last_alert", "2020-02-01T19:16:18Z")
	req.PostForm = form

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 2,
  "sensor": {
    "id": 1,
    "name": "roof window #1",
    "sensor_type": {
      "id": 1,
      "type": "window_state",
      "unit": "closed",
      "symbol": "¬",
      "description": "Shows open/closed state of windows"
    },
    "description": "Shows open/closed state of roof window"
  },
  "name": "Alerterei",
  "description": "Bad alert",
  "alert_condition": "\u003c 20",
  "execution_schedule": "30 * * * *",`)
	assert.Contains(t, body, `"silence_duration_in_minutes": 600`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/alert/2", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{
  "id": 2,
  "sensor": {
    "id": 1,
    "name": "roof window #1",
    "sensor_type": {
      "id": 1,
      "type": "window_state",
      "unit": "closed",
      "symbol": "¬",
      "description": "Shows open/closed state of windows"
    },
    "description": "Shows open/closed state of roof window"
  },
  "name": "Alerterei",
  "description": "Bad alert",
  "alert_condition": "\u003c 20",
  "execution_schedule": "30 * * * *",`)
	assert.Contains(t, body, `"silence_duration_in_minutes": 600`)

	// ----------------------- DELETE -----------------------
	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/alert/2", nil)
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
	req, err = http.NewRequest("GET", "/alert/2", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `"sql: no rows in result set"`)
}
