#!/bin/bash

set -e

echo "Building ./dist/overlay..."
go build -buildvcs=false -o dist ./cmd/overlay
echo "Done."