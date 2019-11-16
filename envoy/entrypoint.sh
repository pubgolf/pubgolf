#!/bin/bash

# Ignore warning about newline escapes in `envsubst` command.
# shellcheck disable=SC2140

# Required for env var support in the Envoy config file.
envsubst ''\
"\$ADMIN_PORT,"\
"\$GRPC_WEB_INBOUND_PORT,"\
"\$API_UPSTREAM_HOST,"\
"\$API_UPSTREAM_PORT" \
< /etc/envoy/config.yaml.tmpl > /etc/envoy/config.yaml

/usr/local/bin/envoy -c /etc/envoy/config.yaml
