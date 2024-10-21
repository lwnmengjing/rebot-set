PROJECT:=mss-boot-admin

.PHONY: test

test:
	go test -coverprofile=coverage.out ./...
deps:
	go mod download
generate:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run -v ./...
fix-lint:
	goimports -w .