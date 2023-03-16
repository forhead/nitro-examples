FROM amazonlinux

RUN yum install python3 gcc net-tools golang -y

ENV GOPATH=/app/go

WORKDIR /app
# COPY libnsm.so ./
COPY requirements.txt ./

RUN pip3 install -r requirements.txt

COPY traffic-forwarder.py ./
COPY run.sh ./
COPY vsockEnclave.go ./
COPY go.mod ./
COPY go.sum ./
COPY vsockEnclave ./

RUN chmod +x run.sh
CMD ["/app/run.sh"]