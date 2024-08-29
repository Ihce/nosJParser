# Define the output directory for the binary
BINDIR := bin

# Define the name of the binary
BINARY := nosj_parser

# Define the source directory for the main package
SRCDIR := ./cmd/nosj_parser

# Define the build target
build:
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(BINARY) $(SRCDIR)

# Define the run target
run: build
	./$(BINDIR)/$(BINARY)

# Define the clean target
clean:
	rm -rf $(BINDIR)

# Define the default target
.PHONY: all
all: build