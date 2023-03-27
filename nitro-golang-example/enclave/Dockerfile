FROM amazonlinux

WORKDIR /app

COPY kmstool_enclave_cli ./
COPY libnsm.so /usr/lib64/
COPY enclaveServer ./

ENV REGION ap-northeast-1

RUN chmod +x kmstool_enclave_cli
RUN chmod +x enclaveServer

CMD ["./app/enclaveServer"]
