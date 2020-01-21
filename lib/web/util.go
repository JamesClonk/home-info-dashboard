package web

import (
	"net/http"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/util"
	"github.com/gorilla/mux"
)

func GetLocation(req *http.Request) (string, string) {
	// first, try to read values from gorilla mux
	vars := mux.Vars(req)
	canton := vars["canton"]
	city := vars["city"]

	// then, parse the form and try to read the values from POST data
	req.ParseForm()
	if len(canton) == 0 {
		canton = req.Form.Get("canton")
	}
	if len(city) == 0 {
		city = req.Form.Get("city")
	}

	return util.GetDefaultLocation(canton, city)
}

func SoilMoistureCapacitive(data []*database.SensorData) []*database.SensorData {
	for d := range data {
		data[d].Value = (data[d].Value - 333) * -1
		if data[d].Value < 0 {
			data[d].Value = 0
		} else if data[d].Value > 333 {
			data[d].Value = 333
		}
	}
	return data
}

func SoilMoisturePercentages(values []*database.SensorValue) []*database.SensorValue {
	for v := range values {
		values[v].Value = int64(float64((values[v].Value-333)*-1) / 3.33)
		if values[v].Value < 0 {
			values[v].Value = 0
		} else if values[v].Value > 100 {
			values[v].Value = 100
		}
	}
	return values
}
