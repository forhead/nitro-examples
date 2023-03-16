#!/bin/bash
FILE=enclaveServer.eif
if [ -f "$FILE" ]; then
    rm $FILE
fi
docker rmi enclaveServer:latest
pkill vsock-proxy

go build enclaveServer.go

docker build . -t enclaveServer
nitro-cli build-enclave --docker-uri enclaveServer:latest --output-file enclaveServer.eif
vsock-proxy 8000 kms.ap-southeast-1.amazonaws.com 443 &

nitro-cli run-enclave --cpu-count 4 --memory 5240 --enclave-cid 16 --eif-path enclaveServer.eif --debug-mode --attach-console