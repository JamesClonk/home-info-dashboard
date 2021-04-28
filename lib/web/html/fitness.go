package html

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/config"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
)

func Fitness(hdb database.HomeInfoDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := &Page{
			Title:  "Home Automation - Fitness",
			Active: "fitness",
		}

		if req.Method == "POST" {
			// parse the form and try to read the values from POST data
			_ = req.ParseForm()
			weight := req.Form.Get("weight")
			calories := req.Form.Get("calories")
			day := req.Form.Get("day")
			if (len(weight) > 0 || len(calories) > 0) && len(day) > 0 {
				timestamp := time.Now()
				if day == "yesterday" {
					timestamp = timestamp.Add(-24 * time.Hour)
				}
				if day == "the day before" {
					timestamp = timestamp.Add(-48 * time.Hour)
				}

				if len(weight) > 0 && weight != "0" {
					weightSensor, err := hdb.GetSensorById(config.Get().Fitness.WeightID)
					if err != nil {
						Error(rw, err)
						return
					}
					value, _ := strconv.ParseInt(weight, 10, 64)
					if value > 0 {
						if err := hdb.InsertSensorData(&database.SensorData{
							Sensor:    *weightSensor,
							Value:     value,
							Timestamp: &timestamp,
						}); err != nil {
							Error(rw, err)
							return
						}
					}
				}
				if len(calories) > 0 && calories != "0" {
					caloriesSensor, err := hdb.GetSensorById(config.Get().Fitness.CaloriesID)
					if err != nil {
						Error(rw, err)
						return
					}
					value, _ := strconv.ParseInt(calories, 10, 64)
					if value > 0 {
						if err := hdb.InsertSensorData(&database.SensorData{
							Sensor:    *caloriesSensor,
							Value:     value,
							Timestamp: &timestamp,
						}); err != nil {
							Error(rw, err)
							return
						}
					}
				}
				req.Method = http.MethodGet
				http.Redirect(rw, req, req.URL.RequestURI(), 303) // redirect to GET
				return
			}
		}

		// collect the graph data
		var graphLabels []string
		graphWeight := make(map[database.Sensor][]*database.SensorValue)
		graphCalories := make(map[database.Sensor][]*database.SensorValue)

		weightSensor, err := hdb.GetSensorById(config.Get().Fitness.WeightID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err := hdb.GetDailyAverages(weightSensor.Id, 99)
		if err != nil {
			Error(rw, err)
			return
		}
		graphWeight[*weightSensor] = values

		caloriesSensor, err := hdb.GetSensorById(config.Get().Fitness.CaloriesID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = hdb.GetDailySums(caloriesSensor.Id, 99)
		if err != nil {
			Error(rw, err)
			return
		}
		graphCalories[*caloriesSensor] = values

		// set labels & target weight
		for d := 0; d < 99; d++ {
			timestamp := time.Now().Add(-time.Duration(d) * 24 * time.Hour)
			graphLabels = append(graphLabels, timestamp.Format("02.01.2006"))

			// add target weight
			targetWeight := database.Sensor{
				Name:       "target",
				SensorType: weightSensor.SensorType,
			}
			graphWeight[targetWeight] = append(graphWeight[targetWeight], &database.SensorValue{
				Value:     74,
				Timestamp: &timestamp,
			})
		}

		// fill missing values
		var lastValue int64
		for _, label := range graphLabels {
			var found bool
			for _, value := range graphWeight[*weightSensor] {
				if value.Timestamp.Format("02.01.2006") == label {
					found = true
					lastValue = value.Value
				}
			}
			if !found {
				timestamp, _ := time.Parse("02.01.2006", label)
				graphWeight[*weightSensor] = append(graphWeight[*weightSensor], &database.SensorValue{
					Value:     lastValue,
					Timestamp: &timestamp,
				})
			}
			found = false
			for _, value := range graphCalories[*caloriesSensor] {
				if value.Timestamp.Format("02.01.2006") == label {
					found = true
				}
			}
			if !found {
				timestamp, _ := time.Parse("02.01.2006", label)
				graphCalories[*caloriesSensor] = append(graphCalories[*caloriesSensor], &database.SensorValue{
					Value:     0,
					Timestamp: &timestamp,
				})
			}
		}
		// sort by timestamp
		sort.Slice(graphWeight[*weightSensor][:], func(i, j int) bool {
			return graphWeight[*weightSensor][i].Timestamp.After(*graphWeight[*weightSensor][j].Timestamp)
		})
		sort.Slice(graphCalories[*caloriesSensor][:], func(i, j int) bool {
			return graphCalories[*caloriesSensor][i].Timestamp.After(*graphCalories[*caloriesSensor][j].Timestamp)
		})

		// collect the top column
		weight, err := hdb.GetSensorData(config.Get().Fitness.WeightID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		if len(weight) == 0 {
			timestamp := time.Now()
			weight = append(weight, &database.SensorData{
				Timestamp: &timestamp,
				Value:     1000,
				Sensor:    *weightSensor,
			})
		}
		calories, err := hdb.GetSensorData(config.Get().Fitness.CaloriesID, 1)
		if err != nil {
			Error(rw, err)
			return
		}
		if len(calories) == 0 {
			timestamp := time.Now()
			calories = append(calories, &database.SensorData{
				Timestamp: &timestamp,
				Value:     2000,
				Sensor:    *caloriesSensor,
			})
		}

		type Graphs struct {
			Labels   []string
			Weight   map[database.Sensor][]*database.SensorValue
			Calories map[database.Sensor][]*database.SensorValue
		}

		graphs := Graphs{
			Labels:   graphLabels,
			Weight:   graphWeight,
			Calories: graphCalories,
		}

		page.Content = struct {
			Graphs   Graphs
			Weight   *database.SensorData
			Calories *database.SensorData
		}{
			Graphs:   graphs,
			Weight:   weight[0],
			Calories: calories[0],
		}

		_ = web.Render().HTML(rw, http.StatusOK, "fitness", page)
	}
}
