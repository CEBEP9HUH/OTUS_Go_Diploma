FROM ubuntu:18.04 AS buildbase

RUN apt update -y && apt upgrade -y && \
    apt install -y wget && \
    apt install -y make && \
    apt install -y unzip && \
    apt-get install -y build-essential


FROM buildbase AS gobuildbase
RUN wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV CGO_ENABLED=1
RUN go env -w GOPATH=/usr/local/go

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v28.2/protoc-28.2-linux-x86_64.zip && \
    unzip -a protoc-28.2-linux-x86_64.zip && \
    cp -r ./include/* /usr/local/include/ && \
    cp ./bin/protoc /usr/local/bin

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    export PATH="$PATH:$(go env GOPATH)/bin"


FROM gobuildbase AS builder
RUN mkdir /opt/diploma
WORKDIR /opt/diploma

CMD [ "make", "all" ]