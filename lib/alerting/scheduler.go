package alerting

import (
	"log"

	"github.com/JamesClonk/home-info-dashboard/lib/database"
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
		fn := func() { Check(*alertCopy) }
		if err := c.AddFunc(alertCopy.Execution, fn); err != nil {
			log.Fatalf("could not register alert check [%s] in scheduler: %v\n", alertCopy.Name, err)
		}
		log.Printf("alert check [%s] with schedule [%s] and condition [%s] registered\n", alertCopy.Name, alertCopy.Execution, alertCopy.Condition)
	}
	StartScheduler()
}

func Check(alert database.Alert) {
	log.Printf("running alert check [%s]\n", alert.Name)
	// if trigger, err := alert.Check(); err != nil {
	// 	log.Printf("scheduled alert check [%s] failed: %v\n", alert.Name, err)
	// }
}
