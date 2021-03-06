BINARY=moulin
CLI_BINARY=moulin-cli

VERSION=1.0.0
BUILD=`git rev-parse HEAD`

# ToDo: set versions stuffs in files
# Setup the -ldflags option for go build here, interpolate the variable values
# LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


build:
	go build -o ${BINARY} server/*.go

cli:
	go build -o ${CLI_BINARY} cli/*.go

run:
	./moulin

test:
	go test ./... -count=1

test-client: build
	{ ./moulin & }; \
	pid=$$!; \
	go test github.com/dhrp/moulin/client; \
	r=$$?; \
	kill $$pid; \
	exit $$r

tests: test test-client

install:
	go install ./...

image:
	docker build -t dhrp/moulin .

push:
	docker push dhrp/moulin

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

.PHONY: build run test install clean image push cli
