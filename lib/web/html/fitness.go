package html

import (
	"net/http"
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
		values, err = hdb.GetDailyAverages(caloriesSensor.Id, 99)
		if err != nil {
			Error(rw, err)
			return
		}
		graphCalories[*caloriesSensor] = values

		// collect labels
		for _, value := range values {
			graphLabels = append(graphLabels, value.Timestamp.Format("02.01.2006"))
		}

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
