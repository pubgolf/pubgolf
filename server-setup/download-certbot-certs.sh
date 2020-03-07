#!/bin/bash

sudo docker run -it --rm \
-v "${HOME}/data/certbot/conf:/etc/letsencrypt" \
-v "${HOME}/data/certbot/www:/var/www/certbot" \
-v "${HOME}/webroot:/data/letsencrypt" \
certbot/certbot \
certonly --webroot \
--webroot-path=/data/letsencrypt \
--register-unsafely-without-email --agree-tos \
-d staging.pubgolf.co -d www-staging.pubgolf.co -d app-staging.pubgolf.co -d api-staging.pubgolf.co
