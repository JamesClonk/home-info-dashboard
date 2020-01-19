package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type HomeInfoDB interface {
	GetTemperature() (int64, error)
	GetSensors() ([]*Sensor, error)
	GetSensorById(int) (*Sensor, error)
	GetSensorsByName(string) ([]*Sensor, error)
	GetSensorsByTypeId(int) ([]*Sensor, error)
	GetSensorTypeById(int) (*SensorType, error)
	GetSensorTypeByType(string) (*SensorType, error)
	GetSensorTypes() ([]*SensorType, error)
	GetSensorData(int, int) ([]*SensorData, error)
	GetSensorValues(int, int) ([]*SensorValue, error)
	GetHourlyAverages(int, int) ([]*SensorValue, error)
	GetDailyAverages(int, int) ([]*SensorValue, error)
	InsertSensor(*Sensor) error
	InsertSensorType(*SensorType) error
	InsertSensorValue(int, int, time.Time) error
	UpdateSensor(*Sensor) error
	UpdateSensorType(*SensorType) error
	DeleteSensor(int) error
	DeleteSensorType(int) error
	DeleteSensorValues(int) error
	GenerateSensorValues(int, int) error
	Housekeeping(int, int) error
}

type homeInfoDB struct {
	*sql.DB
	DatabaseType string
}

func NewHomeInfoDB(adapter Adapter) HomeInfoDB {
	return &homeInfoDB{adapter.GetDatabase(), adapter.GetType()}
}

func (hdb *homeInfoDB) GetTemperature() (int64, error) {
	rows, err := hdb.Query(`
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

func (hdb *homeInfoDB) GetSensors() ([]*Sensor, error) {
	rows, err := hdb.Query(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
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
		if err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.TypeId, &s.Unit, &s.Description); err != nil {
			return nil, err
		}
		ss = append(ss, &s)
	}
	return ss, nil
}

func (hdb *homeInfoDB) GetSensorsByTypeId(id int) ([]*Sensor, error) {
	stmt, err := hdb.Prepare(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
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
		if err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.TypeId, &s.Unit, &s.Description); err != nil {
			return nil, err
		}
		ss = append(ss, &s)
	}
	return ss, nil
}

func (hdb *homeInfoDB) GetSensorById(id int) (*Sensor, error) {
	stmt, err := hdb.Prepare(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
		from sensor s
		join sensor_type st on s.fk_sensor_type_id = st.pk_sensor_type_id
		where s.pk_sensor_id = $1`)
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

func (hdb *homeInfoDB) GetSensorsByName(name string) ([]*Sensor, error) {
	stmt, err := hdb.Prepare(`
		select s.pk_sensor_id, s.name, st.type, st.pk_sensor_type_id, st.unit, s.description
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
		if err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.TypeId, &s.Unit, &s.Description); err != nil {
			return nil, err
		}
		ss = append(ss, &s)
	}
	return ss, nil
}

func (hdb *homeInfoDB) InsertSensor(sensor *Sensor) (err error) {
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
		sensorType, err = hdb.GetSensorTypeById(int(sensorTypeId))
		if err != nil {
			return err
		}
	} else { // by type
		sensorType, err = hdb.GetSensorTypeByType(sensor.Type)
		if err != nil {
			return err
		}
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
	sensor.Type = sensorNew.Type
	sensor.TypeId = sensorNew.TypeId
	sensor.Unit = sensorNew.Unit
	sensor.Description = sensorNew.Description

	return nil
}

func (hdb *homeInfoDB) UpdateSensor(sensor *Sensor) (err error) {
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
		sensorType, err = hdb.GetSensorTypeById(int(sensorTypeId))
		if err != nil {
			return err
		}
	} else { // by type
		sensorType, err = hdb.GetSensorTypeByType(sensor.Type)
		if err != nil {
			return err
		}
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
	sensor.Type = sensorNew.Type
	sensor.TypeId = sensorNew.TypeId
	sensor.Unit = sensorNew.Unit
	sensor.Description = sensorNew.Description

	return nil
}

func (hdb *homeInfoDB) GetSensorTypeById(id int) (*SensorType, error) {
	stmt, err := hdb.Prepare(`
		select pk_sensor_type_id, type, unit, description
		from sensor_type
		where pk_sensor_type_id = $1`)
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

func (hdb *homeInfoDB) GetSensorTypeByType(s string) (*SensorType, error) {
	stmt, err := hdb.Prepare(`
		select pk_sensor_type_id, type, unit, description
		from sensor_type
		where type = $1`)
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

func (hdb *homeInfoDB) GetSensorTypes() ([]*SensorType, error) {
	rows, err := hdb.Query(`
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

func (hdb *homeInfoDB) InsertSensorType(sensorType *SensorType) (err error) {
	stmt, err := hdb.Prepare(`
		insert into sensor_type (type, unit, description) values ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensorType.Type, sensorType.Unit, sensorType.Description); err != nil {
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
	sensorType.Description = sensorTypeNew.Description

	return nil
}

func (hdb *homeInfoDB) UpdateSensorType(sensorType *SensorType) (err error) {
	stmt, err := hdb.Prepare(`
		update sensor_type
		set type = $1,
		unit = $2,
		description = $3
		where pk_sensor_type_id = $4`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(sensorType.Type, sensorType.Unit, sensorType.Description, sensorType.Id); err != nil {
		return err
	}

	// make sure to get updated values
	sensorTypeNew, err := hdb.GetSensorTypeById(sensorType.Id)
	if err != nil {
		return err
	}
	sensorType.Type = sensorTypeNew.Type
	sensorType.Unit = sensorTypeNew.Unit
	sensorType.Description = sensorTypeNew.Description

	return nil
}

func (hdb *homeInfoDB) GetSensorData(id, limit int) ([]*SensorData, error) {
	sql := `
		select s.pk_sensor_id, sd.timestamp, s.name, st.type, st.unit, sd.value
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
		if err := rows.Scan(&d.SensorId, &d.Timestamp, &d.Name, &d.Type, &d.Unit, &d.Value); err != nil {
			return nil, err
		}
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
		if strings.Contains(sensor.Type, "state") {
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

func (hdb *homeInfoDB) Housekeeping(days, rows int) (err error) {
	// housekeeping logic: select count(*) from sensor_data where timestamp < now - $days
	// if count > $rows, then: delete from sensor_data where timestamp < now - $days
	cutOff := time.Now().AddDate(0, -days, 0).UTC()

	getStmt, err := hdb.Prepare(`
		select count(*)
		from sensor_data
		where fk_sensor_id = $1
		and timestamp > $2`)
	if err != nil {
		return err
	}
	defer getStmt.Close()

	deleteStmt, err := hdb.Prepare(`
		delete from sensor_data
		where fk_sensor_id = $1
		and timestamp < $2`)
	if err != nil {
		return err
	}
	defer deleteStmt.Close()

	sensors, err := hdb.GetSensors()
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		row, err := getStmt.Query(sensor.Id, cutOff)
		if err != nil {
			return err
		}
		defer row.Close()

		if row.Next() {
			var count int64
			if err := row.Scan(&count); err != nil {
				return err
			}
			log.Printf("Housekeeping: Sensor [%d:%s] has [%v] minimum rows ...\n", sensor.Id, sensor.Name, count)

			if count > int64(rows) {
				log.Printf("Housekeeping: Will now delete excess rows of sensor [%d:%s] ...\n", sensor.Id, sensor.Name)
				_, err := deleteStmt.Exec(sensor.Id, cutOff)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
