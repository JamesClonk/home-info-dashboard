package html

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/JamesClonk/home-info-dashboard/lib/config"
	"github.com/JamesClonk/home-info-dashboard/lib/database"
	"github.com/JamesClonk/home-info-dashboard/lib/web"
	"github.com/Knetic/govaluate"
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
			bodyfat := req.Form.Get("bodyfat")
			calories := req.Form.Get("calories")
			day := req.Form.Get("day")
			if (len(weight) > 0 || len(bodyfat) > 0 || len(calories) > 0) && len(day) > 0 {
				timestamp := time.Now()
				if day == "yesterday" {
					timestamp = timestamp.Add(-24 * time.Hour)
				}
				if day == "ereyesterday" {
					timestamp = timestamp.Add(-48 * time.Hour)
				}

				if len(weight) > 0 && weight != "0" {
					weightSensor, err := hdb.GetSensorById(config.Get().Fitness.WeightID)
					if err != nil {
						Error(rw, err)
						return
					}

					// evaluate possible expression
					expression, err := govaluate.NewEvaluableExpression(weight)
					if err != nil {
						Error(rw, err)
						return
					}
					result, err := expression.Evaluate(nil)
					if err != nil {
						Error(rw, err)
						return
					}
					value, _ := result.(float64)
					if value >= 66 && value <= 86 { // guard against implausible values, we don't want to be anorexic nor obese
						if err := hdb.InsertSensorData(&database.SensorData{
							Sensor:    *weightSensor,
							Value:     int64(value * 10),
							Timestamp: &timestamp,
						}); err != nil {
							Error(rw, err)
							return
						}
					} else {
						Error(rw, fmt.Errorf("Body weight [%v kg] exceeds bounds, what the heck is wrong with you?!", value))
						return
					}
				}
				if len(bodyfat) > 0 && bodyfat != "0" {
					bodyfatSensor, err := hdb.GetSensorById(config.Get().Fitness.BodyFatID)
					if err != nil {
						Error(rw, err)
						return
					}

					// evaluate possible expression
					expression, err := govaluate.NewEvaluableExpression(bodyfat)
					if err != nil {
						Error(rw, err)
						return
					}
					result, err := expression.Evaluate(nil)
					if err != nil {
						Error(rw, err)
						return
					}
					value, _ := result.(float64)
					if value >= 10 && value <= 26 { // guard against implausible values
						if err := hdb.InsertSensorData(&database.SensorData{
							Sensor:    *bodyfatSensor,
							Value:     int64(value * 10),
							Timestamp: &timestamp,
						}); err != nil {
							Error(rw, err)
							return
						}
					} else {
						Error(rw, fmt.Errorf("Body fat [%v%%] exceeds bounds, what the heck is wrong with you?!", value))
						return
					}
				}
				if len(calories) > 0 && calories != "0" {
					caloriesSensor, err := hdb.GetSensorById(config.Get().Fitness.CaloriesID)
					if err != nil {
						Error(rw, err)
						return
					}

					// evaluate possible expression
					expression, err := govaluate.NewEvaluableExpression(calories)
					if err != nil {
						Error(rw, err)
						return
					}
					result, err := expression.Evaluate(nil)
					if err != nil {
						Error(rw, err)
						return
					}
					value, _ := result.(float64)
					if value > 0 && value < 4444 { // guard against implausible values
						if err := hdb.InsertSensorData(&database.SensorData{
							Sensor:    *caloriesSensor,
							Value:     int64(value),
							Timestamp: &timestamp,
						}); err != nil {
							Error(rw, err)
							return
						}
					} else {
						Error(rw, fmt.Errorf("Calorie intake [%v kcal] exceeds bounds, what the heck is wrong with you?!", value))
						return
					}
				}
			}
			// finish POST and redirect to GET, so the user can safely reload the page
			req.Method = http.MethodGet
			http.Redirect(rw, req, req.URL.RequestURI(), 303) // redirect to GET
			return
		}

		// collect the graph data
		var graphLabels []string
		graphWeight := make(map[database.Sensor][]*database.SensorValue)
		graphBodyFat := make(map[database.Sensor][]*database.SensorValue)
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

		bodyfatSensor, err := hdb.GetSensorById(config.Get().Fitness.BodyFatID)
		if err != nil {
			Error(rw, err)
			return
		}
		values, err = hdb.GetDailyAverages(bodyfatSensor.Id, 99)
		if err != nil {
			Error(rw, err)
			return
		}
		graphBodyFat[*bodyfatSensor] = values

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

		// collect calorie timeline
		calorieIntake, err := hdb.GetTodaysData(caloriesSensor.Id)
		if err != nil {
			Error(rw, err)
			return
		}

		// set labels & target weight / body fat
		for d := 0; d < 99; d++ {
			timestamp := time.Now().Add(-time.Duration(d) * 24 * time.Hour)
			graphLabels = append(graphLabels, timestamp.Format("02.01.2006"))

			// add target weight
			targetWeight := database.Sensor{
				Name:       "target",
				SensorType: weightSensor.SensorType,
			}
			graphWeight[targetWeight] = append(graphWeight[targetWeight], &database.SensorValue{
				Value:     730,
				Timestamp: &timestamp,
			})

			// add target body fat
			targetBodyFat := database.Sensor{
				Name:       "target",
				SensorType: bodyfatSensor.SensorType,
			}
			graphBodyFat[targetBodyFat] = append(graphBodyFat[targetBodyFat], &database.SensorValue{
				Value:     190,
				Timestamp: &timestamp,
			})
		}

		// fill missing values
		var lastWeightValue, lastBodyFatValue int64
		for _, label := range graphLabels {
			var found bool
			for _, value := range graphWeight[*weightSensor] {
				if value.Timestamp.Format("02.01.2006") == label {
					found = true
					lastWeightValue = value.Value
				}
			}
			if !found {
				timestamp, _ := time.Parse("02.01.2006", label)
				graphWeight[*weightSensor] = append(graphWeight[*weightSensor], &database.SensorValue{
					Value:     lastWeightValue,
					Timestamp: &timestamp,
				})
			}

			found = false
			for _, value := range graphBodyFat[*bodyfatSensor] {
				if value.Timestamp.Format("02.01.2006") == label {
					found = true
					lastBodyFatValue = value.Value
				}
			}
			if !found {
				timestamp, _ := time.Parse("02.01.2006", label)
				graphBodyFat[*bodyfatSensor] = append(graphBodyFat[*bodyfatSensor], &database.SensorValue{
					Value:     lastBodyFatValue,
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
		sort.Slice(graphBodyFat[*bodyfatSensor][:], func(i, j int) bool {
			return graphBodyFat[*bodyfatSensor][i].Timestamp.After(*graphBodyFat[*bodyfatSensor][j].Timestamp)
		})
		sort.Slice(graphCalories[*caloriesSensor][:], func(i, j int) bool {
			return graphCalories[*caloriesSensor][i].Timestamp.After(*graphCalories[*caloriesSensor][j].Timestamp)
		})
		// sort ascending
		sort.Slice(calorieIntake[:], func(i, j int) bool {
			return calorieIntake[i].Timestamp.Before(*calorieIntake[j].Timestamp)
		})

		// calculate weekly calorie averages
		calorieAvgWeekly := make([]int, 0)
		daysValue := 0
		for i, day := range graphCalories[*caloriesSensor] {
			if i == 0 {
				continue // skip first day as its still incomplete, dont count it towards history
			}
			daysValue = daysValue + int(day.Value)
			if i%7 == 0 {
				calorieAvgWeekly = append(calorieAvgWeekly, int(daysValue/7))
				daysValue = 0
			}
			if i > 42 {
				break
			}
		}

		type Graphs struct {
			Labels   []string
			Weight   map[database.Sensor][]*database.SensorValue
			BodyFat  map[database.Sensor][]*database.SensorValue
			Calories map[database.Sensor][]*database.SensorValue
		}

		graphs := Graphs{
			Labels:   graphLabels,
			Weight:   graphWeight,
			BodyFat:  graphBodyFat,
			Calories: graphCalories,
		}

		page.Content = struct {
			Graphs           Graphs
			Weight           *database.SensorValue
			BodyFat          *database.SensorValue
			Calories         *database.SensorValue
			CalorieIntake    []*database.SensorValue
			CalorieAvgWeekly []int
		}{
			Graphs:           graphs,
			Weight:           graphWeight[*weightSensor][0],
			BodyFat:          graphBodyFat[*bodyfatSensor][0],
			Calories:         graphCalories[*caloriesSensor][0],
			CalorieIntake:    calorieIntake,
			CalorieAvgWeekly: calorieAvgWeekly,
		}
		_ = web.Render().HTML(rw, http.StatusOK, "fitness", page)
	}
}
