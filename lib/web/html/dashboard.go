package html

import (
	"net/http"
	"time"

	"github.com/anyandrea/weather_app/lib/config"
	"github.com/anyandrea/weather_app/lib/database/weatherdb"
	"github.com/anyandrea/weather_app/lib/web"
)

func Dashboard(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Weather App - Dashboard",
			Active: "dashboard",
		}

		// collect the top column
		roomTemp, err := wdb.GetSensorData(config.Get().Room.TemperatureSensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		roomHum, err := wdb.GetSensorData(config.Get().Room.HumiditySensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}

		// collect the graph data
		var graphLabels []string
		graphTemperature := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)
		graphHumidity := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)
		graphWindows := make(map[weatherdb.Sensor][]*weatherdb.SensorValue)

		roomTempSensor, err := wdb.GetSensorById(config.Get().Room.TemperatureSensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err := wdb.GetHourlyAverages(config.Get().Room.TemperatureSensorID, 72)
		if err != nil {
			Error(rw, err)
			return
		}
		graphTemperature[*roomTempSensor] = values

		forecastTempSensor, err := wdb.GetSensorById(config.Get().Forecast.TemperatureSensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = wdb.GetHourlyAverages(config.Get().Forecast.TemperatureSensorID, 72)
		if err != nil {
			Error(rw, err)
			return
		}
		graphTemperature[*forecastTempSensor] = values

		roomHumSensor, err := wdb.GetSensorById(config.Get().Room.HumiditySensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = wdb.GetHourlyAverages(config.Get().Room.HumiditySensorID, 72)
		if err != nil {
			Error(rw, err)
			return
		}
		graphHumidity[*roomHumSensor] = values

		// graph labels
		for _, value := range values {
			graphLabels = append(graphLabels, value.Timestamp.Format("02.01. - 15:04"))
		}

		// collect the window states
		windows, err := wdb.GetWindowStates()
		if err != nil {
			Error(rw, err)
			return
		}

		// collect window value changes
		windowSensorType, err := wdb.GetSensorTypeByType("window_state")
		if err != nil {
			Error(rw, err)
			return
		}
		windowSensors, err := wdb.GetSensorsByTypeId(windowSensorType.Id)
		if err != nil {
			Error(rw, err)
			return
		}
		for _, sensor := range windowSensors {
			values, err = wdb.GetSensorValues(sensor.Id, 14)
			if err != nil {
				Error(rw, err)
				return
			}
			graphWindows[*sensor] = values
		}

		type Room struct {
			Temperature  *weatherdb.SensorData
			Humidity     *weatherdb.SensorData
			OutdatedTemp bool
			OutdatedHum  bool
		}
		type Graphs struct {
			Labels      []string
			Humidity    map[weatherdb.Sensor][]*weatherdb.SensorValue
			Temperature map[weatherdb.Sensor][]*weatherdb.SensorValue
			Windows     map[weatherdb.Sensor][]*weatherdb.SensorValue
		}

		timeLimit := time.Now().Add(2 * time.Hour).UTC()

		room := Room{
			Temperature:  roomTemp[0],
			Humidity:     roomHum[0],
			OutdatedTemp: roomTemp[0].Timestamp.Before(timeLimit),
			OutdatedHum:  roomHum[0].Timestamp.Before(timeLimit),
		}
		graphs := Graphs{
			Labels:      graphLabels,
			Humidity:    graphHumidity,
			Temperature: graphTemperature,
			Windows:     graphWindows,
		}

		page.Content = struct {
			Windows []*weatherdb.Window
			Room    Room
			Graphs  Graphs
		}{
			Windows: windows,
			Room:    room,
			Graphs:  graphs,
		}

		web.Render().HTML(rw, http.StatusOK, "dashboard", page)
	}
}
