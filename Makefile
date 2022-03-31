# Output directory
OUTDIR = ./bin
OBJDIR = $(OUTDIR)/go-get-started
BINARY = webapp

build:
	mkdir -p $(OBJDIR)
	go build -o $(OBJDIR)	./...

test:
	go test -race ./...

run: build
	$(OBJDIR)/$(BINARY)

clean:
	rm -rf $(OUTDIR)