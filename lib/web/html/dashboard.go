package html

import (
	"net/http"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/config"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func Dashboard(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Info - Dashboard",
			Active: "dashboard",
		}

		// collect the top column
		roomTemp, err := hdb.GetSensorData(config.Get().Room.TemperatureSensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		roomHum, err := hdb.GetSensorData(config.Get().Room.HumiditySensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}

		// collect the graph data
		var graphLabels []string
		graphTemperature := make(map[database.Sensor][]*database.SensorValue)
		graphHumidity := make(map[database.Sensor][]*database.SensorValue)

		roomTempSensor, err := hdb.GetSensorById(config.Get().Room.TemperatureSensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err := hdb.GetHourlyAverages(config.Get().Room.TemperatureSensorID, 72)
		if err != nil {
			Error(rw, err)
			return
		}
		graphTemperature[*roomTempSensor] = values

		forecastTempSensor, err := hdb.GetSensorById(config.Get().Forecast.TemperatureSensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = hdb.GetHourlyAverages(config.Get().Forecast.TemperatureSensorID, 72)
		if err != nil {
			Error(rw, err)
			return
		}
		graphTemperature[*forecastTempSensor] = values

		roomHumSensor, err := hdb.GetSensorById(config.Get().Room.HumiditySensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = hdb.GetHourlyAverages(config.Get().Room.HumiditySensorID, 72)
		if err != nil {
			Error(rw, err)
			return
		}
		graphHumidity[*roomHumSensor] = values

		// graph labels
		for _, value := range values {
			graphLabels = append(graphLabels, value.Timestamp.Format("02.01. - 15:04"))
		}

		type Room struct {
			Temperature  *database.SensorData
			Humidity     *database.SensorData
			OutdatedTemp bool
			OutdatedHum  bool
		}
		type Graphs struct {
			Labels      []string
			Humidity    map[database.Sensor][]*database.SensorValue
			Temperature map[database.Sensor][]*database.SensorValue
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
		}

		page.Content = struct {
			Room   Room
			Graphs Graphs
		}{
			Room:   room,
			Graphs: graphs,
		}

		web.Render().HTML(rw, http.StatusOK, "dashboard", page)
	}
}
