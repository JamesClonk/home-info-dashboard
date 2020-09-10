package alerting

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron"
)

var (
	c             = cron.New()
	alertingError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "homeinfo_dashboard_alerting_errors_total",
		Help: "Total number of Home-Info Dashboard alerting errors.",
	})
	alertsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "homeinfo_dashboard_alerts_total",
		Help: "Total number of Home-Info Dashboard alerts triggered.",
	})
	alertRulesRegistered = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "homeinfo_dashboard_registered_alert_rules_total",
		Help: "Total number of Home-Info Dashboard alert rules.",
	})
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
		alertRulesRegistered.Inc()
	}
	StartScheduler()
}

func Check(alertID int) {
	time.Sleep(time.Duration(rand.Intn(66)) * time.Second)

	// read alert
	alert, err := hdb.GetAlertById(alertID)
	if err != nil {
		message := fmt.Sprintf("alert check [ID:%d] failed: could not select from database: %v", alertID, err)
		log.Println(message)
		if err := Send(message); err != nil {
			log.Printf("could not send alert message: %v\n", err)
			alertingError.Inc()
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
		alertingError.Inc()
		return
	}
	limit, err := strconv.Atoi(params[0])
	if err != nil {
		log.Printf("invalid value for alert average: [%s]\n", params[0])
		alertingError.Inc()
		return
	}
	operator := params[1]
	target, err := strconv.Atoi(params[2])
	if err != nil {
		log.Printf("invalid value for alert target: [%s]\n", params[2])
		alertingError.Inc()
		return
	}

	// deadman check - if no metrics from the last 6h are found it's alerting time baby!
	hours := 6
	count, err := hdb.NumOfSensorDataWithinLastHours(alert.Sensor.Id, hours)
	if err != nil {
		message := fmt.Sprintf("deadman check for alert [%s] failed: could not select data: %v", alert.Name, err)
		log.Println(message)
		if err := Send(message); err != nil {
			log.Printf("could not send alert message: %v\n", err)
			alertingError.Inc()
		}
		return
	}
	if count <= 0 {
		message := fmt.Sprintf("deadman check for alert [%s] found no data within the last [%d] hours", alert.Name, hours)
		log.Println(message)
		if err := Send(message); err != nil {
			log.Printf("could not send alert message: %v\n", err)
			alertingError.Inc()
		}
		return
	}

	// get values
	data, err := hdb.GetSensorData(alert.Sensor.Id, limit)
	if err != nil {
		message := fmt.Sprintf("alert check [%s] failed: could not select data: %v", alert.Name, err)
		log.Println(message)
		if err := Send(message); err != nil {
			log.Printf("could not send alert message: %v\n", err)
			alertingError.Inc()
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
			log.Printf("could not send alert message: %v\n", err)
			alertingError.Inc()
		}
		alertsTotal.Inc()

		// update last alert timestamp
		now := time.Now()
		alert.LastAlert = &now
		if err := hdb.UpdateAlert(alert); err != nil {
			log.Printf("could not update alert status to database: %v\n", err)
			alertingError.Inc()
		}
	}
}
