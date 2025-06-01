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
		sensorId, err := strconv.ParseInt(env.MustGet("CONFIG_LIVING_ROOM_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse living room temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.LivingRoom.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_LIVING_ROOM_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse living room humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.LivingRoom.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_LIVING_ROOM_CO2_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse living room co2 sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.LivingRoom.CO2SensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_PLANT_ROOM_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse plant room temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.PlantRoom.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_PLANT_ROOM_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse plant room humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.PlantRoom.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_BEDROOM_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse bedroom temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.BedRoom.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_BEDROOM_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse bedroom humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.BedRoom.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_BEDROOM_CO2_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse bedroom co2 sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.BedRoom.CO2SensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_GALLERY_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse gallery temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Gallery.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_GALLERY_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse gallery humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Gallery.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_GALLERY_AIR_PRESSURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse gallery air pressure sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Gallery.AirPressureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_BASEMENT_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse basement temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Basement.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_BASEMENT_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse basement humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Basement.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_HOME_OFFICE_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse home office temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.HomeOffice.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_HOME_OFFICE_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse home office humidity sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.HomeOffice.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_FORECAST_TEMPERATURE_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse forecast temperature sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Forecast.TemperatureSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_FORECAST_HUMIDITY_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse forecast humidty sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Forecast.HumiditySensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_FORECAST_WIND_SPEED_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse forecast wind speed sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Forecast.WindSpeedSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_FORECAST_RAIN_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse forecast rain sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Forecast.RainSensorID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_WEIGHT_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse weight sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Fitness.WeightID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_BODYFAT_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse weight sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Fitness.BodyFatID = int(sensorId)

		sensorId, err = strconv.ParseInt(env.MustGet("CONFIG_CALORIES_SENSOR_ID"), 10, 64)
		if err != nil {
			log.Printf("Could not parse calories sensor id: [%v]\n", sensorId)
			log.Fatal(err)
		}
		config.Fitness.CaloriesID = int(sensorId)
	})
	return config
}
