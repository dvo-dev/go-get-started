# Output directory
OUTDIR = ./bin
OBJDIR = $(OUTDIR)/go-get-started
BINARY = webapp

build:
	mkdir -p $(OBJDIR)
	go build -o $(OBJDIR)	./...

lint:
	go vet ./...
	golangci-lint run ./...

test:
	go clean -testcache
	go test -race ./...

run: build
	$(OBJDIR)/$(BINARY)

clean:
	rm -rf $(OUTDIR)