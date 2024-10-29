#!/bin/bash

set -e

echo "Building ./dist/overlay..."
go build -o dist ./cmd/overlay
echo "Done."