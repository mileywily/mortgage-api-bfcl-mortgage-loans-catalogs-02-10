BUILDPATH=$(CURDIR)

PACKAGES := $(shell go list ./... | grep -v '/lib/')

build:
	@echo "Creando Binario ..."
	@mkdir -p $(BUILDPATH)/build/bin
	@go build -ldflags '-s -w' -o $(BUILDPATH)/build/bin/dist ./cmd/api
	@echo "Binario generado en build/bin/dist"

test:
	@echo "Ejecutando tests..."
	@go test -coverpkg=./internal/... $(PACKAGES)

coverage:
	@echo "Coverfile..."
	go test -coverpkg=./internal/... -coverprofile=coverfile_raw.out $(PACKAGES)
	@grep -v -E "/cmd/|dummy_repository.go|/mocks/" coverfile_raw.out > coverfile_out
	@go tool cover -func coverfile_out
	@go tool cover -func coverfile_out | grep total | grep -o '[0-9]*\.[0-9]*' | cut -d' ' -f1 > coverage.txt

.PHONY: test build coverage

