.PHONY: default test lint vet build .build build-with-docker

GO=GO15VENDOREXPERIMENT=1 go

default: test-with-docker

test: vet #lint
	$(GO) test -v . ./cmd/...

lint:
	golint .
	golint ./cmd/...

vet:
	$(GO) vet . ./cmd/...

build: build-with-docker

.build: test
	$(GO) build -o build/robin ./cmd/robin

build-with-docker:
	mkdir -p build && \
	docker run --rm \
		-v $$(pwd):/go/src/github.com/pilu/robin \
		-w /go/src/github.com/pilu/robin \
		gravityblast/go-build \
		make .build

test-with-docker:
	mkdir -p build && \
	docker run --rm \
		-v $$(pwd):/go/src/github.com/pilu/robin \
		-w /go/src/github.com/pilu/robin \
		gravityblast/go-build \
		make test
