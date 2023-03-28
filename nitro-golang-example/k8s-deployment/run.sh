#!/bin/bash -e
# Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

readonly EIF_PATH="/home/enclaveserver.eif"
readonly ENCLAVE_CPU_COUNT=4
readonly ENCLAVE_MEMORY_SIZE=2048

main() {

    cd /home/
    ls -ltr
    vsock-proxy 8000 kms.ap-northeast-1.amazonaws.com 443 &

    nitro-cli run-enclave --cpu-count $ENCLAVE_CPU_COUNT --memory $ENCLAVE_MEMORY_SIZE \
        --eif-path $EIF_PATH --debug-mode 

    local enclave_id=$(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
    echo "-------------------------------"
    echo "Enclave ID is $enclave_id"
    echo "-------------------------------"

    nitro-cli console --enclave-id $enclave_id # blocking call.

}

main