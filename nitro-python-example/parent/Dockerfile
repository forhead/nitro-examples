FROM amazonlinux:2

RUN yum install python3 -y

WORKDIR /app

COPY requirements.txt ./
COPY appClient.py ./
# you need set the ENV based on your environment
ENV REGION ap-southeast-1
RUN pip3 install -r requirements.txt

CMD ["python3", "/app/appClient.py"]
