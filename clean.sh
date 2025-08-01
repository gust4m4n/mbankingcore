#!/bin/bash

# Clean workspace script for MBankingCore Go project
# This script removes temporary files, caches, and build artifacts

set -e

echo "🧹 Cleaning MBankingCore workspace..."

# Remove build artifacts
echo "🗑️  Removing build artifacts..."
rm -f mbankingcore
rm -f migrate
rm -rf bin/
rm -rf dist/
rm -rf build/

# Remove Go build cache
echo "🗑️  Cleaning Go build cache..."
go clean -cache -testcache -modcache

# Remove temporary files
echo "🗑️  Removing temporary files..."
find . -name "*.tmp" -delete 2>/dev/null || true
find . -name "*.temp" -delete 2>/dev/null || true
find . -name "*.log" -delete 2>/dev/null || true
find . -name "*~" -delete 2>/dev/null || true
find . -name "*.bak" -delete 2>/dev/null || true
find . -name "*.backup" -delete 2>/dev/null || true

# Remove OS generated files
echo "🗑️  Removing OS generated files..."
find . -name ".DS_Store" -delete 2>/dev/null || true
find . -name "._*" -delete 2>/dev/null || true
find . -name "Thumbs.db" -delete 2>/dev/null || true
find . -name "ehthumbs.db" -delete 2>/dev/null || true

# Remove editor temporary files
echo "🗑️  Removing editor temporary files..."
find . -name "*.swp" -delete 2>/dev/null || true
find . -name "*.swo" -delete 2>/dev/null || true

# Remove test artifacts
echo "🗑️  Removing test artifacts..."
find . -name "*.test" -delete 2>/dev/null || true
find . -name "*.out" -delete 2>/dev/null || true
find . -name "coverage.*" -delete 2>/dev/null || true

# Remove debug binaries
echo "🗑️  Removing debug binaries..."
find . -name "__debug_bin*" -delete 2>/dev/null || true
find . -name "*.debug" -delete 2>/dev/null || true

# Recreate necessary directories
echo "📁 Recreating necessary directories..."
mkdir -p bin/
mkdir -p certs/

echo "✅ Workspace cleaned successfully!"
echo "📂 Build directory recreated: ./bin/"
echo "🔐 Certificates directory recreated: ./certs/"
echo "🚀 Ready for clean build with: ./build.sh"
