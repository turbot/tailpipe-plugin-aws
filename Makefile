# Paths
PLUGIN_NAME=tailpipe-plugin-aws.plugin
PLUGIN_DIR=~/.tailpipe/plugins

# Build in development mode by default
.PHONY: default
default: install

# Production build, optimized
.PHONY: build
build:
	go build -o $(PLUGIN_NAME) .

# Install the development build
.PHONY: install
install: build
	mv $(PLUGIN_NAME) $(PLUGIN_DIR)

# Run tests
.PHONY: test
test:
	go test ./... -v

# Clean up generated files
.PHONY: clean
clean:
	rm -f $(PLUGIN_NAME)

