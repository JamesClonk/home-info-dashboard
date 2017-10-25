package weatherdb

import (
	"database/sql"

	"github.com/anyandrea/weather_app/lib/database"
)

type WeatherDB interface {
	GetSensors() ([]*Sensor, error)
	GetSensor(id int) (*Sensor, error)
}

type weatherDB struct {
	*sql.DB
	DatabaseType string
}

func NewWeatherDB(adapter database.Adapter) WeatherDB {
	return &weatherDB{adapter.GetDatabase(), adapter.GetType()}
}

func (wdb *weatherDB) GetSensors() ([]*Sensor, error) {
	rows, err := wdb.Query(`select pk_sensor_id, name, type, unit, description from sensor order by type asc, name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := []*Sensor{}
	for rows.Next() {
		var s Sensor
		if err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.Unit, &s.Description); err != nil {
			return nil, err
		}
		ss = append(ss, &s)
	}
	return ss, nil
}

func (wdb *weatherDB) GetSensor(id int) (*Sensor, error) {
	stmt, err := wdb.Prepare(`select pk_sensor_id, name, type, unit, description from sensor where pk_sensor_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	s := &Sensor{}
	if err := stmt.QueryRow(id).Scan(&s.Id, &s.Name, &s.Type, &s.Unit, &s.Description); err != nil {
		return nil, err
	}
	return s, nil
}
