package weatherdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/anyandrea/weather_app/lib/database"
	"github.com/anyandrea/weather_app/lib/util"
)

type WeatherDB interface {
	GetTemperature() (int64, error)
	GetWindowStates() ([]*Window, error)
	//GetDoorStates() ([]*Door, error)
	GetSensors() ([]*Sensor, error)
	GetSensorById(int) (*Sensor, error)
	GetSensorByName(string) (*Sensor, error)
	GetSensorTypeById(int) (*SensorType, error)
	GetSensorTypeByType(string) (*SensorType, error)
	GetSensorTypes() ([]*SensorType, error)
	GetSensorData(int, int) ([]*SensorData, error)
	GetSensorValues(int, int) ([]*SensorValue, error)
	InsertSensor(*Sensor) error
	InsertSensorType(*SensorType) error
	InsertSensorValue(int, int, time.Time) error
	UpdateSensor(*Sensor) error
	UpdateSensorType(*SensorType) error
	DeleteSensor(int) error
	DeleteSensorType(int) error
	DeleteSensorValues(int) error
	GenerateSensorValues(int, int) error
}

type weatherDB struct {
	*sql.DB
	DatabaseType string
}

func NewWeatherDB(adapter database.Adapter) WeatherDB {
	return &weatherDB{adapter.GetDatabase(), adapter.GetType()}
}

