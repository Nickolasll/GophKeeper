REPORTS=.reports
DB_URL=postgresql://admin:admin@localhost:5432/gophkeeper?sslmode=disable

$(REPORTS):
	mkdir $(REPORTS)

setup: 
	$(REPORTS) .tidy .install-tools

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

.DEFAULT_GOAL := all
