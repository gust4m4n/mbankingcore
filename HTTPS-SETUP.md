# HTTPS Setup Guide for MBankingCore API

## 🔒 Overview

MBankingCore API now supports both HTTP and HTTPS modes for secure communication. This guide covers how to enable and configure HTTPS for development and production environments.

## 📋 Quick Start

### 1. Enable HTTPS in Environment

Create or update your `.env` file:

```bash
# Server Configuration
HOST=0.0.0.0
PORT=8080

# HTTPS Configuration
ENABLE_HTTPS=true
HTTPS_PORT=8443
CERT_DIR=./certs

# Other configurations...
```

### 2. Install OpenSSL (macOS)

```bash
# Install OpenSSL if not already installed
brew install openssl
```

### 3. Start the Server

```bash
# Build and run the application
go build -o mbankingcore
./mbankingcore
```

The server will automatically:
- ✅ Create `./certs` directory
- ✅ Generate self-signed SSL certificate
- ✅ Start HTTPS server on port 8443

## 🛠️ Configuration Options

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENABLE_HTTPS` | `false` | Enable HTTPS mode (`true`, `1`, `false`, `0`) |
| `HTTPS_PORT` | `8443` | HTTPS server port |
| `PORT` | `8080` | HTTP server port (fallback) |
| `HOST` | `0.0.0.0` | Server bind address |
| `CERT_DIR` | `./certs` | Directory for SSL certificates |

### SSL Certificate Files

- **Certificate**: `{CERT_DIR}/server.crt`
- **Private Key**: `{CERT_DIR}/server.key`

## 🔐 Certificate Management

### Auto-Generated Development Certificates

The server automatically generates self-signed certificates for development:

```bash
# Certificate details
Subject: /C=ID/ST=Jakarta/L=Jakarta/O=MBankingCore/OU=Development/CN=localhost
Valid for: 365 days
Key size: 4096 bits RSA
```

### Production Certificates

For production, replace auto-generated certificates with proper SSL certificates:

1. **Let's Encrypt (Free)**:
   ```bash
   # Install certbot
   brew install certbot
   
   # Generate certificate
   certbot certonly --standalone -d yourdomain.com
   
   # Copy certificates
   cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ./certs/server.crt
   cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ./certs/server.key
   ```

2. **Commercial SSL Certificate**:
   - Purchase SSL certificate from CA
   - Download certificate files
   - Place them in the `CERT_DIR` directory

### Manual Certificate Generation

```bash
# Create certificates directory
mkdir -p certs

# Generate private key
openssl genrsa -out certs/server.key 4096

# Generate certificate signing request
openssl req -new -key certs/server.key -out certs/server.csr \
  -subj "/C=ID/ST=Jakarta/L=Jakarta/O=MBankingCore/OU=Development/CN=localhost"

# Generate self-signed certificate
openssl x509 -req -days 365 -in certs/server.csr -signkey certs/server.key -out certs/server.crt
```

## 🚀 Server Modes

### HTTPS Mode (Secure)

```bash
# Set environment
export ENABLE_HTTPS=true

# Start server
./mbankingcore
```

**Server URLs**:
- 🔒 HTTPS API: `https://localhost:8443/api`
- 🔐 Health Check: `https://localhost:8443/health`

### HTTP Mode (Default)

```bash
# Set environment (or omit for default)
export ENABLE_HTTPS=false

# Start server
./mbankingcore
```

**Server URLs**:
- 🌐 HTTP API: `http://localhost:8080/api`
- 📋 Health Check: `http://localhost:8080/health`

## 📱 Flutter Client Configuration

### HTTPS URLs

Update your Flutter app's API configuration:

```dart
class ApiConfig {
  static const String baseUrl = _getBaseUrl();
  
  static String _getBaseUrl() {
    if (Platform.isAndroid) {
      // Android emulator - use 10.0.2.2 for host machine
      return 'https://10.0.2.2:8443/api';
    } else if (Platform.isIOS) {
      // iOS simulator - use localhost
      return 'https://localhost:8443/api';
    } else {
      // Physical device - use actual IP address
      return 'https://192.168.123.206:8443/api';
    }
  }
}
```

### SSL Certificate Handling

For development with self-signed certificates:

```dart
// lib/services/http_client.dart
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:http/io_client.dart';

class ApiClient {
  static http.Client _createHttpClient() {
    if (kDebugMode) {
      // Accept self-signed certificates in development
      HttpClient httpClient = HttpClient()
        ..badCertificateCallback = (X509Certificate cert, String host, int port) => true;
      
      return IOClient(httpClient);
    } else {
      // Use secure client in production
      return http.Client();
    }
  }

  static final client = _createHttpClient();
}
```

## 🧪 Testing HTTPS

### Command Line Testing

```bash
# Test HTTPS health endpoint
curl -k https://localhost:8443/health

# Test HTTPS API with verbose output
curl -k -v https://localhost:8443/api/onboardings

# Test with proper certificate validation (production)
curl https://yourdomain.com:8443/health
```

### Browser Testing

1. Open: `https://localhost:8443/health`
2. Accept security warning for self-signed certificate
3. Should see: `{"code":200,"message":"MBankingCore API is running","data":{"status":"ok"}}`

## 🔧 Troubleshooting

### Common Issues

1. **OpenSSL Not Found**:
   ```bash
   # macOS
   brew install openssl
   
   # Ubuntu/Debian
   sudo apt-get install openssl
   ```

2. **Permission Denied**:
   ```bash
   # Fix certificate permissions
   chmod 600 certs/server.key
   chmod 644 certs/server.crt
   ```

3. **Port Already in Use**:
   ```bash
   # Check what's using the port
   lsof -i :8443
   
   # Kill process if needed
   sudo kill -9 <PID>
   ```

4. **Certificate Expired**:
   ```bash
   # Remove old certificates
   rm certs/server.*
   
   # Restart server to generate new ones
   ./mbankingcore
   ```

### Log Messages

**Successful HTTPS startup**:
```
🔒 HTTPS mode enabled
🔐 Generating self-signed SSL certificate for development...
✅ Self-signed SSL certificate generated successfully
🔒 Starting HTTPS server on 0.0.0.0:8443
🔐 Using certificate: ./certs/server.crt
🔑 Using private key: ./certs/server.key
🚀 Starting HTTPS server on 0.0.0.0:8443
```

**Fallback to HTTP**:
```
🔒 HTTPS mode enabled
❌ Failed to generate SSL certificate: ...
🔄 Falling back to HTTP mode...
🌐 HTTP mode enabled
🚀 Starting HTTP server on 0.0.0.0:8080
```

## 🛡️ Security Considerations

### Development
- ✅ Self-signed certificates for local development
- ✅ TLS 1.2+ minimum version
- ✅ Strong cipher suites
- ⚠️ Certificate warnings in browsers (expected)

### Production
- ✅ Valid SSL certificates from trusted CA
- ✅ Regular certificate renewal
- ✅ Strong TLS configuration
- ✅ HSTS headers (recommended)
- ✅ Certificate pinning in mobile apps (recommended)

## 📄 Certificate Information

Generated certificates include:
- **Algorithm**: RSA 4096-bit
- **Validity**: 365 days
- **Extensions**: Basic constraints, key usage
- **Subject Alternative Names**: localhost, 127.0.0.1

For production use, ensure certificates include all domain names and IP addresses your API will serve.
