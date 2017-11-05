package forecasts

import (
	"encoding/xml"
	"time"
)

type WeatherForecast struct {
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
				From   weatherDate `xml:"from,attr"`
				To     weatherDate `xml:"to,attr"`
				Period string      `xml:"period,attr"`
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

type weatherDate struct {
	time.Time
}

func (w *weatherDate) UnmarshalXMLAttr(attr xml.Attr) error {
	const format = "2006-01-02T15:04:05" // 2017-10-22T00:00:00
	parse, err := time.Parse(format, attr.Value)
	if err != nil {
		return err
	}
	*w = weatherDate{parse}
	return nil
}
