#!/bin/bash

# HTTPS Demo Script for MBankingCore API
# This script demonstrates both HTTP and HTTPS modes

echo "🔐 HTTPS Setup Demo for MBankingCore API"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

echo ""
print_info "This demo shows both HTTP and HTTPS modes of the MBankingCore API"
echo ""

# Check if binary exists
if [ ! -f "./mbankingcore" ]; then
    print_error "Binary not found. Building application..."
    go build -o mbankingcore
    if [ $? -eq 0 ]; then
        print_status "Application built successfully"
    else
        print_error "Failed to build application"
        exit 1
    fi
fi

# Stop any existing servers
pkill -f mbankingcore 2>/dev/null

echo ""
print_info "=== DEMO 1: HTTPS Mode (Secure) ==="
echo ""

print_info "Starting HTTPS server on port 8443..."
ENABLE_HTTPS=true ./mbankingcore &
HTTPS_PID=$!

# Wait for server to start
sleep 3

print_info "Testing HTTPS endpoints..."

# Test HTTPS health endpoint
echo ""
echo "🔒 Testing HTTPS Health Endpoint:"
echo "curl -k https://localhost:8443/health"
HTTPS_HEALTH=$(curl -s -k https://localhost:8443/health)
if [[ $HTTPS_HEALTH == *"MBankingCore API is running"* ]]; then
    print_status "HTTPS Health endpoint working!"
    echo "Response: $HTTPS_HEALTH"
else
    print_error "HTTPS Health endpoint failed"
fi

# Test HTTPS API endpoint
echo ""
echo "🔒 Testing HTTPS API Endpoint:"
echo "curl -k https://localhost:8443/api/onboardings"
HTTPS_API=$(curl -s -k https://localhost:8443/api/onboardings | head -c 100)
if [[ $HTTPS_API == *"Onboardings retrieved successfully"* ]]; then
    print_status "HTTPS API endpoint working!"
    echo "Response: ${HTTPS_API}..."
else
    print_error "HTTPS API endpoint failed"
fi

# Show certificate info
echo ""
print_info "SSL Certificate Information:"
if [ -f "certs/server.crt" ]; then
    echo "📄 Certificate: $(openssl x509 -in certs/server.crt -noout -subject 2>/dev/null || echo 'Certificate details unavailable')"
    echo "📅 Valid until: $(openssl x509 -in certs/server.crt -noout -dates 2>/dev/null | grep notAfter | cut -d= -f2 || echo 'Date unavailable')"
    print_status "Self-signed SSL certificate generated automatically"
else
    print_warning "Certificate files not found"
fi

# Stop HTTPS server
kill $HTTPS_PID 2>/dev/null
sleep 2

echo ""
print_info "=== DEMO 2: HTTP Mode (Fallback) ==="
echo ""

print_info "Starting HTTP server on port 8080..."
ENABLE_HTTPS=false ./mbankingcore &
HTTP_PID=$!

# Wait for server to start
sleep 3

print_info "Testing HTTP endpoints..."

# Test HTTP health endpoint
echo ""
echo "🌐 Testing HTTP Health Endpoint:"
echo "curl http://localhost:8080/health"
HTTP_HEALTH=$(curl -s http://localhost:8080/health)
if [[ $HTTP_HEALTH == *"MBankingCore API is running"* ]]; then
    print_status "HTTP Health endpoint working!"
    echo "Response: $HTTP_HEALTH"
else
    print_error "HTTP Health endpoint failed"
fi

# Test HTTP API endpoint
echo ""
echo "🌐 Testing HTTP API Endpoint:"
echo "curl http://localhost:8080/api/onboardings"
HTTP_API=$(curl -s http://localhost:8080/api/onboardings | head -c 100)
if [[ $HTTP_API == *"Onboardings retrieved successfully"* ]]; then
    print_status "HTTP API endpoint working!"
    echo "Response: ${HTTP_API}..."
else
    print_error "HTTP API endpoint failed"
fi

# Stop HTTP server
kill $HTTP_PID 2>/dev/null
sleep 1

echo ""
print_info "=== Configuration Summary ==="
echo ""

echo "Environment Variables:"
echo "• ENABLE_HTTPS=true   → Enables HTTPS on port 8443"
echo "• ENABLE_HTTPS=false  → Enables HTTP on port 8080 (default)"
echo "• HTTPS_PORT=8443     → Custom HTTPS port"
echo "• HTTP_PORT=8080      → Custom HTTP port"
echo "• CERT_DIR=./certs    → SSL certificate directory"

echo ""
echo "Generated Files:"
if [ -f "certs/server.crt" ]; then
    echo "• certs/server.crt    → SSL Certificate ($(stat -f%z certs/server.crt 2>/dev/null || stat -c%s certs/server.crt 2>/dev/null || echo 'unknown') bytes)"
fi
if [ -f "certs/server.key" ]; then
    echo "• certs/server.key    → Private Key ($(stat -f%z certs/server.key 2>/dev/null || stat -c%s certs/server.key 2>/dev/null || echo 'unknown') bytes)"
fi
echo "• mbankingcore        → Application Binary ($(stat -f%z mbankingcore 2>/dev/null || stat -c%s mbankingcore 2>/dev/null || echo 'unknown') bytes)"

echo ""
print_info "=== Flutter Integration URLs ==="
echo ""

echo "For Flutter applications, use these base URLs:"
echo ""
echo "📱 Android Emulator:"
echo "   HTTPS: https://10.0.2.2:8443/api"
echo "   HTTP:  http://10.0.2.2:8080/api"
echo ""
echo "📱 iOS Simulator:"
echo "   HTTPS: https://localhost:8443/api"
echo "   HTTP:  http://localhost:8080/api"
echo ""
echo "📱 Physical Device:"
echo "   HTTPS: https://$(ifconfig | grep 'inet ' | grep -v 127.0.0.1 | head -1 | awk '{print $2}'):8443/api"
echo "   HTTP:  http://$(ifconfig | grep 'inet ' | grep -v 127.0.0.1 | head -1 | awk '{print $2}'):8080/api"

echo ""
print_status "HTTPS setup complete! 🎉"
echo ""
print_info "To start the server:"
echo "• HTTPS: ENABLE_HTTPS=true ./mbankingcore"
echo "• HTTP:  ENABLE_HTTPS=false ./mbankingcore (or just ./mbankingcore)"
echo ""
print_info "For production, replace self-signed certificates with proper SSL certificates from a trusted CA."
