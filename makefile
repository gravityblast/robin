.PHONY: default test lint vet build release build-with-docker test-with-docker

GO=GO15VENDOREXPERIMENT=1 go

default: test-with-docker

test: vet #lint
	$(GO) test -v . ./cmd/...

lint:
	golint .
	golint ./cmd/...

vet:
	$(GO) vet . ./cmd/...

build: test
	$(GO) build

release:
	$(GO) build -o ./release/robin ./cmd/robin

build-with-docker:
	mkdir -p build && \
	docker run --rm \
		-v $$(pwd):/go/src/github.com/pilu/robin \
		-w /go/src/github.com/pilu/robin \
		gravityblast/go-build \
		make build

test-with-docker:
	mkdir -p build && \
	docker run --rm \
		-v $$(pwd):/go/src/github.com/pilu/robin \
		-w /go/src/github.com/pilu/robin \
		gravityblast/go-build \
		make test
