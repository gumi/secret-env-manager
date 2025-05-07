#!/bin/bash

# Functions for colored output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Log output functions
log_success() {
  echo "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
  echo "${RED}[ERROR]${NC} $1"
}

log_info() {
  echo "${YELLOW}[INFO]${NC} $1"
}

# Test execution function
run_test() {
  base_name=$1
  input_file="tests/input/$base_name"
  expected_file="tests/expected/$base_name"
  output_file=".cache.tests_input_$base_name"
  
  # Check if expected result file exists
  if [ ! -f "$expected_file" ]; then
    log_error "Expected result file does not exist: $expected_file"
    return 1
  fi
  
  log_info "Running test: $input_file (expected: $expected_file)"
  
  # Execute update command
  go run main.go update -i "$input_file" --endpoint-url http://localhost:4566
  
  if [ $? -ne 0 ]; then
    log_error "Failed to execute update command: $input_file"
    return 1
  fi

  # Check if output file exists
  if [ ! -f "$output_file" ]; then
    log_error "Output file was not generated: $output_file"
    return 1
  fi

  # Execute load command and capture its output
  load_output=$(go run main.go load -i "$input_file")
  
  if [ $? -ne 0 ]; then
    log_error "Failed to execute load command: $input_file"
    return 1
  fi

  # Compare load command output with expected result
  expected_content=$(cat "$expected_file")
  if [ "$load_output" == "$expected_content" ]; then
    log_success "Test passed: $input_file"
    return 0
  else
    log_error "Test failed: $input_file - Load output does not match expected"
    echo "Expected:"
    echo "$expected_content"
    echo "Got:"
    echo "$load_output"
    return 1
  fi
}

# Main process
main() {
  # Check project root directory
  if [ ! -f "main.go" ] || [ ! -d "tests" ]; then
    log_error "Please run this script from the project root directory"
    log_error "Current directory: $(pwd)"
    exit 1
  fi
  
  log_info "Please ensure LocalStack is running and secrets are initialized"
  
  # Dynamically find all input env files
  input_files=()
  for file in $(find tests/input -name "*.env" -type f); do
    # Extract just the filename without path
    filename=$(basename "$file")
    input_files+=("$filename")
  done
  
  # Test variables
  total_tests=${#input_files[@]}
  passed_tests=0
  failed_tests=0
  
  log_info "Found $total_tests test files to process"
  
  # Process each test file
  for file in "${input_files[@]}"; do
    if run_test "$file"; then
      ((passed_tests++))
    else
      ((failed_tests++))
    fi
    echo "----------------------------------------"
  done
  
  # Output result summary
  echo "Test Results Summary:"
  echo "  Total Tests: $total_tests"
  echo "  Passed: ${GREEN}$passed_tests${NC}"
  echo "  Failed: ${RED}$failed_tests${NC}"
  
  if [ $failed_tests -eq 0 ]; then
    log_success "All tests passed!"
    exit 0
  else
    log_error "$failed_tests test(s) failed"
    exit 1
  fi
}

# Run script
main