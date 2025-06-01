package forecasts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	memo   map[string]WeatherForecast
	omMemo map[string]OpenMeteoWeatherForecast
	mutex  = &sync.Mutex{}
)

func init() {
	memo = make(map[string]WeatherForecast)
	omMemo = make(map[string]OpenMeteoWeatherForecast)
}

func getData(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	// identify ourselves for yr.no / api.met.no
	req.Header.Set("User-Agent", "home-info.jamesclonk.io github.com/JamesClonk/home-info-dashboard")

	resp, err := client.Do(req)
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
	if err := json.Unmarshal(data, &forecast); err != nil {
		return WeatherForecast{}, err
	}
	forecast.Timestamp = time.Now()
	return forecast, nil
}

func parseOpenMeteoWeatherForecast(data []byte) (OpenMeteoWeatherForecast, error) {
	var forecast OpenMeteoWeatherForecast
	if err := json.Unmarshal(data, &forecast); err != nil {
		return OpenMeteoWeatherForecast{}, err
	}
	forecast.Timestamp = time.Now()
	return forecast, nil
}

func Get(lat, lon, alt string) (data WeatherForecast, err error) {
	if len(lat) == 0 || len(lon) == 0 || len(alt) == 0 {
		// default
		lat = "47.02115"
		lon = "7.44914"
		alt = "555"
	}

	// yr.no - check memory first
	if forecast, ok := readMemo(lat, lon, alt); ok {
		// is it older than ~1 hour?
		if time.Now().After(forecast.Timestamp.Add(58 * time.Minute)) {
			if err := updateWeatherForecasts(lat, lon, alt); err != nil {
				return WeatherForecast{}, err
			}
		}
	} else {
		if err := updateWeatherForecasts(lat, lon, alt); err != nil {
			return WeatherForecast{}, err
		}
	}

	forecast, _ := readMemo(lat, lon, alt)
	return forecast, nil
}

func GetOpenMeteo(lat, lon string) OpenMeteoWeatherForecast {
	if len(lat) == 0 || len(lon) == 0 {
		// default
		lat = "47.02115"
		lon = "7.44914"
	}

	forecast, ok := readOmMemo(lat, lon)
	if !ok {
		return OpenMeteoWeatherForecast{}
	}
	return forecast
}

func readMemo(lat, lon, alt string) (WeatherForecast, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	value, ok := memo[fmt.Sprintf("%s:%s:%s", lat, lon, alt)]
	return value, ok
}

func readOmMemo(lat, lon string) (OpenMeteoWeatherForecast, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	value, ok := omMemo[fmt.Sprintf("%s:%s", lat, lon)]
	return value, ok
}

func updateWeatherForecasts(lat, lon, alt string) error {
	mutex.Lock()
	defer mutex.Unlock()

	// yr.no
	log.Printf("update yr.no weather forecast data for [lat:%s / lon:%s / alt:%s] ...\n", lat, lon, alt)
	data, err := getData(fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%s&lon=%s&altitude=%s", lat, lon, alt))
	if err != nil {
		return err
	}
	forecast, err := parseWeatherForecast(data)
	if err != nil {
		return err
	}
	memo[fmt.Sprintf("%s:%s:%s", lat, lon, alt)] = forecast

	// open-meteo.com
	log.Printf("update open-meteo.com weather forecast data for [lat:%s / lon:%s] ...\n", lat, lon)
	data, err = getData(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m,precipitation,precipitation_probability,relative_humidity_2m,rain,wind_speed_10m,showers&current=relative_humidity_2m,temperature_2m,precipitation,rain,wind_speed_10m&timezone=Europe%%2FBerlin&forecast_days=1", lat, lon))
	if err != nil {
		return err
	}
	omForecast, err := parseOpenMeteoWeatherForecast(data)
	if err != nil {
		return err
	}
	omMemo[fmt.Sprintf("%s:%s", lat, lon)] = omForecast

	return nil
}
