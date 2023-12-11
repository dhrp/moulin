BINARY=moulin
CLI_BINARY=moulin-cli

VERSION=1.0.0
BUILD=`git rev-parse HEAD`

# ToDo: set versions stuffs in files
# Setup the -ldflags option for go build here, interpolate the variable values
# LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

IMAGE_TAG=v0.4.1


build:
	go build -o ${BINARY} cmd/moulin/*.go

cli:
	go build -o ${CLI_BINARY} cmd/miller/*.go

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
	docker build -t dhrp/moulin:$(IMAGE_TAG) .

publish:
	docker buildx build --platform linux/amd64 -t dhrp/moulin:$(IMAGE_TAG) --push .

docker-run:
	docker run --link redis:redis -e REDIS_ADDRESS=redis:6379 --name moulin -d dhrp/moulin:$(IMAGE_TAG)

docker-run-cli:
	docker run --link redis:redis -e MOULIN_SERVER=moulin:8042 dhrp/moulin:$(IMAGE_TAG)

docker-debug-cli:
	docker run -it --link moulin:moulin -e MOULIN_SERVER=moulin:8042 dhrp/moulin:$(IMAGE_TAG) bash

push:
	docker push dhrp/moulin:$(IMAGE_TAG)

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

.PHONY: build run test install clean image push cli

redis:
	@echo "starting redis if not already running"
	docker start redis || docker run -p 6379:6379 --name redis -d redis && sleep 4

