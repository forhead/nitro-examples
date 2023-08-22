####################################################################################################
## Builder
####################################################################################################
FROM rust:latest AS builder

ARG TARGETPLATFORM
# x86_64-unknown-linux-gnu
# aarch64-unknown-linux-gnu

RUN rustup target add $TARGETPLATFORM
RUN apt update && apt install -y musl-tools musl-dev
RUN update-ca-certificates

WORKDIR /app

COPY ./ .

RUN cargo test --package enclave
RUN cargo build --target $TARGETPLATFORM --release

####################################################################################################
## Final image
####################################################################################################
FROM ubuntu
ARG TARGETPLATFORM

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

# Copy our build
COPY --from=builder /app/target/${TARGETPLATFORM}/release/enclave-server /app/enclave-server
# Copy kmstool_enclave_cli
COPY --from=builder /app/kmstool_enclave_cli /app/kmstool_enclave_cli
COPY --from=builder /app/libnsm.so /usr/lib64/libnsm.so

ENV REGION ap-east-1
ENV LD_LIBRARY_PATH="/usr/lib64":"/lib"

# RUN ls -lha /app
# RUN ls -lha /lib
RUN chmod +x /app/kmstool_enclave_cli
RUN chmod +x /app/enclave-server

CMD ["./app/enclave-server", "-p", "7878"]
