#!/bin/bash

# Ignore warning about newline escapes in `envsubst` command.
# shellcheck disable=SC2140

# Required for env var support in the Nginx config file.
envsubst ''\
"\$GRPC_WEB_PORT,"\
"\$WEB_APP_PORT"\
< /etc/nginx/conf.d/proxy.conf.tmpl > /etc/nginx/conf.d/proxy.conf

nginx -g "daemon off;"
