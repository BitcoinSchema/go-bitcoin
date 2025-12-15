#!/usr/bin/env bash
# ------------------------------------------------------------------------------------
# Parse Test Label Script
#
# Helper function to generate human-readable test labels from artifact names.
# Sourced by workflow steps that need consistent test labeling.
#
# Usage: parse_test_label "artifact-name" "jsonl-filename"
# Output: "Unit Tests (Ubuntu, Go 1.22)" or similar
# ------------------------------------------------------------------------------------

parse_test_label() {
  local artifact_name="$1"
  local jsonl_name="$2"

  # Determine test type from artifact prefix or JSONL name
  local test_type="Tests"
  if [[ "$artifact_name" == test-results-fuzz-* ]] || [[ "$jsonl_name" == *fuzz* ]]; then
    test_type="Fuzz Tests"
  elif [[ "$artifact_name" == ci-results-* ]]; then
    test_type="Unit Tests"
  fi

  # Extract OS from artifact name
  local os_name=""
  if [[ "$artifact_name" =~ ubuntu ]]; then
    os_name="Ubuntu"
  elif [[ "$artifact_name" =~ windows ]]; then
    os_name="Windows"
  elif [[ "$artifact_name" =~ macos ]]; then
    os_name="macOS"
  fi

  # Extract Go version (last segment like "1.22", "1.24.x", or "go1.22")
  local go_version=""
  go_version=$(echo "$artifact_name" | grep -oE '[0-9]+\.[0-9]+(\.[x0-9]+)?' | tail -1 || echo "")

  # Build label
  if [[ -n "$os_name" && -n "$go_version" ]]; then
    echo "$test_type ($os_name, Go $go_version)"
  elif [[ -n "$os_name" ]]; then
    echo "$test_type ($os_name)"
  elif [[ -n "$go_version" ]]; then
    echo "$test_type (Go $go_version)"
  else
    echo "$test_type"
  fi
}
