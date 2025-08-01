#!/bin/bash

# Build script for MBankingCore
# This builds the binary in a separate directory to keep workspace clean

set -e

echo "🧹 Cleaning previous builds..."
rm -f ./mbankingcore
go clean -cache

echo "🏗️  Building MBankingCore..."
go build -o ./bin/mbankingcore

echo "✅ Build completed successfully!"
echo "📍 Binary location: ./bin/mbankingcore"
echo "🚀 To run: ./bin/mbankingcore"
