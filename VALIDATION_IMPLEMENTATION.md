# Banking Login Validation System

## Implemented Validation Rules

The `/api/login` endpoint now includes comprehensive validation as requested:

### 1. Field Length Validation
- ✅ **Name**: Minimum 8 characters
- ✅ **Account Number**: Minimum 8 characters  
- ✅ **Phone**: Minimum 8 characters
- ✅ **Mother Name**: Minimum 8 characters
- ✅ **PIN ATM**: Exactly 6 digits (numeric only)

### 2. Phone Number Registration Logic
- ✅ **Existing Users**: Validates account data and PIN before sending OTP
- ✅ **New Users**: Auto-registers after OTP verification
- ✅ **Security**: All validations checked before OTP generation

## API Request/Response Examples

### Valid Login Request
```json
POST /api/login
{
  "name": "John Doe Smith",           // ✅ 8+ characters
  "account_number": "12345678",       // ✅ 8+ characters
  "phone": "6281234567890",          // ✅ 8+ characters
  "mother_name": "Jane Doe Williams", // ✅ 8+ characters
  "pin_atm": "123456",               // ✅ Exactly 6 digits
  "device_info": {
    "device_type": "mobile",
    "device_id": "device123",
    "device_name": "iPhone 13"
  }
}
```

### Success Response (Existing User)
```json
{
  "code": 200,
  "message": "OTP sent to your registered phone number",
  "data": {
    "login_token": "a1b2c3d4e5f6...",
    "expires_in": 300,
    "is_new_user": false
  }
}
```

### Success Response (New User)
```json
{
  "code": 200,
  "message": "Phone number will be registered. OTP sent for verification",
  "data": {
    "login_token": "a1b2c3d4e5f6...",
    "expires_in": 300,
    "is_new_user": true
  }
}
```

## Validation Error Examples

### Name Too Short
```json
{
  "code": 400,
  "message": "Name must be at least 8 characters long"
}
```

### Account Number Too Short
```json
{
  "code": 400,
  "message": "Account number must be at least 8 characters long"
}
```

### Phone Too Short
```json
{
  "code": 400,
  "message": "Phone number must be at least 8 characters long"
}
```

### Mother Name Too Short
```json
{
  "code": 400,
  "message": "Mother name must be at least 8 characters long"
}
```

### Invalid PIN Length
```json
{
  "code": 400,
  "message": "PIN ATM must be exactly 6 digits"
}
```

### Non-Numeric PIN
```json
{
  "code": 400,
  "message": "PIN ATM must contain only numeric digits"
}
```

### Account Mismatch (Existing User)
```json
{
  "code": 401,
  "message": "Account information does not match our records"
}
```

## Flow Diagram

### For Existing Users:
1. **Validation** → Name(8+), Account(8+), Phone(8+), Mother(8+), PIN(6 digits)
2. **Database Check** → Phone exists in database
3. **Data Verification** → Name, Account, Mother Name match records
4. **PIN Verification** → PIN matches stored hash
5. **OTP Generation** → Send OTP to registered phone
6. **Response** → login_token with is_new_user: false

### For New Users:
1. **Validation** → Name(8+), Account(8+), Phone(8+), Mother(8+), PIN(6 digits)
2. **Database Check** → Phone not in database
3. **OTP Generation** → Send OTP for registration
4. **Response** → login_token with is_new_user: true
5. **After OTP Verification** → Auto-register user with hashed PIN

## Security Features

- ✅ **Input Validation**: All fields validated before processing
- ✅ **PIN Security**: 6-digit numeric validation with secure hashing
- ✅ **Data Integrity**: Account info verified for existing users
- ✅ **Auto Registration**: Seamless new user onboarding
- ✅ **Secure Tokens**: login_token system for OTP verification
- ✅ **Session Management**: Device-specific session handling

## Test Cases

### Valid Inputs:
- Name: "John Doe Smith" ✅
- Account: "87654321" ✅  
- Phone: "6281234567890" ✅
- Mother: "Jane Doe Williams" ✅
- PIN: "123456" ✅

### Invalid Inputs:
- Name: "John" ❌ (too short)
- Account: "1234" ❌ (too short)
- Phone: "628123" ❌ (too short)
- Mother: "Jane" ❌ (too short)
- PIN: "12345" ❌ (not 6 digits)
- PIN: "12345a" ❌ (contains non-numeric)

The validation system ensures data quality and security before any OTP generation or database operations.
