package config

type Configuration struct {
	Room struct {
		TemperatureSensorID int
		HumiditySensorID    int
	}
	Forecast struct {
		TemperatureSensorID int
	}
}
