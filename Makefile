BIN_DIR=bin

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

build:
	go build -o ./bin/wt ./main.go

build-cross-platform:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(BIN_DIR)/ --out wt .

build-project-archive:
	tar -czvf $(BIN_DIR)/wt.tar.gz $(BIN_DIR)

prepare-release:
	make build-cross-platform
	make build-project-archive