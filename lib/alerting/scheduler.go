package alerting

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
)

var (
	c = cron.New()
)

func StartScheduler() {
	c.Start()
}

func StopScheduler() {
	c.Stop()
}

func RegisterAlerts() {
	log.Println("registering monitoring alerts in scheduler")

	alerts, err := hdb.GetAlerts()
	if err != nil {
		log.Fatalf("could not read alert checks from database: %v\n", err)
	}

	for _, alert := range alerts {
		alertCopy := alert
		alertID := alertCopy.Id
		fn := func() { Check(alertID) }
		if err := c.AddFunc("0 "+alertCopy.Execution, fn); err != nil {
			log.Fatalf("could not register alert check [%s] in scheduler: %v\n", alertCopy.Name, err)
		}
		log.Printf("alert check [%s] with schedule [%s] and condition [%s] registered\n", alertCopy.Name, alertCopy.Execution, alertCopy.Condition)
	}
	StartScheduler()
}

func Check(alertID int) {
	// read alert
	alert, err := hdb.GetAlertById(alertID)
	if err != nil {
		message := fmt.Sprintf("alert check [ID:%d] failed: could not select from database: %v", alertID, err)
		log.Println(message)
		if err := Send(message); err != nil {
			log.Fatalf("could not send alert message: %v\n", err)
		}
		return
	}

	if alert.LastAlert != nil &&
		time.Since(*alert.LastAlert).Minutes() < float64(alert.SilenceDuration) { // too soon, not enough time has passed since last alert
		log.Printf("not running alert check [%s] because of silence duration [%dm]\n", alert.Name, alert.SilenceDuration)
		return
	}

	log.Printf("running alert check [%s]\n", alert.Name)

	// get condition
	params := strings.SplitN(alert.Condition, ";", 3)
	if len(params) != 3 {
		log.Printf("invalid alert condition: [%s]\n", alert.Condition)
		return
	}
	limit, err := strconv.Atoi(params[0])
	if err != nil {
		log.Printf("invalid value for alert average: [%s]\n", params[0])
		return
	}
	operator := params[1]
	target, err := strconv.Atoi(params[2])
	if err != nil {
		log.Printf("invalid value for alert target: [%s]\n", params[2])
		return
	}

	// get values
	data, err := hdb.GetSensorData(alert.Sensor.Id, limit)
	if err != nil {
		message := fmt.Sprintf("alert check [%s] failed: could not select data: %v", alert.Name, err)
		log.Println(message)
		if err := Send(message); err != nil {
			log.Fatalf("could not send alert message: %v\n", err)
		}
		return
	}
	value := Average(data)

	// check data against given conditions
	var match bool
	switch operator {
	case "<":
		if value < target {
			match = true
		}
	case ">":
		if value > target {
			match = true
		}
	case "==":
		if value == target {
			match = true
		}
	default: // always match/alert if none of the above operators are given
		match = true
	}
	if match {
		log.Printf("matched alert condition [%s] for sensor [%s], value: [%d]\n", alert.Condition, alert.Sensor.Name, value)
		if err := Send(fmt.Sprintf("Alert [*%s*] with Condition `%s` for Sensor [*%s*] triggered with a value of `%d`", alert.Name, alert.Condition, alert.Sensor.Name, value)); err != nil {
			log.Fatalf("could not send alert message: %v\n", err)
		}

		// update last alert timestamp
		now := time.Now()
		alert.LastAlert = &now
		if err := hdb.UpdateAlert(alert); err != nil {
			log.Printf("could not update alert status to database: %v\n", err)
		}
	}
}
