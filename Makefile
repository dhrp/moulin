.PHONY: build run test install clean image push cli

BINARY=moulin
CLI_BINARY=moulin-cli
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
IMAGE_TAG=$(VERSION)

VERSION=$(shell git describe --tags --always --dirty)
DATE=$(shell date +%Y%m%d.%H%M)

define build-server
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) \
	go build -o ${BINARY}-$(1)-$(2) \
	-ldflags "-X main.Version=${VERSION} -X main.Build=${DATE}" \
	cmd/moulin/*.go
endef

define build-cli
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) \
	go build -o ${CLI_BINARY}-$(1)-$(2) \
	-ldflags "-X main.Version=${VERSION} -X main.Build=${DATE}" \
	cmd/miller/main.go
endef


build:
	$(call build-server,${GOOS},$(GOARCH))
	cp ${BINARY}-${GOOS}-${GOARCH} ${BINARY}

cli:
	$(call build-cli,${GOOS},$(GOARCH))
	cp ${CLI_BINARY}-${GOOS}-${GOARCH} ${CLI_BINARY}

build-all:
	$(call build-server,linux,amd64)
	$(call build-server,darwin,arm64)
	$(call build-cli,darwin,arm64)
	$(call build-cli,linux,amd64)

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
	docker build --platform linux/amd64 --build-arg APP_VERSION=$(VERSION) -t dhrp/moulin:$(IMAGE_TAG) .

publish:
	docker buildx build --platform linux/amd64 --build-arg APP_VERSION=$(VERSION) -t dhrp/moulin:$(IMAGE_TAG) --push .

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

redis:
	@echo "starting redis if not already running"
	docker start redis || docker run -p 6379:6379 --name redis -d redis && sleep 4

