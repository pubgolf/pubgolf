#!/bin/bash

# Ignore warning about newline escapes in `envsubst` command.
# shellcheck disable=SC2140

# Required for env var support in the Envoy config file.
envsubst ''\
"\$ENVOY_ADMIN_PORT,"\
"\$GRPC_WEB_PORT,"\
"\$API_HOST,"\
"\$API_PORT" \
< /etc/envoy/config.yaml.tmpl > /etc/envoy/config.yaml

/usr/local/bin/envoy -c /etc/envoy/config.yaml
