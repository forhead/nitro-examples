
FROM amazonlinux

# Set Environment Variables
ENV REGION = ap-northeast-1

COPY appClient /home

RUN chmod +x /home/appClient

CMD ["./home/appClient"]