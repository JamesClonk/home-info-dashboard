package util

import (
	"strings"

	"github.com/anyandrea/weather_app/lib/forecasts"
)

func GetWindowImage(state int64) (string, error) {
	var image string
	switch state {
	case 0:
		forecast, err := forecasts.Get(GetDefaultLocation("", ""))
		if err != nil {
			return "", err
		}
		weather := strings.ToLower(forecast.Forecast.Tabular.Time[0].Symbol.Name)

		image = "window_open_rainy.png"
		if strings.Contains(weather, "rain") || strings.Contains(weather, "shower") {
			image = "window_open_rainy.png"
		} else if strings.Contains(weather, "cloud") {
			image = "window_open_cloudy.png"
		} else if strings.Contains(weather, "snow") {
			image = "window_open_snowy.png"
		} else if strings.Contains(weather, "sun") || strings.Contains(weather, "fair") {
			image = "window_open_sunny.png"
		}
	case 1:
		image = "window_closed.png"
	}
	return image, nil
}
