# Dashboard API

## Overview
The Dashboard API provides comprehensive statistics for the banking application, designed for admin users to monitor system activity and performance.

## Endpoint

### GET /api/admin/dashboard

Get dashboard statistics including user counts, admin counts, and transaction summaries.

**Authentication Required:** Yes (Admin Token)

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_access_token>
```

**Response:**

```json
{
    "code": 200,
    "message": "Dashboard data retrieved successfully",
    "data": {
        "total_users": 6067,
        "total_admins": 50,
        "total_transactions": {
            "today": 159,
            "this_month": 5133,
            "this_year": 35375
        },
        "topup_transactions": {
            "today": 43,
            "this_month": 1258,
            "this_year": 8805
        },
        "withdraw_transactions": {
            "today": 50,
            "this_month": 1300,
            "this_year": 8840
        },
        "transfer_transactions": {
            "today": 66,
            "this_month": 2575,
            "this_year": 17730
        },
        "total_transactions_amount": {
            "today": 15000000,
            "this_month": 513300000,
            "this_year": 3537500000
        },
        "total_topup_amount": {
            "today": 4300000,
            "this_month": 125800000,
            "this_year": 880500000
        },
        "total_transfer_amount": {
            "today": 6600000,
            "this_month": 257500000,
            "this_year": 1773000000
        }
    }
}
```

## Data Fields

### Dashboard Statistics
- **total_users** (integer): Total number of registered users in the system
- **total_admins** (integer): Total number of admin users in the system

### Transaction Periods
Each transaction type includes three time period counts:
- **today** (integer): Transactions created today (from 00:00:00 to 23:59:59)
- **this_month** (integer): Transactions created this month (from 1st to last day of current month)
- **this_year** (integer): Transactions created this year (from January 1st to December 31st)

### Transaction Types
- **total_transactions**: All transactions combined (topup, withdraw, transfer_in, transfer_out, etc.)
- **topup_transactions**: Only topup transactions
- **withdraw_transactions**: Only withdraw transactions  
- **transfer_transactions**: Both transfer_in and transfer_out transactions combined

### Transaction Amount Types  
- **total_transactions_amount**: Total amount of all transactions combined (topup, withdraw, transfer_in, transfer_out, etc.)
- **total_topup_amount**: Total amount of topup transactions only
- **total_transfer_amount**: Total amount of transfer_out transactions only (to avoid double counting)

**Note**: All amounts are in Indonesian Rupiah (IDR) and represented as integers.

## Authentication

This endpoint requires admin authentication. Use the admin login endpoint first:

```bash
curl -X POST "http://localhost:8080/api/admin/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@mbankingcore.com",
    "password": "Admin123?"
  }'
```

Then use the returned `access_token` in the Authorization header.

## Error Responses

### 301 - Invalid Token
```json
{
    "code": 301,
    "message": "Invalid or expired token",
    "data": null
}
```

### 251 - Database Error
```json
{
    "code": 251,
    "message": "Failed to retrieve dashboard data",
    "data": null
}
```

## Example Usage

### Using cURL
```bash
# 1. Login as admin
TOKEN=$(curl -X POST "http://localhost:8080/api/admin/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@mbankingcore.com",
    "password": "Admin123?"
  }' | jq -r '.data.access_token')

# 2. Get dashboard statistics
curl -X GET "http://localhost:8080/api/admin/dashboard" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"
```

### Using JavaScript/Fetch
```javascript
// Login and get dashboard data
async function getDashboardData() {
  // Login
  const loginResponse = await fetch('/api/admin/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      email: 'admin@mbankingcore.com',
      password: 'Admin123?'
    })
  });
  
  const loginData = await loginResponse.json();
  const token = loginData.data.access_token;
  
  // Get dashboard
  const dashboardResponse = await fetch('/api/admin/dashboard', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
  
  const dashboardData = await dashboardResponse.json();
  console.log(dashboardData.data);
}
```

## Time Zone Considerations

All time calculations are performed using the server's local timezone. The time periods are calculated as follows:

- **Today**: From 00:00:00 to 23:59:59 of the current date
- **This Month**: From the 1st day 00:00:00 to the last day 23:59:59 of the current month
- **This Year**: From January 1st 00:00:00 to December 31st 23:59:59 of the current year

## Notes

- The endpoint uses efficient database queries with proper indexing for optimal performance
- Transfer transactions include both `transfer_in` and `transfer_out` types to provide complete transfer activity overview
- All counts exclude soft-deleted records (records with `deleted_at` not null)
- The API is designed to be called frequently for real-time dashboard updates
