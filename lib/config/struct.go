package config

type Configuration struct {
	PlantRoom  Room
	LivingRoom Room
	BedRoom    Room
	Gallery    Room
	Basement   Room
	HomeOffice Room
	Forecast   struct {
		TemperatureSensorID int
		HumiditySensorID    int
		WindSpeedSensorID   int
	}
	Fitness struct {
		WeightID   int
		BodyFatID  int
		CaloriesID int
	}
}

type Room struct {
	TemperatureSensorID int
	HumiditySensorID    int
	CO2SensorID         int
	AirPressureSensorID int
	AltitudeSensorID    int
}
