# Define the name of the binary
BINARY := nosj_parser

# Define the source directory for the main package
SRCDIR := ./cmd/nosj_parser

# Define the build target
build:
	@go build -o $(BINARY) $(SRCDIR)

# Define the run target
run:
	@./$(BINARY) $(FILE)

# Define the clean target
clean:
	rm -rf $(BINDIR)

# Define the default target
.PHONY: all