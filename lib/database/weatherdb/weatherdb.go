package weatherdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/anyandrea/weather_app/lib/database"
)

type WeatherDB interface {
	GetSensors() ([]*Sensor, error)
	GetSensor(id int) (*Sensor, error)
	GetSensorType(id int) (*SensorType, error)
	GetSensorTypes() ([]*SensorType, error)
	GetSensorData(id int, limit int) ([]*SensorData, error)
	GetSensorValues(id int, limit int) ([]*SensorValue, error)
	InsertSensorValue(sensorId, value int, timestamp time.Time) error
	GenerateSensorValues(id int, num int) error
	DropSensorValues(id int) error
}

type weatherDB struct {
	*sql.DB
	DatabaseType string
}

func NewWeatherDB(adapter database.Adapter) WeatherDB {
	return &weatherDB{adapter.GetDatabase(), adapter.GetType()}
}

func (wdb *weatherDB) GetSensors() ([]*Sensor, error) {
	rows, err := wdb.Query(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, s.unit, s.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		order by st.type asc, s.name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := []*Sensor{}
	for rows.Next() {
		var s Sensor
		if err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.TypeId, &s.Unit, &s.Description); err != nil {
			return nil, err
		}
		ss = append(ss, &s)
	}
	return ss, nil
}

func (wdb *weatherDB) GetSensor(id int) (*Sensor, error) {
	stmt, err := wdb.Prepare(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, s.unit, s.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.pk_sensor_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	s := &Sensor{}
	if err := stmt.QueryRow(id).Scan(&s.Id, &s.Name, &s.Type, &s.TypeId, &s.Unit, &s.Description); err != nil {
		return nil, err
	}
	return s, nil
}

func (wdb *weatherDB) GetSensorType(id int) (*SensorType, error) {
	stmt, err := wdb.Prepare(`
		select pk_sensor_type_id, type, description
		from sensor_type
		where pk_sensor_type_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := &SensorType{}
	if err := stmt.QueryRow(id).Scan(&t.Id, &t.Type, &t.Description); err != nil {
		return nil, err
	}
	return t, nil
}

func (wdb *weatherDB) GetSensorTypes() ([]*SensorType, error) {
	rows, err := wdb.Query(`
		select pk_sensor_type_id, type, description
		from sensor_type
		order by type asc, description asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	st := []*SensorType{}
	for rows.Next() {
		var t SensorType
		if err := rows.Scan(&t.Id, &t.Type, &t.Description); err != nil {
			return nil, err
		}
		st = append(st, &t)
	}
	return st, nil
}

func (wdb *weatherDB) GetSensorData(id, limit int) ([]*SensorData, error) {
	sql := `
		select s.pk_sensor_id, sd.timestamp, s.name, st.type, s.unit, sd.value
		from sensor_data sd
		join sensor s on s.pk_sensor_id = sd.fk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.pk_sensor_id = ?
		order by sd.timestamp desc`
	if limit > 0 {
		sql += fmt.Sprintf(" limit %d", limit)
	}

	stmt, err := wdb.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*SensorData{}
	for rows.Next() {
		var d SensorData
		if err := rows.Scan(&d.SensorId, &d.Timestamp, &d.Name, &d.Type, &d.Unit, &d.Value); err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

func (wdb *weatherDB) GetSensorValues(id, limit int) ([]*SensorValue, error) {
	sql := `
		select timestamp, value
		from sensor_data
		where fk_sensor_id = ?
		order by timestamp desc`
	if limit > 0 {
		sql += fmt.Sprintf(" limit %d", limit)
	}

	stmt, err := wdb.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []*SensorValue{}
	for rows.Next() {
		var value SensorValue
		if err := rows.Scan(&value.Timestamp, &value.Value); err != nil {
			return nil, err
		}
		values = append(values, &value)
	}
	return values, nil
}

func (wdb *weatherDB) InsertSensorValue(sensorId, value int, timestamp time.Time) error {
	stmt, err := wdb.Prepare(`
		insert into sensor_data (fk_sensor_id, value, timestamp) values (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorId, value, timestamp)
	return err
}

func (wdb *weatherDB) GenerateSensorValues(id, num int) error {
	sensor, err := wdb.GetSensor(id)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		value := rand.Intn(100)
		if sensor.Type == "state" {
			value = value % 2
		}
		timestamp := time.Unix(rand.Int63n(time.Now().Unix()-94608000)+94608000, 0)

		if err := wdb.InsertSensorValue(id, value, timestamp); err != nil {
			return err
		}
	}
	return nil
}

func (wdb *weatherDB) DropSensorValues(sensorId int) error {
	stmt, err := wdb.Prepare(`
		delete from sensor_data
		where fk_sensor_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorId)
	return err
}
