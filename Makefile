default: build

build: go_build docker_build

DOCKER_IMAGE ?= rossf7/annotator
DOCKER_TAG ?= $(shell cat VERSION)
BINARY ?= operator

go_build:
	GOOS=linux go build -o $(BINARY) cmd/operator/main.go

docker_build:
	docker build \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg VERSION=`cat VERSION` \
		--build-arg VCS_URL=`git config --get remote.origin.url` \
		--build-arg VCS_REF=`git rev-parse --short HEAD` \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) .

clean:
	rm $(BINARY)

test:
	go test $(shell go list ./...)
