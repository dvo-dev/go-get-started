# Output directory
OUTDIR = ./bin
OBJDIR = $(OUTDIR)/go-get-started
BINARY = webapp
TEST_COVERAGE = coverage.out

build:
	mkdir -p $(OBJDIR)
	go build -o $(OBJDIR)	./...

lint:
	go vet ./...
	golangci-lint run ./...

test:
	go clean -testcache
	go test -race ./...

test-coverage:
	go clean -testcache
	go test -race -coverprofile=$(TEST_COVERAGE) -covermode=atomic ./...
	go tool cover -func $(TEST_COVERAGE)

run: build
	$(OBJDIR)/$(BINARY)

clean:
	rm -rf $(OUTDIR)

db-pgcli:
    pgcli -h localhost -p 5432 -U admin app
