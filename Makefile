PROJECT=wordlistdb

default: help

.PHONY: help
help: ## list makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:
	@go build -o $(PROJECT) -v cmd/$(PROJECT)/main.go

.PHONY: fmtcheck
fmtcheck: ## run gofmt and print detected files
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: clean
clean:
	@sh -c "rm ./$(PROJECT)"
	go mod tidy
	rm ./*.rpm ./*.apk ./*.deb

.PHONY: dist
dist:
	@nfpm package -p deb
	nfpm package -p apk
	nfpm package -p rpm

PHONY: test
test: ## run go tests
	go test -v ./...

PHONY: fmt
fmt: ## format go files
	gofumpt -w .

.PHONY: pre-commit
pre-commit:	## run pre-commit hooks
	pre-commit run
