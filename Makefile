.PHONY: build
build: vet
	go build ./...

.PHONY: install
install: test
	go install ./...

.PHONY: vet
vet:
	 go vet ./...

.PHONY: lint
lint:
	 golint ./...

.PHONY: test
test: vet
	go test -race ./...
