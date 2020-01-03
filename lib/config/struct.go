package config

type Configuration struct {
	LivingRoom Room
	BedRoom    Room
	HomeOffice Room
	Forecast   struct {
		TemperatureSensorID int
	}
}

type Room struct {
	TemperatureSensorID int
	HumiditySensorID    int
}
