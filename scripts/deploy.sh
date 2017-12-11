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
echo "deploying home-info-dashboard ..."

wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.32.0&source=github-rel' -qO cf-cli.tgz
tar -xvzf cf-cli.tgz 1>/dev/null
chmod +x cf
rm -f cf-cli.tgz || true

./cf login -a "https://api.lyra-836.appcloud.swisscom.com" -u "${APC_USERNAME}" -p "${APC_PASSWORD}" -o "${APC_ORGANIZATION}" -s "${APC_SPACE}"

# make sure routes will be ready
./cf create-route "${APC_SPACE}" scapp.io --hostname weather
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname weather
./cf create-route "${APC_SPACE}" scapp.io --hostname home-info
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname home-info
./cf create-route "${APC_SPACE}" scapp.io --hostname home-info-dashboard
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname home-info-dashboard
./cf create-route "${APC_SPACE}" scapp.io --hostname home-info-blue-green
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname home-info-blue-green
sleep 2

# secure working app
./cf rename home-info-dashboard home-info-dashboard-old || true
./cf unmap-route home-info-dashboard-old scapp.io --hostname home-info-blue-green || true
sleep 2

# push new app
./cf push home-info-dashboard-new --no-route
./cf map-route home-info-dashboard-new scapp.io --hostname home-info-blue-green
./cf map-route home-info-dashboard-new applicationcloud.io --hostname home-info-blue-green
sleep 5

# test app
response=$(curl -sIL -w "%{http_code}" -o /dev/null "home-info-blue-green.scapp.io")
if [[ "${response}" != "200" ]]; then
    ./cf delete -f home-info-dashboard-new || true
    echo "App did not respond as expected, HTTP [${response}]"
    exit 1
fi

# finish blue-green deployment of app
./cf delete -f home-info-dashboard || true
./cf rename home-info-dashboard-new home-info-dashboard
./cf map-route home-info-dashboard scapp.io --hostname weather
./cf map-route home-info-dashboard applicationcloud.io --hostname weather
./cf map-route home-info-dashboard scapp.io --hostname home-info
./cf map-route home-info-dashboard applicationcloud.io --hostname home-info
./cf map-route home-info-dashboard scapp.io --hostname home-info-dashboard
./cf map-route home-info-dashboard applicationcloud.io --hostname home-info-dashboard
./cf unmap-route home-info-dashboard scapp.io --hostname home-info-blue-green || true
./cf unmap-route home-info-dashboard applicationcloud.io --hostname home-info-blue-green || true
./cf delete -f home-info-dashboard-old

# show status
./cf apps
./cf app home-info-dashboard

./cf logout

rm -f cf || true
rm -f LICENSE || true
rm -f NOTICE || true
