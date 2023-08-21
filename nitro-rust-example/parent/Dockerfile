####################################################################################################
## Builder
####################################################################################################
FROM rust:latest AS builder

ARG TARGETPLATFORM
# x86_64-unknown-linux-musl
# aarch64-unknown-linux-musl

RUN rustup target add $TARGETPLATFORM
RUN apt update && apt install -y musl-tools musl-dev
RUN update-ca-certificates

# Create appuser
ENV USER=myip
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"


WORKDIR /myip

COPY ./ .

RUN cargo build --target $TARGETPLATFORM --release

####################################################################################################
## Final image
####################################################################################################
FROM alpine
ARG TARGETPLATFORM

# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /myip

# Copy our build
COPY --from=builder /myip/target/${TARGETPLATFORM}/release/enclave-client ./

# Use an unprivileged user.
USER myip:myip

CMD ["./myip/enclave-client", "-c", "16", "-p", "7878", "-k", "205b2e71-5e9a-4446-9b9e-962c5cb042bd", "-t", "AccountTable"]
