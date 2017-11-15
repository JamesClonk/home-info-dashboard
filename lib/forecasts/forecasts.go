package forecasts

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var memo map[string]WeatherForecast

func init() {
	memo = make(map[string]WeatherForecast, 0)
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

func parseWeatherForecast(data []byte) (WeatherForecast, error) {
	var forecast WeatherForecast
	if err := xml.Unmarshal(data, &forecast); err != nil {
		return WeatherForecast{}, err
	}
	forecast.Timestamp = time.Now()
	return forecast, nil
}

func Get(canton, city string) (data WeatherForecast, err error) {
	if len(canton) == 0 {
		canton = "Bern"
	}
	if len(city) == 0 {
		city = "Bern"
	}

	// check memory first
	if forecast, ok := memo[fmt.Sprintf("%s:%s", canton, city)]; ok {
		// is it older than ~1 hour?
		if time.Now().After(forecast.Timestamp.Add(58 * time.Minute)) {
			if err := updateWeatherForecast(canton, city); err != nil {
				return WeatherForecast{}, err
			}
		}
	} else {
		if err := updateWeatherForecast(canton, city); err != nil {
			return WeatherForecast{}, err
		}
	}
	return memo[fmt.Sprintf("%s:%s", canton, city)], nil
}

func updateWeatherForecast(canton, city string) error {
	log.Printf("update weather forecast data for [%s/%s] ...\n", canton, city)

	data, err := getData(fmt.Sprintf("http://www.yr.no/place/Switzerland/%s/%s/forecast_hour_by_hour.xml", canton, city))
	//data, err := getData(fmt.Sprintf("http://www.yr.no/place/Switzerland/%s/%s/forecast.xml", canton, city))
	if err != nil {
		return err
	}

	forecast, err := parseWeatherForecast(data)
	if err != nil {
		return err
	}

	memo[fmt.Sprintf("%s:%s", canton, city)] = forecast
	return nil
}
