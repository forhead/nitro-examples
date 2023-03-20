FROM amazonlinux

RUN yum install python3 gcc net-tools -y

WORKDIR /app

COPY requirements.txt ./
COPY enclaveServer.py ./
COPY kmstool_enclave_cli ./
COPY libnsm.so /usr/lib64/
COPY kms.py ./
# you need set the ENV based on your environment
ENV REGION ap-southeast-1
RUN pip3 install -r requirements.txt

RUN chmod +x kmstool_enclave_cli

CMD ["python3", "/app/enclaveServer.py"]
