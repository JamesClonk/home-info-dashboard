package database

import (
	"database/sql"
	"time"
)

type SensorType struct {
	Id          int    `json:"id" xml:"id,attr"`
	Type        string `json:"type" xml:"type"`
	Unit        string `json:"unit" xml:"unit"`
	Symbol      string `json:"symbol" xml:"symbol"`
	Description string `json:"description" xml:"description"`
}

type Sensor struct {
	Id          int        `json:"id" xml:"id,attr"`
	Name        string     `json:"name" xml:"name"`
	SensorType  SensorType `json:"sensor_type" xml:"sensor_type"`
	Description string     `json:"description" xml:"description"`
}

type SensorData struct {
	Sensor    Sensor     `json:"sensor" xml:"sensor"`
	Timestamp *time.Time `json:"timestamp" xml:"timestamp"`
	Value     int64      `json:"value" xml:"value"`
}

type SensorValue struct {
	Timestamp *time.Time `json:"timestamp" xml:"timestamp"`
	Value     int64      `json:"value" xml:"value"`
}

type Alert struct {
	Id              int          `json:"id" xml:"id,attr"`
	Sensor          Sensor       `json:"sensor" xml:"sensor"`
	Name            string       `json:"name" xml:"name"`
	Description     string       `json:"description" xml:"description"`
	Condition       string       `json:"alert_condition" xml:"alert_condition"`
	Execution       string       `json:"execution_schedule" xml:"execution_schedule"`
	LastAlert       *time.Time   `json:"last_alert_timestamp" xml:"last_alert_timestamp"`
	SilenceDuration int64        `json:"silence_duration_in_minutes" xml:"silence_duration_in_minutes"`
	Active          sql.NullBool `json:"active" xml:"active"`
}

type Queue struct {
	Name         string     `json:"name" xml:"name"`
	MessageCount int64      `json:"message_count" xml:"message_count"`
	LastMessage  *time.Time `json:"last_message_timestamp" xml:"last_message_timestamp"`
}

type Message struct {
	Id        int        `json:"id" xml:"id,attr"`
	Queue     string     `json:"queue" xml:"queue"`
	Message   string     `json:"message" xml:"message"`
	CreatedAt *time.Time `json:"created_at" xml:"created_at"`
}
