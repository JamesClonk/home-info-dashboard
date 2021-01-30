package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type HomeInfoDB interface {
	GetAllAlerts() ([]*Alert, error)
	GetActiveAlerts() ([]*Alert, error)
	GetAlertsByState(bool) ([]*Alert, error)
	GetAlertById(int) (*Alert, error)
	GetAlertsByName(string) ([]*Alert, error)
	GetSensors() ([]*Sensor, error)
	GetSensorById(int) (*Sensor, error)
	GetSensorsByName(string) ([]*Sensor, error)
	GetSensorsByTypeId(int) ([]*Sensor, error)
	GetSensorTypeById(int) (*SensorType, error)
	GetSensorTypeByType(string) (*SensorType, error)
	GetSensorTypes() ([]*SensorType, error)
	NumOfSensorDataWithinLastHours(int, int) (int, error)
	GetSensorData(int, int) ([]*SensorData, error)
	GetSensorValues(int, int) ([]*SensorValue, error)
	GetHourlyAverages(int, int) ([]*SensorValue, error)
	GetDailyAverages(int, int) ([]*SensorValue, error)
	InsertAlert(*Alert) error
	InsertSensor(*Sensor) error
	InsertSensorType(*SensorType) error
	InsertSensorValue(int, int, time.Time) error
	UpdateAlert(*Alert) error
	UpdateSensor(*Sensor) error
	UpdateSensorType(*SensorType) error
	DeleteAlert(int) error
	DeleteSensor(int) error
	DeleteSensorType(int) error
	DeleteSensorValues(int) error
	GenerateSensorValues(int, int) error
	Housekeeping() error
	GetQueues() ([]*Queue, error)
	GetQueueByName(string) (*Queue, error)
	GetMessages() ([]*Message, error)
	GetMessagesFromQueue(string) ([]*Message, error)
	GetMessageById(int) (*Message, error)
	InsertMessage(*Message) error
	DeleteQueue(string) error
	DeleteMessage(int) error
}

type homeInfoDB struct {
	*sql.DB
	DatabaseType string
}

func NewHomeInfoDB(adapter Adapter) HomeInfoDB {
	return &homeInfoDB{adapter.GetDatabase(), adapter.GetType()}
}

