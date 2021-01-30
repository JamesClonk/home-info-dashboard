#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
source .env

# =============================================================================================
retry() {
    local -r -i max_attempts="$1"; shift
    local -r cmd="$@"
    local -i attempt_num=1

    until $cmd
    do
        if (( attempt_num == max_attempts ))
        then
            echo "Attempt $attempt_num failed and there are no more attempts left!"
            return 1
        else
            echo "Attempt $attempt_num failed! Trying again in $attempt_num seconds..."
            sleep $(( attempt_num++ ))
        fi
    done
}

# =============================================================================================
echo "waiting on postgres ..."
export PGPASSWORD=dev-secret
retry 10 psql -h 127.0.0.1 -U dev-user -d home_info_db -c '\q'
echo "postgres is up!"

# =============================================================================================
echo "setting up home-info db ..."
sleep 1

psql -h 127.0.0.1 -U dev-user -d home_info_db  < lib/database/migrations/postgres/001_sensor_tables.up.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < lib/database/migrations/postgres/002_setup_sensors.up.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < lib/database/migrations/postgres/003_alert_tables.up.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < lib/database/migrations/postgres/004_setup_alerts.up.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < lib/database/migrations/postgres/005_message_queue.up.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < lib/database/migrations/postgres/006_alert_active_flag.up.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < _fixtures/migration.sql
#psql -h 127.0.0.1 -U dev-user -d home_info_db  < _fixtures/values.sql
psql -h 127.0.0.1 -U dev-user -d home_info_db  < _fixtures/export.sql
