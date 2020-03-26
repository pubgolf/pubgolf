#!/bin/bash
# Generate client certificates for a docker host.
# Command line params:
# - $1 - Path to ca.pem
# - $2 - Path to ca-key.pem
# - $3 - Env name to label key files with.

if [ "$1" = "" ] || [ "$2" = "" ] || [ "$3" = "" ]; then
    echo "Please provide paths to the ca.pem and ca-key.pem files, and the env name."
    exit 2
fi

echo 01 > ca.srl
cp "$1" ca.pem
cp "$2" ca-key.pem

openssl genrsa -out "../keys/docker-${3}client-key.pem"
openssl req -new -key "../keys/docker-${3}client-key.pem" -out client.csr

echo 'extendedKeyUsage = clientAuth' > extfile.cnf

openssl x509 -req -days 365 \
    -in client.csr \
    -out "../keys/docker-${3}client-cert.pem" \
    -CA ca.pem \
    -CAkey ca-key.pem \
    -extfile extfile.cnf

# Cleanup
rm ca.srl ca.pem ca-key.pem client.csr extfile.cnf
