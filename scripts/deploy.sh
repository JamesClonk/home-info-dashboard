#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [ -z "${APC_USERNAME}" ]; then
	echo "APC_USERNAME must be set!"
	exit 1
fi
if [ -z "${APC_PASSWORD}" ]; then
	echo "APC_PASSWORD must be set!"
	exit 1
fi
if [ -z "${APC_ORGANIZATION}" ]; then
	echo "APC_ORGANIZATION must be set!"
	exit 1
fi
if [ -z "${APC_SPACE}" ]; then
	echo "APC_SPACE must be set!"
	exit 1
fi

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
echo "deploying weather_app ..."

wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.32.0&source=github-rel' -qO cf-cli.tgz
tar -xvzf cf-cli.tgz 1>/dev/null
chmod +x cf
rm -f cf-cli.tgz || true

./cf login -a "https://api.lyra-836.appcloud.swisscom.com" -u "${APC_USERNAME}" -p "${APC_PASSWORD}" -o "${APC_ORGANIZATION}" -s "${APC_SPACE}"

# make sure routes will be ready
./cf create-route "${APC_SPACE}" scapp.io --hostname weather-app
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname weather-app
./cf create-route "${APC_SPACE}" scapp.io --hostname weatherapp
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname weatherapp
./cf create-route "${APC_SPACE}" scapp.io --hostname smarttemperature
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname smarttemperature
./cf create-route "${APC_SPACE}" scapp.io --hostname weather-app-blue-green
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname weather-app-blue-green
sleep 2

# secure working app
./cf rename weather_app weather_app_old || true
./cf unmap-route weather_app_old scapp.io --hostname weather-app-blue-green || true
sleep 2

# push new app
./cf push weather_app_new --no-route
./cf map-route weather_app_new scapp.io --hostname weather-app-blue-green
./cf map-route weather_app_new applicationcloud.io --hostname weather-app-blue-green
sleep 5

# test app
response=$(curl -sIL -w "%{http_code}" -o /dev/null "weather-app-blue-green.scapp.io")
if [[ "${response}" != "200" ]]; then
    ./cf delete -f weather_app_new || true
    echo "App did not respond as expected, HTTP [${response}]"
    exit 1
fi

# finish blue-green deployment of app
./cf delete -f weather_app || true
./cf rename weather_app_new weather_app
./cf map-route weather_app scapp.io --hostname weather-app
./cf map-route weather_app applicationcloud.io --hostname weather-app
./cf map-route weather_app scapp.io --hostname weatherapp
./cf map-route weather_app applicationcloud.io --hostname weatherapp
./cf map-route weather_app scapp.io --hostname smarttemperature
./cf map-route weather_app applicationcloud.io --hostname smarttemperature
./cf unmap-route weather_app scapp.io --hostname weather-app-blue-green || true
./cf unmap-route weather_app applicationcloud.io --hostname weather-app-blue-green || true
./cf delete -f weather_app_old

./cf logout

rm -f cf || true
rm -f LICENSE || true
rm -f NOTICE || true
