#!/bin/bash

# Ignore warning about newline escapes in `envsubst` command.
# shellcheck disable=SC2140

# Required for env var support in the Nginx config file.
envsubst ''\
"\$ROOT_DOMAIN"\
"\$WWW_SUBDOMAIN"\
"\$WEB_APP_SUBDOMAIN"\
"\$WEB_APP_PORT"\
"\$API_SUBDOMAIN"\
"\$GRPC_WEB_PORT,"\
"\$EVENT_KEY_REDIRECT,"\
< /etc/nginx/conf.d/proxy.conf.tmpl > /etc/nginx/conf.d/proxy.conf

while :; do
  sleep 6h & wait ${!};
  nginx -s reload;
done &

nginx -g "daemon off;"
