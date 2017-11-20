#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
source .env_mysql

# =============================================================================================
echo "setting up weatherdb ..."
sleep 5

mysql --host=127.0.0.1 --user=blubb --password=blabb --database=weather_db -v < lib/database/migrations/mysql/001_weatherdb_sensor_tables.up.sql
mysql --host=127.0.0.1 --user=blubb --password=blabb --database=weather_db -v < lib/database/migrations/mysql/002_weatherdb_setup_sensors.up.sql
mysql --host=127.0.0.1 --user=blubb --password=blabb --database=weather_db -v < lib/database/migrations/mysql/003_weatherdb_config_table.up.sql
mysql --host=127.0.0.1 --user=blubb --password=blabb --database=weather_db -v < _fixtures/migration.sql
mysql --host=127.0.0.1 --user=blubb --password=blabb --database=weather_db -v < _fixtures/sensors.sql
mysql --host=127.0.0.1 --user=blubb --password=blabb --database=weather_db -v < _fixtures/values.sql
