#!/usr/bin/env sh
set -ue

: "${BASIC_USERNAME:=mockbucket}"
: "${BASIC_PASSWORD:=mockbucket123!}"
: "${VERBOSE:=false}"

if [ "${VERBOSE}" = "true" ]; then
  set -x
fi

echo "Creating Credentials"
htpasswd -bc /www/htpasswd $BASIC_USERNAME $BASIC_PASSWORD

spawn-fcgi -s /var/run/fcgiwrap.socket /usr/bin/fcgiwrap
/apiserver/apiserver --port 8081 &

exec "$@"
