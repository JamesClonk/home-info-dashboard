package forecasts

import (
	"encoding/xml"
	"time"
)

type OpenMeteoWeatherForecast struct {
	Timestamp time.Time `json:"time"`
	Timezone  string    `json:"timezone"`

	CurrentUnits struct {
		Time             string `json:"time"`
		Interval         string `json:"interval"`
		RelativeHumidity string `json:"relative_humidity_2m"`
		Temperature      string `json:"temperature_2m"`
		Precipitation    string `json:"precipitation"`
		Rain             string `json:"rain"`
		WindSpeed        string `json:"wind_speed_10m"`
	} `json:"current_units"`

	CurrentValues struct {
		Time             string  `json:"time"`
		Interval         int     `json:"interval"`
		RelativeHumidity float64 `json:"relative_humidity_2m"`
		Temperature      float64 `json:"temperature_2m"`
		Precipitation    float64 `json:"precipitation"`
		Rain             float64 `json:"rain"`
		WindSpeed        float64 `json:"wind_speed_10m"`
	} `json:"current"`

	HourlyUnits struct {
		Time                     string `json:"time"`
		RelativeHumidity         string `json:"relative_humidity_2m"`
		Temperature              string `json:"temperature_2m"`
		Precipitation            string `json:"precipitation"`
		PrecipitationProbability string `json:"precipitation_probability"`
		Rain                     string `json:"rain"`
		WindSpeed                string `json:"wind_speed_10m"`
		Showers                  string `json:"showers"`
	} `json:"hourly_units"`

	HourlyValues struct {
		Time                     []string  `json:"time"`
		RelativeHumidity         []float64 `json:"relative_humidity_2m"`
		Temperature              []float64 `json:"temperature_2m"`
		Precipitation            []float64 `json:"precipitation"`
		PrecipitationProbability []float64 `json:"precipitation_probability"`
		Rain                     []float64 `json:"rain"`
		WindSpeed                []float64 `json:"wind_speed_10m"`
		Showers                  []float64 `json:"showers"`
	} `json:"hourly"`
}

type WeatherForecast struct {
	Timestamp time.Time `json:"time"`

	Location struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`

	Properties struct {
		Meta struct {
			UpdatedAt time.Time `json:"updated_at"`
			Units     struct {
				AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
				AirTemperature        string `json:"air_temperature"`
				CloudAreaFraction     string `json:"cloud_area_fraction"`
				PrecipitationAmount   string `json:"precipitation_amount"`
				RelativeHumidity      string `json:"relative_humidity"`
				WindFromDirection     string `json:"wind_from_direction"`
				WindSpeed             string `json:"wind_speed"`
			} `json:"units"`
		} `json:"meta"`
		Timeseries []struct {
			Time time.Time `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
						AirTemperature        float64 `json:"air_temperature"`
						CloudAreaFraction     float64 `json:"cloud_area_fraction"`
						RelativeHumidity      float64 `json:"relative_humidity"`
						WindFromDirection     float64 `json:"wind_from_direction"`
						WindSpeed             float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next1Hour struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_1_hours"`
				Next6Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_6_hours"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

type OldXMLWeatherForecast struct {
	Timestamp time.Time `xml:"timestamp"`
	Location  struct {
		Name     string `xml:"name"`
		Country  string `xml:"country"`
		Location struct {
			Altitude  string `xml:"altitude,attr"`
			Latitude  string `xml:"latitude,attr"`
			Longitude string `xml:"longitude,attr"`
			GeoBaseID string `xml:"geobaseid,attr"`
		} `xml:"location"`
	} `xml:"location"`
	Credit struct {
		Link struct {
			Text string `xml:"text,attr"`
			URL  string `xml:"url,attr"`
		} `xml:"link"`
	} `xml:"credit"`
	Forecast struct {
		Tabular struct {
			Time []struct {
				From   oldXMLweatherDate `xml:"from,attr"`
				To     oldXMLweatherDate `xml:"to,attr"`
				Period string            `xml:"period,attr"`
				Symbol struct {
					Name string `xml:"name,attr"`
				} `xml:"symbol"`
				Precipitation struct {
					Value string `xml:"value,attr"`
				} `xml:"precipitation"`
				WindDirection struct {
					Degree string `xml:"deg,attr"`
					Code   string `xml:"code,attr"`
					Name   string `xml:"name,attr"`
				} `xml:"windDirection"`
				WindSpeed struct {
					MPS  string `xml:"mps,attr"`
					Name string `xml:"name,attr"`
				} `xml:"windSpeed"`
				Temperature struct {
					Unit  string `xml:"unit,attr"`
					Value string `xml:"value,attr"`
				} `xml:"temperature"`
				Pressure struct {
					Unit  string `xml:"unit,attr"`
					Value string `xml:"value,attr"`
				} `xml:"pressure"`
			} `xml:"time"`
		} `xml:"tabular"`
	} `xml:"forecast"`
}

type oldXMLweatherDate struct {
	time.Time
}

func (w *oldXMLweatherDate) UnmarshalXMLAttr(attr xml.Attr) error {
	const format = "2006-01-02T15:04:05" // 2017-10-22T00:00:00
	parse, err := time.Parse(format, attr.Value)
	if err != nil {
		return err
	}
	*w = oldXMLweatherDate{parse}
	return nil
}
