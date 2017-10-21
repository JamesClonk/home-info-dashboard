package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeatherForecast struct {
	Location struct {
		Name     string `xml:"name"`
		Country  string `xml:"country"`
		Location struct {
			Altitude  string `xml:"altitude,attr"`
			Latitude  string `xml:"latitude,attr"`
			Longitude string `xml:"longitude,attr"`
			GeoBaseID string `xml:"geobaseid,attr"`
		} `xml:"location"`
	} `xml:"location"`
}

func getData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}

func parseWeatherForecast(data []byte) (*WeatherForecast, error) {
	var forecast WeatherForecast
	if err := xml.Unmarshal(data, &forecast); err != nil {
		return nil, err
	}
	return &forecast, nil
}

func GetWeatherForecast(canton, city string) (*WeatherForecast, error) {
	if len(canton) == 0 {
		canton = "Bern"
	}
	if len(city) == 0 {
		city = "Bern"
	}

	data, err := getData(fmt.Sprintf("http://www.yr.no/place/Switzerland/%s/%s/forecast.xml", canton, city))
	if err != nil {
		return nil, err
	}
	return parseWeatherForecast(data)
}
