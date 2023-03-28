# Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

######## full image ########

FROM amazonlinux:2 as full_image

# Install packages
#
RUN amazon-linux-extras install aws-nitro-enclaves-cli && \
    yum install aws-nitro-enclaves-cli-devel jq -y

WORKDIR /ne-deps

# Copy only the required binaries to /ne-deps folder.
#
RUN BINS="\
    /usr/bin/nitro-cli \
    /usr/bin/lscpu \
    /usr/bin/vsock-proxy \
    /usr/bin/jq \
    " && \
    for bin in $BINS; do \
        { echo "$bin"; ldd "$bin" | grep -Eo "/.*lib.*/[^ ]+"; } | \
            while read path; do \
                mkdir -p ".$(dirname $path)"; \
                cp -fL "$path" ".$path"; \
            done \
    done

# Prepare other required files and folders for the final image.
#
RUN mkdir -p /ne-deps/etc/nitro_enclaves && \
    mkdir -p /ne-deps/run/nitro_enclaves && \
    mkdir -p /ne-deps/var/log/nitro_enclaves && \
    cp -rf /usr/local/ /ne-deps/usr/ && \
    cp -rf /usr/share/nitro_enclaves/ /ne-deps/usr/share/ && \
    cp -f /etc/nitro_enclaves/vsock-proxy.yaml /ne-deps/etc/nitro_enclaves

######## kms image ########

FROM amazonlinux:2 as image

COPY --from=full_image /ne-deps /

# Set Environment Variables

ENV REGION = ap-northeast-1

COPY enclaveserver.eif /home
COPY run.sh  /home

RUN chmod +x /home/run.sh
# RUN systemctl restart nitro-enclaves-allocator.service

CMD ["/home/run.sh"]