func (wdb *weatherDB) GetTemperature() (int64, error) {
	rows, err := wdb.Query(`
		select sd.value
		from sensor_data sd
		join sensor s on s.pk_sensor_id = sd.fk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where st.type = 'temperature'
		order by s.name asc, sd.timestamp desc
		limit 1`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var temperature int64
	if rows.Next() {
		if err := rows.Scan(&temperature); err != nil {
			return 0, err
		}
	}
	return temperature, nil
}

func (wdb *weatherDB) GetWindowStates() ([]*Window, error) {
	sensors, err := wdb.GetSensors()
	if err != nil {
		return nil, err
	}

	var windows []*Window
	for _, sensor := range sensors {
		if sensor.Type == "window_state" {
			values, err := wdb.GetSensorValues(sensor.Id, 1)
			if err != nil {
				return nil, err
			}

			if len(values) > 0 {
				image, err := util.GetWindowImage(values[0].Value)
				if err != nil {
					return nil, err
				}
				windows = append(windows, &Window{Image: image, Title: sensor.Name})
			}
		}
	}
	return windows, nil
}

func (wdb *weatherDB) GetSensors() ([]*Sensor, error) {
	rows, err := wdb.Query(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
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

func (wdb *weatherDB) GetSensorById(id int) (*Sensor, error) {
	stmt, err := wdb.Prepare(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
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

func (wdb *weatherDB) GetSensorByName(name string) (*Sensor, error) {
	stmt, err := wdb.Prepare(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.name = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	s := &Sensor{}
	if err := stmt.QueryRow(name).Scan(&s.Id, &s.Name, &s.Type, &s.TypeId, &s.Unit, &s.Description); err != nil {
		return nil, err
	}
	return s, nil
}

func (wdb *weatherDB) InsertSensor(sensor *Sensor) (err error) {
	var sensorTypeId int64
	if len(sensor.TypeId) > 0 {
		sensorTypeId, err = strconv.ParseInt(sensor.TypeId, 10, 64)
		if err != nil {
			return err
		}
	}

	// figure out sensor type
	var sensorType *SensorType
	if sensorTypeId > 0 { // by id
		sensorType, err = wdb.GetSensorTypeById(int(sensorTypeId))
		if err != nil {
			return err
		}
	} else { // by type
		sensorType, err = wdb.GetSensorTypeByType(sensor.Type)
		if err != nil {
			return err
		}
	}

	stmt, err := wdb.Prepare(`
		insert into sensor (name, fk_sensor_type_id, description) values (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensor.Name, sensorType.Id, sensor.Description); err != nil {
		return err
	}

	// make sure to get new values
	sensorNew, err := wdb.GetSensorByName(sensor.Name)
	if err != nil {
		return err
	}
	sensor.Id = sensorNew.Id
	sensor.Name = sensorNew.Name
	sensor.Type = sensorNew.Type
	sensor.TypeId = sensorNew.TypeId
	sensor.Unit = sensorNew.Unit
	sensor.Description = sensorNew.Description

	return nil
}

func (wdb *weatherDB) UpdateSensor(sensor *Sensor) (err error) {
	var sensorTypeId int64
	if len(sensor.TypeId) > 0 {
		sensorTypeId, err = strconv.ParseInt(sensor.TypeId, 10, 64)
		if err != nil {
			return err
		}
	}

	// figure out sensor type
	var sensorType *SensorType
	if sensorTypeId > 0 { // by id
		sensorType, err = wdb.GetSensorTypeById(int(sensorTypeId))
		if err != nil {
			return err
		}
	} else { // by type
		sensorType, err = wdb.GetSensorTypeByType(sensor.Type)
		if err != nil {
			return err
		}
	}

	stmt, err := wdb.Prepare(`
		update sensor
		set name = ?,
		fk_sensor_type_id = ?,
		description = ?
		where pk_sensor_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensor.Name, sensorType.Id, sensor.Description, sensor.Id); err != nil {
		return err
	}

	// make sure to get updated values
	sensorNew, err := wdb.GetSensorById(sensor.Id)
	if err != nil {
		return err
	}
	sensor.Name = sensorNew.Name
	sensor.Type = sensorNew.Type
	sensor.TypeId = sensorNew.TypeId
	sensor.Unit = sensorNew.Unit
	sensor.Description = sensorNew.Description

	return nil
}

func (wdb *weatherDB) GetSensorTypeById(id int) (*SensorType, error) {
	stmt, err := wdb.Prepare(`
		select pk_sensor_type_id, type, unit, description
		from sensor_type
		where pk_sensor_type_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := &SensorType{}
	if err := stmt.QueryRow(id).Scan(&t.Id, &t.Type, &t.Unit, &t.Description); err != nil {
		return nil, err
	}
	return t, nil
}

func (wdb *weatherDB) GetSensorTypeByType(s string) (*SensorType, error) {
	stmt, err := wdb.Prepare(`
		select pk_sensor_type_id, type, unit, description
		from sensor_type
		where type = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := &SensorType{}
	if err := stmt.QueryRow(s).Scan(&t.Id, &t.Type, &t.Unit, &t.Description); err != nil {
		return nil, err
	}
	return t, nil
}

func (wdb *weatherDB) GetSensorTypes() ([]*SensorType, error) {
	rows, err := wdb.Query(`
		select pk_sensor_type_id, type, unit, description
		from sensor_type
		order by type asc, description asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	st := []*SensorType{}
	for rows.Next() {
		var t SensorType
		if err := rows.Scan(&t.Id, &t.Type, &t.Unit, &t.Description); err != nil {
			return nil, err
		}
		st = append(st, &t)
	}
	return st, nil
}

func (wdb *weatherDB) InsertSensorType(sensorType *SensorType) (err error) {
	stmt, err := wdb.Prepare(`
		insert into sensor_type (type, unit, description) values (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensorType.Type, sensorType.Unit, sensorType.Description); err != nil {
		return err
	}

	// make sure to get new values
	sensorTypeNew, err := wdb.GetSensorTypeByType(sensorType.Type)
	if err != nil {
		return err
	}
	sensorType.Id = sensorTypeNew.Id
	sensorType.Type = sensorTypeNew.Type
	sensorType.Unit = sensorTypeNew.Unit
	sensorType.Description = sensorTypeNew.Description

	return nil
}

func (wdb *weatherDB) UpdateSensorType(sensorType *SensorType) (err error) {
	stmt, err := wdb.Prepare(`
		update sensor_type
		set type = ?,
		unit = ?,
		description = ?
		where pk_sensor_type_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensorType.Type, sensorType.Unit, sensorType.Description, sensorType.Id); err != nil {
		return err
	}

	// make sure to get updated values
	sensorTypeNew, err := wdb.GetSensorTypeById(sensorType.Id)
	if err != nil {
		return err
	}
	sensorType.Type = sensorTypeNew.Type
	sensorType.Unit = sensorTypeNew.Unit
	sensorType.Description = sensorTypeNew.Description

	return nil
}

func (wdb *weatherDB) GetSensorData(id, limit int) ([]*SensorData, error) {
	sql := `
		select s.pk_sensor_id, sd.timestamp, s.name, st.type, st.unit, sd.value
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

func (wdb *weatherDB) DeleteSensor(sensorId int) error {
	stmt, err := wdb.Prepare(`
		delete from sensor
		where pk_sensor_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorId)
	return err
}

func (wdb *weatherDB) DeleteSensorType(sensorTypeId int) error {
	stmt, err := wdb.Prepare(`
		delete from sensor_type
		where pk_sensor_type_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorTypeId)
	return err
}

func (wdb *weatherDB) DeleteSensorValues(sensorId int) error {
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

func (wdb *weatherDB) GenerateSensorValues(id, num int) error {
	sensor, err := wdb.GetSensorById(id)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		value := rand.Intn(100)
		if strings.Contains(sensor.Type, "state") {
			value = value % 2
		}
		timestamp := time.Unix(rand.Int63n(time.Now().Unix()-94608000)+94608000, 0)

		if err := wdb.InsertSensorValue(id, value, timestamp); err != nil {
			return err
		}
	}
	return nil
}
