GODOG_VERSION=v0.11.0
GODOG := github.com/cucumber/godog/cmd/godog@$(GODOG_VERSION)
GOLANGCI_VERSION=v1.39.0
GOLANGCI := github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)
GOPUML_VERSION=v0.2.1
GOPUML := github.com/lonnblad/gopuml/cmd/gopuml@$(GOPUML_VERSION)
SWAG_GO_VERSION := v1.7.0
SWAG_GO := github.com/swaggo/swag/cmd/swag@$(SWAG_GO_VERSION)
GO_SWAGGER_VERSION := v0.27.0
GO_SWAGGER := github.com/go-swagger/go-swagger/cmd/swagger@$(GO_SWAGGER_VERSION)
TOOLS := $(GODOG) $(GOLANGCI) $(GOPUML) $(SWAG_GO) $(GO_SWAGGER)

SERVICE_API_V1_SPEC_SOURCES := $(shell find boundaries/rest/v1 -type f -name '*.go')
SERVICE_API_V1_SPEC_TARGET := boundaries/rest/v1/generated/swagger/swagger.json

UML_DOC_GEN_SOURCES :=  $(shell find docs/diagrams -type f -name '*.puml')
UML_DOC_GEN_TARGET :=  $(shell find docs/diagrams -type f -name '*.puml' | xargs -0 -n1 echo | sed "s/puml/svg/" )

BACKEND_GO_FILES := $(shell find . -type f -name '*.go')

BINARIES := \
	cmd/rest-api/main

$(TOOLS):
	@echo "[tools] Installing $@"
	@go install $@

.PHONY: gen
gen: $(SERVICE_API_V1_SPEC_TARGET) $(UML_DOC_GEN_TARGET)

$(SERVICE_API_V1_SPEC_TARGET): $(SERVICE_API_V1_SPEC_SOURCES) | $(SWAG_GO) $(GO_SWAGGER)
	@echo "[swag] Generating swagger documentation for API v1"
	@cd boundaries/rest/v1 && $(shell go env GOPATH)/bin/swag init --parseDepth 1 --parseDependency -o generated/swagger -g v1.go
	@echo "[go-swagger] Validating generated swagger documentation for API v1"
	@$(shell go env GOPATH)/bin/swagger validate $@


$(UML_DOC_GEN_TARGET): $(UML_DOC_GEN_SOURCES) | $(GOPUML)
	@echo "[gopuml] Generating UML diagrams"
	@rm $(UML_DOC_GEN_TARGET)
	@$(shell go env GOPATH)/bin/gopuml build $(UML_DOC_GEN_SOURCES)

.PHONY: run-local
run-local: gen
	@env SERVICE_NAME=shipment-service \
		SERVICE_VERSION=dev \
		ENVIRONMENT=local \
	go run cmd/rest-api/main.go

.PHONY: run-behaviour-test
run-behaviour-test: gen | $(GODOG)
	@cd cmd/rest-api-test && \
		env MODE=unit-test \
		godog run ../../behaviour

.PHONY: build
build: $(BINARIES)

$(BINARIES): cmd/%/main: gen $(BACKEND_GO_FILES)
	@echo "[golang] Building" $@
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-w -extldflags "-static"' -o ./$@ ./$@.go

.PHONY: fmt
fmt: go-fmt

.PHONY: go-fmt
go-fmt: gen
	@echo "[golang] Cleaning dependencies"
	@go mod tidy
	@echo "[golang] Formatting go code"
	@go fmt ./...
	@$(shell go env GOPATH)/bin/goimports -w -local github.com/lonnblad .
	
.PHONY: lint
lint: go-lint

.PHONY: go-lint
go-lint: gen $(GOLANGCI)
	@echo "[golang] Running golangci-lint"
	@$(shell go env GOPATH)/bin/golangci-lint run ./...

.PHONY: test
test: go-unit

.PHONY: go-unit
go-unit: gen
	@echo "[golang] Running unit tests"
	@go test -v -short --race ./...

.PHONY: clean
clean:
	@rm -rf $(BINARIES)
	@rm -rf boundaries/rest/v1/generated
