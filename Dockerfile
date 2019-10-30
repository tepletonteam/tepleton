# Simple usage with a mounted data directory:
# > docker build -t ton .
# > docker run -v $HOME/.tond:/root/.tond ton init
# > docker run -v $HOME/.tond:/root/.tond ton start

FROM alpine:edge

# Install minimum necessary dependencies

ENV PACKAGES go glide make git libc-dev bash
RUN apk add --no-cache $PACKAGES

# Set up GOPATH & PATH

ENV GOPATH       /root/go
ENV BASE_PATH    $GOPATH/src/github.com/tepleton
ENV REPO_PATH    $BASE_PATH/tepleton-sdk
ENV WORKDIR      /tepleton/
ENV PATH         $GOPATH/bin:$PATH

# Link expected Go repo path

RUN mkdir -p $WORKDIR $GOPATH/pkg $ $GOPATH/bin $BASE_PATH && ln -sf $WORKDIR $REPO_PATH

# Add source files

ADD . $WORKDIR

# Build tepleton-sdk

RUN cd $REPO_PATH && make get_tools && make get_vendor_deps && make all && make install

# Remove packages

RUN apk del $PACKAGES

# Set entrypoint

ENTRYPOINT ["/root/go/bin/tond"]
