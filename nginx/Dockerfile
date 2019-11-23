FROM nginx:latest

# Ensure that the `envsubst` command is installed.
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN rm /etc/nginx/conf.d/default.conf

COPY ./nginx/proxy.conf.tmpl /etc/nginx/conf.d/
COPY ./nginx/entrypoint.sh /etc/nginx/

CMD ["/etc/nginx/entrypoint.sh"]
