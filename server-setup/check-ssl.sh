#!/bin/bash

if [ "$1" = "" ]; then
    echo "Please provide hostname to check cert for."
    exit 2
fi

echo | \
openssl s_client -showcerts -servername "$1" -connect "$1:443" 2>/dev/null | \
openssl x509 -inform pem -noout -text
