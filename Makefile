BIN_DIR=bin

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

build:
	make vendor
	go build -o ./bin/wt ./main.go

build-cross-platform:
	make vendor
	env CGO_ENABLED=0 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(BIN_DIR)/ --out wt .

build-project-archive:
	tar -czvf $(BIN_DIR)/wt.tar.gz $(BIN_DIR)

lint:
	golint -set_exit_status ./services/...

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

prepare-release:
	make build-cross-platform
	make build-project-archive