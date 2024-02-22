REPORTS=.reports

$(REPORTS):
	mkdir $(REPORTS)

setup: $(REPORTS) .tidy .install-tools

clean:
	rm -rf $(REPORTS)
	go clean -r -i -testcache -modcache

format:
	golangci-lint run --fix ./...

lint:
	golangci-lint run ./... --verbose --out-format json > $(REPORTS)/golangci-lint.json

test:
	go test -race -coverprofile=$(REPORTS)/coverage.out -v 2>&1 ./... | go-junit-report -out $(REPORTS)/junit.xml -iocopy -set-exit-code
	go tool cover -html=$(REPORTS)/coverage.out -o $(REPORTS)/coverage.html

.tidy:
	go mod tidy

.install-tools:
	 cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

all: format lint test

.DEFAULT_GOAL := all
