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
docker rmi walletenclave:latest
pkill vsock-proxy

docker build -t walletenclave:latest .
nitro-cli build-enclave --docker-uri walletenclave:latest  --output-file walletenclave.eif > EnclaveImage.log

vsock-proxy 8000 kms.ap-southeast-1.amazonaws.com 443 &

nitro-cli run-enclave --cpu-count 4 --memory 5240 --enclave-cid 16 --eif-path walletenclave.eif --debug-mode --attach-console
