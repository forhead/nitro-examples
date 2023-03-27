FROM amazonlinux

WORKDIR /app

COPY kmstool_enclave_cli ./
COPY libnsm.so ./
COPY enclaveServer ./

RUN chmod +x kmstool_enclave_cli
RUN chmod +x enclaveServer

CMD ["./app/enclaveServer"]
