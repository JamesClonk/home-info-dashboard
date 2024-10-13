package html

import (
	"math"
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
		plantRoomTemp, err := hdb.GetSensorData(config.Get().PlantRoom.TemperatureSensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
		plantRoomHum, err := hdb.GetSensorData(config.Get().PlantRoom.HumiditySensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
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
		livingRoomCO2, err := hdb.GetSensorData(config.Get().LivingRoom.CO2SensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomTemp, err := hdb.GetSensorData(config.Get().BedRoom.TemperatureSensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomHum, err := hdb.GetSensorData(config.Get().BedRoom.HumiditySensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		bedRoomCO2, err := hdb.GetSensorData(config.Get().BedRoom.CO2SensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		galleryTemp, err := hdb.GetSensorData(config.Get().Gallery.TemperatureSensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		galleryHum, err := hdb.GetSensorData(config.Get().Gallery.HumiditySensorID, 5)
		if err != nil {
			Error(rw, err)
			return
		}
		galleryPa, err := hdb.GetSensorData(config.Get().Gallery.AirPressureSensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
		basementTemp, err := hdb.GetSensorData(config.Get().Basement.TemperatureSensorID, 3)
		if err != nil {
			Error(rw, err)
			return
		}
		basementHum, err := hdb.GetSensorData(config.Get().Basement.HumiditySensorID, 3)
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
			graphLabels = append(graphLabels, value.Timestamp.Format("02.01.2006 - 15:04"))
		}

		type Room struct {
			Temperature *database.SensorData
			Humidity    *database.SensorData
			CO2         *database.SensorData
			AirPressure *database.SensorData
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
			return int64(math.RoundToEven(float64(value) / float64(counter)))
		}
		// average values because of multiple sensor for the same room
		if len(plantRoomHum) >= 2 {
			plantRoomHum[0].Value = getAverage(plantRoomHum, 2)
		}
		if len(plantRoomTemp) >= 2 {
			plantRoomTemp[0].Value = getAverage(plantRoomTemp, 2)
		}
		if len(livingRoomHum) >= 3 {
			livingRoomHum[0].Value = getAverage(livingRoomHum, 3)
		}
		if len(livingRoomTemp) >= 3 {
			livingRoomTemp[0].Value = getAverage(livingRoomTemp, 3)
		}
		if len(livingRoomCO2) >= 3 {
			livingRoomCO2[0].Value = getAverage(livingRoomCO2, 3)
		}
		if len(bedRoomTemp) >= 3 {
			bedRoomTemp[0].Value = getAverage(bedRoomTemp, 3)
		}
		if len(bedRoomHum) >= 3 {
			bedRoomHum[0].Value = getAverage(bedRoomHum, 3)
		}
		if len(bedRoomCO2) >= 3 {
			bedRoomCO2[0].Value = getAverage(bedRoomCO2, 3)
		}
		if len(galleryTemp) >= 3 {
			galleryTemp[0].Value = getAverage(galleryTemp, 3)
		}
		if len(galleryHum) >= 3 {
			galleryHum[0].Value = getAverage(galleryHum, 3)
		}
		if len(galleryPa) >= 2 {
			galleryPa[0].Value = getAverage(galleryPa, 2)
		}
		if len(basementTemp) >= 2 {
			basementTemp[0].Value = getAverage(basementTemp, 2)
		}
		if len(basementHum) >= 2 {
			basementHum[0].Value = getAverage(basementHum, 2)
		}
		if len(homeOfficeTemp) >= 2 {
			homeOfficeTemp[0].Value = getAverage(homeOfficeTemp, 2)
		}
		if len(homeOfficeHum) >= 2 {
			homeOfficeHum[0].Value = getAverage(homeOfficeHum, 2)
		}

		rooms := make([]Room, 0)
		if len(plantRoomTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: plantRoomTemp[0],
				Humidity:    plantRoomHum[0],
			})
		}
		if len(livingRoomTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: livingRoomTemp[0],
				Humidity:    livingRoomHum[0],
				CO2:         livingRoomCO2[0],
			})
		}
		if len(bedRoomTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: bedRoomTemp[0],
				Humidity:    bedRoomHum[0],
				CO2:         bedRoomCO2[0],
			})
		}
		if len(galleryTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: galleryTemp[0],
				Humidity:    galleryHum[0],
				AirPressure: galleryPa[0],
			})
		}
		if len(basementTemp) > 0 {
			rooms = append(rooms, Room{
				Temperature: basementTemp[0],
				Humidity:    basementHum[0],
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
				d, err := hdb.GetSensorData(sensor.Id, 6)
				if err != nil {
					Error(rw, err)
					return
				}
				if len(d) > 1 {
					d[0].Value = getAverage(d, 5)
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
