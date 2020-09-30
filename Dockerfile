# Simple usage with a mounted data directory:
# > docker build -t ton .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.tond:/root/.tond -v ~/.toncli:/root/.toncli ton tond init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.tond:/root/.tond -v ~/.toncli:/root/.toncli ton tond start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/tepleton/tepleton-sdk

# Add source files
COPY . .

# Install minimum necessary dependencies, build Tepleton SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make get_tools && \
    make get_vendor_deps && \
    make build && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/tond /usr/bin/tond
COPY --from=build-env /go/bin/toncli /usr/bin/toncli

# Run tond by default, omit entrypoint to ease using container with toncli
CMD ["tond"]
