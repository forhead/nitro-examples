#!/bin/bash
FILE=vsockdemo.eif
if [ -f "$FILE" ]; then
    rm $FILE
fi
#docker rmi -f $(docker images -a -q)
docker rmi kms-demo:latest
pkill vsock-proxy

docker build . -t vsockdemo
nitro-cli build-enclave --docker-uri vsockdemo:latest --output-file vsockdemo.eif
vsock-proxy 8000 kms.ap-southeast-1.amazonaws.com 443 &

nitro-cli run-enclave --cpu-count 4 --memory 5240 --enclave-cid 16 --eif-path vsockdemo.eif --debug-mode --attach-console
