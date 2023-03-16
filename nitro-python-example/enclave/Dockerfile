FROM amazonlinux

RUN yum install python3 gcc net-tools -y

WORKDIR /app

COPY requirements.txt ./

COPY traffic-forwarder.py ./
COPY enclaveServer.py ./
COPY run.sh ./

RUN pip3 install -r requirements.txt

RUN chmod +x run.sh
CMD ["/app/run.sh"]
