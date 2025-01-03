#!/bin/bash

set -e

# Ensure gocovmerge is installed
if ! command -v gocovmerge &> /dev/null; then
    echo "gocovmerge could not be found, installing..."
    go install github.com/wadey/gocovmerge@latest
fi

# Remove any existing coverage files
rm -rf ./coverage

# Create a directory to store individual coverage profiles
mkdir -p coverage

# Find all packages and run tests with coverage
for pkg in $(go list ./...); do
    go test -coverprofile=coverage/$(echo $pkg | tr / -).out -coverpkg=./... $pkg
done

# Merge all coverage profiles into one
gocovmerge coverage/*.out > ./coverage/coverage.out

go tool cover -func=./coverage/coverage.out | grep "total:"

# Generate an HTML report
go tool cover -html=./coverage/coverage.out -o coverage/coverage.html

if [ "$1" != "quiet" ]; then
  nohup open coverage/coverage.html > /dev/null 2>&1&
fi
 
