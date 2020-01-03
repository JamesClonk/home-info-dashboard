package html

import (
	"net/http"

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
		livingRoomTemp, err := hdb.GetSensorData(config.Get().LivingRoom.TemperatureSensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		livingRoomHum, err := hdb.GetSensorData(config.Get().LivingRoom.HumiditySensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomTemp, err := hdb.GetSensorData(config.Get().BedRoom.TemperatureSensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomHum, err := hdb.GetSensorData(config.Get().BedRoom.HumiditySensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		homeOfficeTemp, err := hdb.GetSensorData(config.Get().HomeOffice.TemperatureSensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		homeOfficeHum, err := hdb.GetSensorData(config.Get().HomeOffice.HumiditySensorID, 1)
		if err != nil {
			Error(rw, err)
			return
		}

		// collect the graph data
		var graphLabels []string
		graphTemperature := make(map[database.Sensor][]*database.SensorValue)
		graphHumidity := make(map[database.Sensor][]*database.SensorValue)

		roomTempSensor, err := hdb.GetSensorById(config.Get().LivingRoom.TemperatureSensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err := hdb.GetHourlyAverages(config.Get().LivingRoom.TemperatureSensorID, 72)
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

		roomHumSensor, err := hdb.GetSensorById(config.Get().LivingRoom.HumiditySensorID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = hdb.GetHourlyAverages(config.Get().LivingRoom.HumiditySensorID, 72)
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
			Temperature *database.SensorData
			Humidity    *database.SensorData
		}
		type Graphs struct {
			Labels      []string
			Humidity    map[database.Sensor][]*database.SensorValue
			Temperature map[database.Sensor][]*database.SensorValue
		}

		graphs := Graphs{
			Labels:      graphLabels,
			Humidity:    graphHumidity,
			Temperature: graphTemperature,
		}

		rooms := make([]Room, 0)
		if len(livingRoomTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: livingRoomTemp[0],
				Humidity:    livingRoomHum[0],
			})
		}
		if len(bedRoomTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: bedRoomTemp[0],
				Humidity:    bedRoomHum[0],
			})
		}
		if len(homeOfficeTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: homeOfficeTemp[0],
				Humidity:    homeOfficeHum[0],
			})
		}

		page.Content = struct {
			Rooms  []Room
			Graphs Graphs
		}{
			Rooms:  rooms,
			Graphs: graphs,
		}

		_ = web.Render().HTML(rw, http.StatusOK, "dashboard", page)
	}
}
