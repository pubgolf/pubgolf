#!/bin/bash

sudo docker run \
  -v "${HOME}/nginx-conf/nginx.conf:/etc/nginx/conf.d/default.conf" \
  -v "${HOME}/webroot:/usr/share/nginx/html" \
  -p "80:80" \
  nginx:latest 
