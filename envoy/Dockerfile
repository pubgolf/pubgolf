FROM envoyproxy/envoy:latest

# Ensure that the `envsubst` command is installed.
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
CMD ["/etc/envoy/entrypoint.sh"]
