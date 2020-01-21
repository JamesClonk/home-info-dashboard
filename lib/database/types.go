package database

import "time"

type SensorType struct {
	Id          int    `json:"id" xml:"id,attr"`
	Type        string `json:"type" xml:"type"`
	Unit        string `json:"unit" xml:"unit"`
	Symbol      string `json:"symbol" xml:"symbol"`
	Description string `json:"description" xml:"description"`
}

type Sensor struct {
	Id          int    `json:"id" xml:"id,attr"`
	Name        string `json:"name" xml:"name"`
	Type        string `json:"type" xml:"type"`
	TypeId      string `json:"type_id" xml:"type_id,attr"`
	Unit        string `json:"unit" xml:"unit"`
	Symbol      string `json:"symbol" xml:"symbol"`
	Description string `json:"description" xml:"description"`
}

type SensorData struct {
	SensorId  int        `json:"sensor_id" xml:"sensor_id,attr"`
	Timestamp *time.Time `json:"timestamp" xml:"timestamp"`
	Name      string     `json:"name" xml:"name"`
	Type      string     `json:"type" xml:"type"`
	Unit      string     `json:"unit" xml:"unit"`
	Symbol    string     `json:"symbol" xml:"symbol"`
	Value     int64      `json:"value" xml:"value"`
}

type SensorValue struct {
	Timestamp *time.Time `json:"timestamp" xml:"timestamp"`
	Value     int64      `json:"value" xml:"value"`
}
