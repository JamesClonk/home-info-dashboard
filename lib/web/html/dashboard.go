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
			Title:  "Home Automation - Dashboard",
			Active: "dashboard",
		}

		// collect the top column
		livingRoomTemp, err := hdb.GetSensorData(config.Get().LivingRoom.TemperatureSensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		livingRoomHum, err := hdb.GetSensorData(config.Get().LivingRoom.HumiditySensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomTemp, err := hdb.GetSensorData(config.Get().BedRoom.TemperatureSensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomHum, err := hdb.GetSensorData(config.Get().BedRoom.HumiditySensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
		homeOfficeTemp, err := hdb.GetSensorData(config.Get().HomeOffice.TemperatureSensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
		homeOfficeHum, err := hdb.GetSensorData(config.Get().HomeOffice.HumiditySensorID, 3)
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
		type Plant struct {
			Data *database.SensorData
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

		getAverage := func(data []*database.SensorData, rows int) int64 {
			var counter, value int64
			for r := 0; r < rows; r++ {
				if len(data) > r {
					value += data[r].Value
					counter += 1
				}
			}
			if counter == 0 {
				return 0
			}
			return value / counter
		}
		// average values because of multiple sensor for the same room
		if len(livingRoomHum) > 3 {
			livingRoomHum[0].Value = getAverage(livingRoomHum, 4)
		}
		if len(livingRoomTemp) > 3 {
			livingRoomTemp[0].Value = getAverage(livingRoomTemp, 4)
		}
		if len(bedRoomTemp) > 1 {
			bedRoomTemp[0].Value = getAverage(bedRoomTemp, 2)
		}
		if len(bedRoomHum) > 1 {
			bedRoomHum[0].Value = getAverage(bedRoomHum, 2)
		}
		if len(homeOfficeTemp) > 1 {
			homeOfficeTemp[0].Value = getAverage(homeOfficeTemp, 2)
		}
		if len(homeOfficeHum) > 1 {
			homeOfficeHum[0].Value = getAverage(homeOfficeHum, 2)
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

		plants := make([]Plant, 0)
		sensors, err := hdb.GetSensors()
		if err != nil {
			Error(rw, err)
			return
		}
		for _, sensor := range sensors {
			switch sensor.SensorType.Type {
			case "soil":
				d, err := hdb.GetSensorData(sensor.Id, 2)
				if err != nil {
					Error(rw, err)
					return
				}
				if len(d) > 1 {
					d[0].Value = getAverage(d, 2)
					plants = append(plants, Plant{
						Data: d[0],
					})
				}
			default:
				continue
			}
		}

		page.Content = struct {
			Rooms  []Room
			Plants []Plant
			Graphs Graphs
		}{
			Rooms:  rooms,
			Plants: plants,
			Graphs: graphs,
		}

		_ = web.Render().HTML(rw, http.StatusOK, "dashboard", page)
	}
}
