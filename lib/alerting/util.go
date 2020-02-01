package alerting

import (
	"github.com/JamesClonk/home-info-dashboard/lib/database"
)

func Average(data []*database.SensorData) int {
	if len(data) == 0 {
		return 0
	}

	var value int
	for _, d := range data {
		value += int(d.Value)
	}
	return value / len(data)
}
