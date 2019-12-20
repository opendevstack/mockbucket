#!/usr/bin/env bash
set -ue

: "${BASIC_USERNAME:=mockbucket}"
: "${BASIC_PASSWORD:=mockbucket123!}"
: "${VERBOSE:=false}"
: "${REPOS:=}"

if [ "${VERBOSE}" = "true" ]; then
  set -x
fi

echo "Creating Credentials"
htpasswd -bc /www/htpasswd $BASIC_USERNAME $BASIC_PASSWORD

IFS=';' read -ra REPOS <<< "${REPOS}"

for REPO in "${REPOS[@]}"
do
    mkdir -p "/scm/${REPO}"
    git -C "/scm/${REPO}" init --bare
done


spawn-fcgi -s /var/run/fcgiwrap.socket /usr/bin/fcgiwrap
/apiserver/apiserver --port 8081 &

exec "$@"
