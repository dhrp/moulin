FROM golang:1.24 AS builder

ENV PROTOC_VERSION="3.5.1"

RUN apt-get update && \
	apt-get install unzip &&\
	rm -rf /var/lib/apt/lists/

# install protoc
RUN curl -L https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > protoc.zip && \
	unzip protoc.zip -d .

# Download modules
COPY go.mod go.sum /go/src/github.com/dhrp/moulin/
WORKDIR /go/src/github.com/dhrp/moulin/
RUN go mod download

# RUN go install google.golang.org/grpc && \
# 	go install -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
# 	go install -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
# 	go install -u github.com/golang/protobuf/protoc-gen-go

# install and build go app
COPY . /go/src/github.com/dhrp/moulin/

ARG APP_VERSION

# RUN make -C protobuf
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o moulin \
	-ldflags "-X main.version=${APP_VERSION}" \
	cmd/moulin/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o moulin-cli \
	-ldflags "-X main.version=${APP_VERSION}" \
	cmd/miller/main.go

FROM ubuntu:latest
WORKDIR /
COPY --from=0 /go/src/github.com/dhrp/moulin/moulin /usr/local/bin/moulin
COPY --from=0 /go/src/github.com/dhrp/moulin/moulin-cli /usr/local/bin/moulin-cli

ENV REDIS_HOST="localhost"
EXPOSE 8042
CMD ["moulin"]
