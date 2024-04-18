# Paths
PLUGIN_NAME=tailpipe-plugin-aws.so
PLUGIN_DIR=~/.tailpipe/plugins

# Build flags
PROD_FLAGS=-buildmode=plugin
DEV_FLAGS=-buildmode=plugin -gcflags='all=-N -l'

# Build in development mode by default
.PHONY: default
default: dev-install

# Production build, optimized
.PHONY: build
build:
	go build $(PROD_FLAGS) -o $(PLUGIN_NAME) .

# Development build, removes optimizations to allow debugging
.PHONY: dev
dev:
	go build $(DEV_FLAGS) -o $(PLUGIN_NAME) .

# Install the development build
.PHONY: dev-install
dev-install: dev
	cp $(PLUGIN_NAME) $(PLUGIN_DIR)

# Run tests
.PHONY: test
test:
	go test ./... -v

# Clean up generated files
.PHONY: clean
clean:
	rm -f $(PLUGIN_NAME)
