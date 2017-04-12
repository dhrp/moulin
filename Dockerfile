from golang:1

RUN apt-get update && \
	apt-get install unzip &&\
	rm -rf /var/lib/apt/lists/

# install go dependencies
RUN go get github.com/mediocregopher/radix.v2


# install protoc
WORKDIR protobuf
RUN curl -L https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip > protoc.zip && \
	unzip protoc.zip -d . && \
	cp bin/protoc /go/bin/
RUN go get google.golang.org/grpc

# install and build go app
COPY . /go/src/github.com/nerdalize/moulin/
RUN go build -o /go/bin/server /go/src/github.com/nerdalize/moulin/server/main.go 

EXPOSE 50051
ENTRYPOINT ["/go/bin/server"]
