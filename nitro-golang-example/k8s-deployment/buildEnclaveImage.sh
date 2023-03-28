#!/bin/bash

docker build -t enclaveserver:latest ../enclave/

nitro-cli build-enclave --docker-uri enclaveserver:latest  --output-file enclaveserver.eif > EnclaveImage.log