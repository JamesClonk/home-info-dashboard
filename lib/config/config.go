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
		// TODO: read config values from database (and not just once, but always anew)
		sensorId, err := strconv.ParseInt(env.MustGet("CONFIG_ROOM_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Println("Could not parse room temperature sensor id: [%v]", sensorId)
			log.Fatal(err)
		}
		config.Room.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_ROOM_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Println("Could not parse room humidity sensor id: [%v]", sensorId)
			log.Fatal(err)
		}
		config.Room.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_FORECAST_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Println("Could not parse forecast temperature sensor id: [%v]", sensorId)
			log.Fatal(err)
		}
		config.Forecast.TemperatureSensorID = int(sensorId)
	})
	return config
}
