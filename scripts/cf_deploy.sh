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
if [ -z "${AUTH_USERNAME}" ]; then
	echo "AUTH_USERNAME must be set!"
	exit 1
fi
if [ -z "${AUTH_PASSWORD}" ]; then
	echo "AUTH_PASSWORD must be set!"
	exit 1
fi
if [ -z "${TELEBOT_TOKEN}" ]; then
	echo "TELEBOT_TOKEN must be set!"
	exit 1
fi
if [ -z "${TELEBOT_USERNAME}" ]; then
	echo "TELEBOT_USERNAME must be set!"
	exit 1
fi
if [ -z "${SLACK_TOKEN}" ]; then
	echo "SLACK_TOKEN must be set!"
	exit 1
fi
if [ -z "${SLACK_CHANNEL}" ]; then
	echo "SLACK_CHANNEL must be set!"
	exit 1
fi

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
echo "deploying home-info-dashboard ..."

wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.43.0&source=github-rel' -qO cf-cli.tgz
tar -xvzf cf-cli.tgz 1>/dev/null
chmod +x cf
rm -f cf-cli.tgz || true
rm -f LICENSE || true
rm -f NOTICE || true

./cf login -a "https://api.lyra-836.appcloud.swisscom.com" -u "${APC_USERNAME}" -p "${APC_PASSWORD}" -o "${APC_ORGANIZATION}" -s "${APC_SPACE}"

# push app
./cf push home-info -f manifest.yml \
  --var auth_username=${AUTH_USERNAME} --var auth_password=${AUTH_PASSWORD} --var telebot_token=${TELEBOT_TOKEN} --var telebot_username=${TELEBOT_USERNAME} --var slack_token=${SLACK_TOKEN} --var slack_channel="${SLACK_CHANNEL}"
sleep 5

# show status
./cf app home-info

./cf logout

rm -f cf || true
