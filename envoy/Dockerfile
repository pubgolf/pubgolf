FROM envoyproxy/envoy:latest

# Ensure that the `envsubst` command is installed.
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY ./envoy/config.yaml.tmpl /etc/envoy/config.yaml.tmpl
COPY ./envoy/entrypoint.sh /etc/envoy/entrypoint.sh
CMD ["/etc/envoy/entrypoint.sh"]
