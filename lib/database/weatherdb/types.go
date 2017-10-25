package weatherdb

import "time"

type Sensor struct {
	Id          int    `json:"id" xml:"id,attr"`
	Name        string `json:"name" xml:"name"`
	Type        string `json:"type" xml:"type"`
	Unit        string `json:"unit" xml:"unit"`
	Description string `json:"description" xml:"description"`
}

type SensorValue struct {
	SensorId  int        `json:"sensor_id" xml:"sensor_id,attr"`
	Timestamp *time.Time `json:"timestamp" xml:"timestamp"`
	Name      string     `json:"name" xml:"name"`
	Type      string     `json:"type" xml:"type"`
	Unit      string     `json:"unit" xml:"unit"`
	Value     int64      `json:"value" xml:"value"`
}
