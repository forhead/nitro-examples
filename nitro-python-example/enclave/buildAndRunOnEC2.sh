#!/bin/bash
FILE=walletenclave.eif
if [ -f "$FILE" ]; then
    rm $FILE
fi

RunningEnclave=$(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
if [ -n "$RunningEnclave" ]; then
	nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID");
fi

#docker rmi -f $(docker images -a -q)
docker rmi enclaveserver:latest
pkill vsock-proxy

docker build -t enclaveserver:latest .
nitro-cli build-enclave --docker-uri enclaveserver:latest  --output-file enclaveserver.eif > EnclaveImage.log

vsock-proxy 8000 kms.ap-southeast-1.amazonaws.com 443 &

nitro-cli run-enclave --cpu-count 4 --memory 5240 --enclave-cid 16 --eif-path enclaveserver.eif --debug-mode --attach-console
