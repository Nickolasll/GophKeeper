REPORTS=.reports
CLIENT_BIN=bin
CLIENT_DIR=cmd/client
DB_URL=postgresql://admin:admin@localhost:5432/gophkeeper?sslmode=disable

$(REPORTS):
	mkdir $(REPORTS)

$(CLIENT_BIN):
	mkdir $(CLIENT_BIN)

setup: 
	$(REPORTS) $(CLIENT_BIN) .tidy .install-tools

clean:
	rm -rf $(REPORTS)
	go clean -r -i -testcache -modcache

format:
	golangci-lint run --fix ./...

lint:
	golangci-lint run ./... || golangci-lint run ./... --out-format json > $(REPORTS)/golangci-lint.json

test:
	go test -race -coverpkg=./... -coverprofile=$(REPORTS)/coverage.out -v 2>&1 ./... | go-junit-report -out $(REPORTS)/junit.xml -iocopy -set-exit-code
	go tool cover -html=$(REPORTS)/coverage.out -o $(REPORTS)/coverage.html
	go tool cover -func $(REPORTS)/coverage.out

.tidy:
	go mod tidy

.install-tools:
	 cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

all: 
	format lint test

migration-up: 
	migrate -path ./migrations -database $(DB_URL) -verbose up

migration-down: 
	migrate -path ./migrations -database $(DB_URL) -verbose down

migration-fix: 
	migrate -path ./migrations -database $(DB_URL) force 1

build-client:
	rm -rf $(CLIENT_BIN)
	GOOS=linux GOARCH=amd64 go build -o $(CLIENT_BIN)/cli_linux_amd64 $(CLIENT_DIR)/main.go
	GOOS=windows GOARCH=amd64 go build -o $(CLIENT_BIN)/cli_windows_amd64.exe $(CLIENT_DIR)/main.go
	GOOS=darwin GOARCH=amd64 go build -o $(CLIENT_BIN)/cli_darwin_amd64 $(CLIENT_DIR)/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(CLIENT_BIN)/cli_darwin_arm64 $(CLIENT_DIR)/main.go
	cp $(CLIENT_DIR)/ca.crt $(CLIENT_BIN)/ca.crt

.DEFAULT_GOAL := all
