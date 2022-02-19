package config

type Configuration struct {
	PlantRoom  Room
	LivingRoom Room
	BedRoom    Room
	HomeOffice Room
	Forecast   struct {
		TemperatureSensorID int
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
}
