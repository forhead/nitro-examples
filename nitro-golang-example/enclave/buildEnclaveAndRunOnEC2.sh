#!/bin/bash

FILE=enclave_server_go.eif
ENCLAVE_CPU_COUNT=4
ENCLAVE_MEMORY_SIZE=768
ENCLAVE_CLIENT_CID=16

if [ -f "$FILE" ]; then
    rm $FILE
fi

RunningEnclave=$(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
if [ -n "$RunningEnclave" ]; then
	nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID");
fi

docker rmi enclave_server_go:latest
pkill vsock-proxy

go build enclaveServer.go

docker build -t enclave_server_go:latest .
nitro-cli build-enclave --docker-uri enclave_server_go:latest  --output-file $FILE > enclave_server_go.log
# change this to your kms's regional endpoint
vsock-proxy 8000 kms.ap-northeast-1.amazonaws.com 443 &

nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE --enclave-cid $ENCLAVE_CLIENT_CID --eif-path $FILE --debug-mode --attach-console
