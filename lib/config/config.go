package config

import (
	"log"
	"strconv"
	"sync"

	"github.com/JamesClonk/home-info-dashboard/lib/env"
)

var (
	config *Configuration = nil
	once   sync.Once
)

func Get() *Configuration {
	once.Do(func() {
		config = &Configuration{}
		sensorId, err := strconv.ParseInt(env.MustGet("CONFIG_ROOM_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse room temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Room.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_ROOM_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse room humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Room.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_FORECAST_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse forecast temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Forecast.TemperatureSensorID = int(sensorId)
	})
	return config
}
