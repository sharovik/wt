BIN_DIR=bin

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

build:
	make vendor
	go build -o ./bin/wt ./main.go

release:
	make vendor
	env CGO_ENABLED=0 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(BIN_DIR)/ --out wt .

lint:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./services/...
	staticcheck ./configuration/...
	staticcheck ./dto/...

imports:
	goimports -d -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(shell go list ./... | grep -v /vendor/)

tests:
	go test ./...

code-check:
	make lint
	make tests

code-clean:
	make imports
	make format

check-security:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...