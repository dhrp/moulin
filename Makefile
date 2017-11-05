BINARY=moulin

VERSION=1.0.0
BUILD=`git rev-parse HEAD`

# ToDo: set verions stuffs in files
# Setup the -ldflags option for go build here, interpolate the variable values
# LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"


build:
	go build -o ${BINARY} *.go

run:
	./moulin &

test:
	go test github.com/nerdalize/moulin/rouge
	go test github.com/nerdalize/moulin

test-client: build
	{ ./moulin & }; \
	pid=$$!; \
	go test github.com/nerdalize/moulin/client; \
	r=$$?; \
	kill $$pid; \
	exit $$r

tests: test test-client

install:
	go install

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

.PHONY: build run test install clean