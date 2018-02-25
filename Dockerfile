FROM golang:1.9 AS builder

ENV DEP_VERSION="0.4.1"
ENV PROTOC_VERSION="3.5.1"

RUN apt-get update && \
	apt-get install unzip &&\
	rm -rf /var/lib/apt/lists/

# Install Go Dep
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 && chmod +x /usr/local/bin/dep

# install protoc
RUN curl -L https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > protoc.zip && \						
	unzip protoc.zip -d .

RUN go get google.golang.org/grpc && \
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
		go get -u github.com/golang/protobuf/protoc-gen-go

# Install dependencies
COPY Gopkg.toml Gopkg.lock /go/src/github.com/dhrp/moulin/
WORKDIR /go/src/github.com/dhrp/moulin/
RUN dep ensure -vendor-only

# install and build go app
COPY . /go/src/github.com/dhrp/moulin/

RUN make -C protobuf
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/moulin *.go


# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /go/bin/
# COPY --from=0 /go/bin/moulin /go/bin/moulin
# COPY --from=0 /go/src/github.com/dhrp/moulin/certificates/* /go/bin/certificates/

EXPOSE 8042
ENTRYPOINT ["/go/bin/moulin"]
