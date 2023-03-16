FROM amazonlinux

RUN yum install python3 gcc net-tools golang -y

WORKDIR /app

COPY requirements.txt ./
COPY kmstool_enclave_cli ./
COPY libnsm.so ./

COPY traffic-forwarder.py ./
COPY enclaveServer ./
COPY run.sh ./

RUN pip3 install -r requirements.txt

RUN chmod +x run.sh
CMD ["/app/run.sh"]
