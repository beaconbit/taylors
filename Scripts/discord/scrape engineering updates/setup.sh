#!/bin/bash
# Setup script for Go Discord Message Scraper

echo "=== Go Discord Message Scraper Setup ==="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

echo "1. Installing Go modules..."
go mod download

echo "2. Building the scraper..."
go build -o discord-scraper main.go

