# Binary names and generated baseline location
BIN_DIR=bin
GENERATOR_BINARY=$(BIN_DIR)/generate-baseline
SCANNER_BINARY=$(BIN_DIR)/secure_fim
BASELINE_FILE=internal/baseline/baseline.json

.PHONY: all generate build build-generator build-scanner clean test

# Default target runs everything in order
all: generate build

# Generate the baseline embedded by the scanner
generate:
	@echo "[*] Running the dynamic generator..."
	go run ./cmd/generate-baseline
	@echo "[*] Checking baseline output..."
	@cat $(BASELINE_FILE)

# Build each CLI independently
build: build-generator build-scanner

build-generator:
	@mkdir -p $(BIN_DIR)
	go build -o $(GENERATOR_BINARY) ./cmd/generate-baseline

build-scanner:
	@mkdir -p $(BIN_DIR)
	go build -o $(SCANNER_BINARY) ./cmd/secure-fim
	@echo "[+] Scanner compiled as ./$(SCANNER_BINARY)"

# Clean up build artifacts
clean:
	@echo "[*] Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	@echo "[+] Clean complete."

# Helper target to quickly run a test generation step
test: generate
	@echo "[*] Verifying baseline target counts..."
	@wc -l $(BASELINE_FILE)
