REPORTS=.reports
CLIENT_BIN=bin
CLIENT_DIR=cmd/client
CLIENT_V=v0.0.1
BUILD_DATE=$$(date +'%Y/%m/%d %H:%M:%S')
CLIENT_BUILD_FLAGS=-ldflags "-X main.Version=$(CLIENT_V) -X 'main.BuildDate=$(BUILD_DATE)'"
DB_URL=postgresql://admin:admin@localhost:5432/gophkeeper?sslmode=disable

$(REPORTS):
	mkdir $(REPORTS)

$(CLIENT_BIN):
	mkdir $(CLIENT_BIN)

setup: $(REPORTS) $(CLIENT_BIN) .tidy .install-tools  ## Установка окружения

clean:  # Очистка окружения
	rm -rf $(REPORTS)
	go clean -r -i -testcache -modcache

format:  ## Запуск форматтеров
	golangci-lint run --fix ./...

lint:  ## Последовательный запуск всех линтеров
	golangci-lint run ./... || golangci-lint run ./... --out-format json > $(REPORTS)/golangci-lint.json

test:  ## Запуск тестов
	go test -race -coverpkg=./... -coverprofile=$(REPORTS)/coverage.out -v 2>&1 ./... | go-junit-report -out $(REPORTS)/junit.xml -iocopy -set-exit-code
	go tool cover -html=$(REPORTS)/coverage.out -o $(REPORTS)/coverage.html
	go tool cover -func $(REPORTS)/coverage.out

.tidy:
	go mod tidy

.install-tools:
	 cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

all:  ## Запуск форматтеров, линтеров и тестов
	format lint test

migration-up:  ## Установка миграций бд сервера
	migrate -path ./migrations -database $(DB_URL) -verbose up

migration-down:  ## Откат миграций бд сервера
	migrate -path ./migrations -database $(DB_URL) -verbose down

migration-fix:  ## Фиксация миграций бд сервера (в случае неудачной установки/отката)
	migrate -path ./migrations -database $(DB_URL) force 1

build-client:  ## Сборка бинарников клиента под разные платформы
	rm -rf $(CLIENT_BIN)
	GOOS=linux GOARCH=amd64 go build $(CLIENT_BUILD_FLAGS) -o $(CLIENT_BIN)/cli_linux_amd64 $(CLIENT_DIR)/main.go
	GOOS=windows GOARCH=amd64 go build $(CLIENT_BUILD_FLAGS) -o $(CLIENT_BIN)/cli_windows_amd64.exe $(CLIENT_DIR)/main.go
	GOOS=darwin GOARCH=amd64 go build $(CLIENT_BUILD_FLAGS) -o $(CLIENT_BIN)/cli_darwin_amd64 $(CLIENT_DIR)/main.go
	GOOS=darwin GOARCH=arm64 go build $(CLIENT_BUILD_FLAGS) -o $(CLIENT_BIN)/cli_darwin_arm64 $(CLIENT_DIR)/main.go
	cp $(CLIENT_DIR)/config.json $(CLIENT_BIN)/config.json

generate-spec:  ## Генерация swagger spec
	swagger generate spec -o ./api/swagger.json

help:  ## Показать описание всех команд
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := all