func (hdb *homeInfoDB) GetAllAlerts() ([]*Alert, error) {
	rows, err := hdb.Query(`
		select
			a.pk_alert_id, a.name, a.description, a.condition, a.execution, a.last_alert, a.silence_duration,
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from alert a
		join sensor s on a.fk_sensor_id = s.pk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		order by a.name asc, s.name asc, st.type asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := []*Alert{}
	for rows.Next() {
		var a Alert
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&a.Id, &a.Name, &a.Description, &a.Condition, &a.Execution, &a.LastAlert, &a.SilenceDuration,
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		a.Sensor = s
		aa = append(aa, &a)
	}
	return aa, nil
}

func (hdb *homeInfoDB) GetActiveAlerts() ([]*Alert, error) {
	return hdb.GetAlertsByState(true)
}

func (hdb *homeInfoDB) GetAlertsByState(active bool) ([]*Alert, error) {
	stmt, err := hdb.Prepare(`
		select
			a.pk_alert_id, a.name, a.description, a.condition, a.execution, a.last_alert, a.silence_duration,
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from alert a
		join sensor s on a.fk_sensor_id = s.pk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where a.active = $1
		order by a.name asc, s.name asc, st.type asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := []*Alert{}
	for rows.Next() {
		var a Alert
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&a.Id, &a.Name, &a.Description, &a.Condition, &a.Execution, &a.LastAlert, &a.SilenceDuration,
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		a.Sensor = s
		aa = append(aa, &a)
	}
	return aa, nil
}

func (hdb *homeInfoDB) GetAlertById(id int) (*Alert, error) {
	stmt, err := hdb.Prepare(`
		select
			a.pk_alert_id, a.name, a.description, a.condition, a.execution, a.last_alert, a.silence_duration,
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from alert a
		join sensor s on a.fk_sensor_id = s.pk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where a.pk_alert_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var a Alert
	var s Sensor
	var st SensorType
	if err := stmt.QueryRow(id).Scan(
		&a.Id, &a.Name, &a.Description, &a.Condition, &a.Execution, &a.LastAlert, &a.SilenceDuration,
		&s.Id, &s.Name, &s.Description,
		&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
	); err != nil {
		return nil, err
	}
	s.SensorType = st
	a.Sensor = s
	return &a, nil
}

func (hdb *homeInfoDB) GetAlertsByName(name string) ([]*Alert, error) {
	stmt, err := hdb.Prepare(`
		select
			a.pk_alert_id, a.name, a.description, a.condition, a.execution, a.last_alert, a.silence_duration,
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from alert a
		join sensor s on a.fk_sensor_id = s.pk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where a.name = $1
		order by a.pk_alert_id desc, s.name asc, st.type asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := []*Alert{}
	for rows.Next() {
		var a Alert
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&a.Id, &a.Name, &a.Description, &a.Condition, &a.Execution, &a.LastAlert, &a.SilenceDuration,
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		a.Sensor = s
		aa = append(aa, &a)
	}
	return aa, nil
}

func (hdb *homeInfoDB) GetSensors() ([]*Sensor, error) {
	rows, err := hdb.Query(`
		select
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		order by s.name asc, st.type asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := []*Sensor{}
	for rows.Next() {
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		ss = append(ss, &s)
	}
	return ss, nil
}

func (hdb *homeInfoDB) GetSensorsByTypeId(id int) ([]*Sensor, error) {
	stmt, err := hdb.Prepare(`
		select
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where st.pk_sensor_type_id = $1
		order by s.name asc, st.type asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := []*Sensor{}
	for rows.Next() {
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		ss = append(ss, &s)
	}
	return ss, nil
}

func (hdb *homeInfoDB) GetSensorById(id int) (*Sensor, error) {
	stmt, err := hdb.Prepare(`
		select
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.pk_sensor_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var s Sensor
	var st SensorType
	if err := stmt.QueryRow(id).Scan(
		&s.Id, &s.Name, &s.Description,
		&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
	); err != nil {
		return nil, err
	}
	s.SensorType = st
	return &s, nil
}

func (hdb *homeInfoDB) GetSensorsByName(name string) ([]*Sensor, error) {
	stmt, err := hdb.Prepare(`
		select
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.name = $1
		order by s.pk_sensor_id desc, st.type asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := []*Sensor{}
	for rows.Next() {
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		ss = append(ss, &s)
	}
	return ss, nil
}

func (hdb *homeInfoDB) InsertAlert(alert *Alert) (err error) {
	// require sensor id
	if alert.Sensor.Id <= 0 {
		return fmt.Errorf("sensor_id missing!")
	}

	stmt, err := hdb.Prepare(`
		insert into alert (name, fk_sensor_id, description, condition, execution, silence_duration) values ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(alert.Name, alert.Sensor.Id, alert.Description, alert.Condition, alert.Execution, alert.SilenceDuration); err != nil {
		return err
	}

	// make sure to get new values
	alerts, err := hdb.GetAlertsByName(alert.Name)
	if err != nil {
		return err
	}
	if len(alerts) == 0 {
		return fmt.Errorf("no alert found with name %s", alert.Name)
	}
	alertNew := alerts[0] // first one is newest
	alert.Id = alertNew.Id
	alert.Name = alertNew.Name
	alert.Sensor = alertNew.Sensor
	alert.LastAlert = alertNew.LastAlert

	return nil
}

func (hdb *homeInfoDB) InsertSensor(sensor *Sensor) (err error) {
	// figure out sensor type
	var sensorType SensorType
	if sensor.SensorType.Id > 0 {
		sensorType = sensor.SensorType
	} else if len(sensor.SensorType.Type) > 0 { // by type name
		st, err := hdb.GetSensorTypeByType(sensor.SensorType.Type)
		if err != nil {
			return err
		}
		sensorType = *st
	} else {
		return fmt.Errorf("sensor_type missing!")
	}

	stmt, err := hdb.Prepare(`
		insert into sensor (name, fk_sensor_type_id, description) values ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensor.Name, sensorType.Id, sensor.Description); err != nil {
		return err
	}

	// make sure to get new values
	sensors, err := hdb.GetSensorsByName(sensor.Name)
	if err != nil {
		return err
	}
	if len(sensors) == 0 {
		return fmt.Errorf("no sensor found with name %s", sensor.Name)
	}
	sensorNew := sensors[0] // first one is newest
	sensor.Id = sensorNew.Id
	sensor.Name = sensorNew.Name
	sensor.SensorType = sensorNew.SensorType
	sensor.Description = sensorNew.Description

	return nil
}

func (hdb *homeInfoDB) UpdateAlert(alert *Alert) (err error) {
	// require sensor id
	if alert.Sensor.Id <= 0 {
		return fmt.Errorf("sensor_id missing!")
	}

	stmt, err := hdb.Prepare(`
		update alert
		set name = $1,
		fk_sensor_id = $2,
		description = $3,
		condition = $4,
		execution = $5,
		silence_duration = $6,
		last_alert = $7
		where pk_alert_id = $8`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(alert.Name, alert.Sensor.Id, alert.Description, alert.Condition, alert.Execution, alert.SilenceDuration, alert.LastAlert, alert.Id); err != nil {
		return err
	}

	// make sure to get updated values
	alertNew, err := hdb.GetAlertById(alert.Id)
	if err != nil {
		return err
	}
	alert.Name = alertNew.Name
	alert.Id = alertNew.Id
	alert.Sensor = alertNew.Sensor
	alert.Description = alertNew.Description
	alert.Condition = alertNew.Condition
	alert.Execution = alertNew.Execution
	alert.LastAlert = alertNew.LastAlert
	alert.SilenceDuration = alertNew.SilenceDuration

	return nil
}

func (hdb *homeInfoDB) UpdateSensor(sensor *Sensor) (err error) {
	// figure out sensor type
	var sensorType SensorType
	if sensor.SensorType.Id > 0 {
		sensorType = sensor.SensorType
	} else if len(sensor.SensorType.Type) > 0 { // by type name
		st, err := hdb.GetSensorTypeByType(sensor.SensorType.Type)
		if err != nil {
			return err
		}
		sensorType = *st
	} else {
		return fmt.Errorf("sensor_type missing!")
	}

	stmt, err := hdb.Prepare(`
		update sensor
		set name = $1,
		fk_sensor_type_id = $2,
		description = $3
		where pk_sensor_id = $4`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensor.Name, sensorType.Id, sensor.Description, sensor.Id); err != nil {
		return err
	}

	// make sure to get updated values
	sensorNew, err := hdb.GetSensorById(sensor.Id)
	if err != nil {
		return err
	}
	sensor.Name = sensorNew.Name
	sensor.SensorType = sensorNew.SensorType
	sensor.Description = sensorNew.Description

	return nil
}

func (hdb *homeInfoDB) GetSensorTypeById(id int) (*SensorType, error) {
	stmt, err := hdb.Prepare(`
		select pk_sensor_type_id, type, unit, symbol, description
		from sensor_type
		where pk_sensor_type_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := &SensorType{}
	if err := stmt.QueryRow(id).Scan(&t.Id, &t.Type, &t.Unit, &t.Symbol, &t.Description); err != nil {
		return nil, err
	}
	return t, nil
}

func (hdb *homeInfoDB) GetSensorTypeByType(s string) (*SensorType, error) {
	stmt, err := hdb.Prepare(`
		select pk_sensor_type_id, type, unit, symbol, description
		from sensor_type
		where type = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := &SensorType{}
	if err := stmt.QueryRow(s).Scan(&t.Id, &t.Type, &t.Unit, &t.Symbol, &t.Description); err != nil {
		return nil, err
	}
	return t, nil
}

func (hdb *homeInfoDB) GetSensorTypes() ([]*SensorType, error) {
	rows, err := hdb.Query(`
		select pk_sensor_type_id, type, unit, symbol, description
		from sensor_type
		order by type asc, description asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	st := []*SensorType{}
	for rows.Next() {
		var t SensorType
		if err := rows.Scan(&t.Id, &t.Type, &t.Unit, &t.Symbol, &t.Description); err != nil {
			return nil, err
		}
		st = append(st, &t)
	}
	return st, nil
}

func (hdb *homeInfoDB) InsertSensorType(sensorType *SensorType) (err error) {
	stmt, err := hdb.Prepare(`
		insert into sensor_type (type, unit, symbol, description) values ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensorType.Type, sensorType.Unit, sensorType.Symbol, sensorType.Description); err != nil {
		return err
	}

	// make sure to get new values
	sensorTypeNew, err := hdb.GetSensorTypeByType(sensorType.Type)
	if err != nil {
		return err
	}
	sensorType.Id = sensorTypeNew.Id
	sensorType.Type = sensorTypeNew.Type
	sensorType.Unit = sensorTypeNew.Unit
	sensorType.Symbol = sensorTypeNew.Symbol
	sensorType.Description = sensorTypeNew.Description

	return nil
}

func (hdb *homeInfoDB) UpdateSensorType(sensorType *SensorType) (err error) {
	stmt, err := hdb.Prepare(`
		update sensor_type
		set type = $1,
		unit = $2,
		symbol = $3,
		description = $4
		where pk_sensor_type_id = $5`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensorType.Type, sensorType.Unit, sensorType.Symbol, sensorType.Description, sensorType.Id); err != nil {
		return err
	}

	// make sure to get updated values
	sensorTypeNew, err := hdb.GetSensorTypeById(sensorType.Id)
	if err != nil {
		return err
	}
	sensorType.Type = sensorTypeNew.Type
	sensorType.Unit = sensorTypeNew.Unit
	sensorType.Symbol = sensorTypeNew.Symbol
	sensorType.Description = sensorTypeNew.Description

	return nil
}

func (hdb *homeInfoDB) NumOfSensorDataWithinLastHours(id, hours int) (int, error) {
	sql := `
		select count(*)
		from sensor_data
		where fk_sensor_id = $1
		and timestamp > now() - $2 * interval '1 hour'`
	if hdb.DatabaseType == "sqlite" {
		sql = `
		select count(*)
		from sensor_data
		where fk_sensor_id = $1
		and timestamp > datetime('now', '-' || $2 || ' hour'`
	}
	stmt, err := hdb.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id, hours)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (hdb *homeInfoDB) GetSensorData(id, limit int) ([]*SensorData, error) {
	sql := `
		select
			sd.timestamp, sd.value,
			s.pk_sensor_id, s.name, s.description,
			st.pk_sensor_type_id, st.type, st.unit, st.symbol, st.description
		from sensor_data sd
		join sensor s on s.pk_sensor_id = sd.fk_sensor_id
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.pk_sensor_id = $1
		order by sd.timestamp desc`
	if limit > 0 {
		sql += fmt.Sprintf(" limit %d", limit)
	}

	stmt, err := hdb.Prepare(sql)
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
		var s Sensor
		var st SensorType
		if err := rows.Scan(
			&d.Timestamp, &d.Value,
			&s.Id, &s.Name, &s.Description,
			&st.Id, &st.Type, &st.Unit, &st.Symbol, &st.Description,
		); err != nil {
			return nil, err
		}
		s.SensorType = st
		d.Sensor = s
		data = append(data, &d)
	}
	return data, nil
}

func (hdb *homeInfoDB) GetSensorValues(id, limit int) ([]*SensorValue, error) {
	sql := `
		select timestamp, value
		from sensor_data
		where fk_sensor_id = $1
		order by timestamp desc`
	if limit > 0 {
		sql += fmt.Sprintf(" limit %d", limit)
	}

	stmt, err := hdb.Prepare(sql)
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

func (hdb *homeInfoDB) GetHourlyAverages(id, limit int) ([]*SensorValue, error) {
	stmt, err := hdb.Prepare(`
    select hour, value
    from (select the_day + (the_hour * interval '1 hour') as hour, the_value as value
        from (select
            date(sd.timestamp) as the_day,
            extract(hour from sd.timestamp) as the_hour,
            round(avg(sd.value)) as the_value
            from sensor_data sd
            where sd.fk_sensor_id = $1
            group by 1,2
            order by 1 desc, 2 desc
        ) avg
        limit $2) sort
    order by 1 asc`)
	if hdb.DatabaseType == "sqlite" {
		stmt, err = hdb.Prepare(`
	    select hour, value
	    from (select datetime(the_day, '+' || the_hour || ' hour') as hour, the_value as value
	        from (select
	            date(sd.timestamp) as the_day,
	            strftime('%H', sd.timestamp) as the_hour,
	            round(avg(sd.value)) as the_value
	            from sensor_data sd
	            where sd.fk_sensor_id = $1
	            group by 1,2
	            order by 1 desc, 2 desc
	        ) avg
	        limit $2) sort
	    order by 1 asc`)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []*SensorValue{}
	for rows.Next() {
		var value SensorValue
		if hdb.DatabaseType == "sqlite" {
			var toConvert string
			if err := rows.Scan(&toConvert, &value.Value); err != nil {
				return nil, err
			}
			hour, _ := time.Parse("2006-01-02 15:04:05", toConvert)
			value.Timestamp = &hour
		} else {
			if err := rows.Scan(&value.Timestamp, &value.Value); err != nil {
				return nil, err
			}
		}
		values = append(values, &value)
	}
	return values, nil
}

func (hdb *homeInfoDB) GetDailyAverages(id, limit int) ([]*SensorValue, error) {
	stmt, err := hdb.Prepare(`
    select day, value
    from (select
        date(sd.timestamp) as day,
        round(avg(sd.value)) as value
        from sensor_data sd
        where sd.fk_sensor_id = $1
        group by 1
        order by 1 desc
        limit $2) sort
    order by 1 asc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []*SensorValue{}
	for rows.Next() {
		var value SensorValue
		if hdb.DatabaseType == "sqlite" {
			var toConvert string
			if err := rows.Scan(&toConvert, &value.Value); err != nil {
				return nil, err
			}
			hour, _ := time.Parse("2006-01-02 15:04:05", toConvert)
			value.Timestamp = &hour
		} else {
			if err := rows.Scan(&value.Timestamp, &value.Value); err != nil {
				return nil, err
			}
		}
		values = append(values, &value)
	}
	return values, nil
}

func (hdb *homeInfoDB) InsertSensorValue(sensorId, value int, timestamp time.Time) error {
	stmt, err := hdb.Prepare(`
		insert into sensor_data (fk_sensor_id, value, timestamp) values ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorId, value, timestamp)
	return err
}

func (hdb *homeInfoDB) DeleteAlert(alertId int) error {
	stmt, err := hdb.Prepare(`
		delete from alert
		where pk_alert_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(alertId)
	return err
}

func (hdb *homeInfoDB) DeleteSensor(sensorId int) error {
	stmt, err := hdb.Prepare(`
		delete from sensor
		where pk_sensor_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorId)
	return err
}

func (hdb *homeInfoDB) DeleteSensorType(sensorTypeId int) error {
	stmt, err := hdb.Prepare(`
		delete from sensor_type
		where pk_sensor_type_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorTypeId)
	return err
}

func (hdb *homeInfoDB) DeleteSensorValues(sensorId int) error {
	stmt, err := hdb.Prepare(`
		delete from sensor_data
		where fk_sensor_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sensorId)
	return err
}

func (hdb *homeInfoDB) GenerateSensorValues(id, num int) error {
	sensor, err := hdb.GetSensorById(id)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		value := rand.Intn(100)
		if strings.Contains(sensor.SensorType.Type, "state") {
			value = value % 2
		}
		timestamp := time.Unix(rand.Int63n(time.Now().Unix()-666666)+666666, 0).UTC()

		if err := hdb.InsertSensorValue(id, value, timestamp); err != nil {
			return err
		}
	}
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (hdb *homeInfoDB) Housekeeping() (err error) {
	deleteStmt, err := hdb.Prepare(`
		delete from sensor_data
		where timestamp < NOW() - INTERVAL '111 days'`)
	if err != nil {
		return err
	}
	defer deleteStmt.Close()

	log.Println("Housekeeping: delete rows older than [111] days ...")
	if _, err := deleteStmt.Exec(); err != nil {
		return err
	}
	return nil
}

func (hdb *homeInfoDB) GetQueues() ([]*Queue, error) {
	rows, err := hdb.Query(`
		select
			distinct mq.queue,
			count(pk_message_id),
			max(created_at)
		from message_queue mq
		group by mq.queue
		order by 1 asc, 2 desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qs := []*Queue{}
	for rows.Next() {
		var q Queue
		if err := rows.Scan(
			&q.Name, &q.MessageCount, &q.LastMessage,
		); err != nil {
			return nil, err
		}
		qs = append(qs, &q)
	}
	return qs, nil
}

func (hdb *homeInfoDB) GetQueueByName(name string) (*Queue, error) {
	stmt, err := hdb.Prepare(`
		select
			distinct mq.queue,
			count(pk_message_id),
			max(created_at)
		from message_queue mq
		where mq.queue = $1
		group by mq.queue`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var q Queue
	if err := stmt.QueryRow(name).Scan(
		&q.Name, &q.MessageCount, &q.LastMessage,
	); err != nil {
		return nil, err
	}
	return &q, nil
}

func (hdb *homeInfoDB) GetMessages() ([]*Message, error) {
	rows, err := hdb.Query(`
		select
			mq.pk_message_id,
			mq.queue,
			mq.message,
			mq.created_at
		from message_queue mq
		order by 4 desc, 2 asc, 1 desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := []*Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(
			&m.Id, &m.Queue, &m.Message, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		ms = append(ms, &m)
	}
	return ms, nil
}

func (hdb *homeInfoDB) GetMessagesFromQueue(name string) ([]*Message, error) {
	stmt, err := hdb.Prepare(`
		select
			mq.pk_message_id,
			mq.queue,
			mq.message,
			mq.created_at
		from message_queue mq
		where mq.queue = $1
		order by 4 desc, 2 asc, 1 desc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := []*Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(
			&m.Id, &m.Queue, &m.Message, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		ms = append(ms, &m)
	}
	return ms, nil
}

func (hdb *homeInfoDB) GetMessageById(id int) (*Message, error) {
	stmt, err := hdb.Prepare(`
		select
			mq.pk_message_id,
			mq.queue,
			mq.message,
			mq.created_at
		from message_queue mq
		where mq.pk_message_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var m Message
	if err := stmt.QueryRow(id).Scan(
		&m.Id, &m.Queue, &m.Message, &m.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (hdb *homeInfoDB) InsertMessage(message *Message) (err error) {
	stmt, err := hdb.Prepare(`
		insert into message_queue (queue, message) values ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(message.Queue, message.Message); err != nil {
		return err
	}
	return nil
}

func (hdb *homeInfoDB) DeleteQueue(name string) error {
	stmt, err := hdb.Prepare(`
		delete from message_queue
		where queue = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	return err
}

func (hdb *homeInfoDB) DeleteMessage(id int) error {
	stmt, err := hdb.Prepare(`
		delete from message_queue
		where pk_message_id = $1`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
