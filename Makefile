MIGRATE_CMD = go run cmd/db.go
RUN_CMD = go run main.go
APP_NAME = textdiff
VERSION  = v1.0.0
BIN_DIR  = bin

.PHONY: all clean

migrate:
	${MIGRATE_CMD} migrate

reset:
	${MIGRATE_CMD} reset

init:
	${RUN_CMD} init

diff:
	${RUN_CMD} diff

commit:
	${RUN_CMD} commit --message="yo wassup man"

revert:
	${RUN_CMD} revert --back=3

log:
	${RUN_CMD} log

plain:
	${RUN_CMD}

test:
	go test -v -count=1 ./...

build: $(BIN_DIR)/$(APP_NAME)-linux \
       $(BIN_DIR)/$(APP_NAME)-macos-arm64 \
       $(BIN_DIR)/$(APP_NAME)-win.exe

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN_DIR)/$(APP_NAME)-linux: $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ ./main.go

$(BIN_DIR)/$(APP_NAME)-macos-arm64: $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $@ ./main.go

$(BIN_DIR)/$(APP_NAME)-win.exe: $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $@ ./main.go

zip: $(BIN_DIR)/$(APP_NAME)-linux.zip \
     $(BIN_DIR)/$(APP_NAME)-macos-arm64.zip \
     $(BIN_DIR)/$(APP_NAME)-win.zip

$(BIN_DIR)/$(APP_NAME)-linux.zip: $(BIN_DIR)/$(APP_NAME)-linux
	zip -j $@ $<

$(BIN_DIR)/$(APP_NAME)-macos-arm64.zip: $(BIN_DIR)/$(APP_NAME)-macos-arm64
	zip -j $@ $<

$(BIN_DIR)/$(APP_NAME)-win.zip: $(BIN_DIR)/$(APP_NAME)-win.exe
	zip -j $@ $<

clean:
	rm -rf $(BIN_DIR)
