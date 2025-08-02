# 📋 Postman Collection - Example Responses Update

## 🎯 Overview
Koleksi Postman MBankingCore telah diperbarui untuk menyertakan **contoh response sukses** pada endpoint-endpoint utama. Ini memudahkan developer untuk memahami struktur response yang diharapkan dari API.

## ✅ Endpoint yang Sudah Ditambahkan Example Response

### 1. **Health Check** (`GET /health`)
- ✅ **Success Response (200)**: Server sehat, database terhubung
- ❌ **Error Response (503)**: Database tidak terhubung

```json
{
  "api_status": "healthy",
  "code": 200,
  "current_time": "2025-08-02 08:46:19 WIB",
  "database_status": "healthy",
  "message": "MBankingCore API Health Check",
  "server_started_at": "2025-08-02 08:40:01 WIB",
  "total_uptime_seconds": 378,
  "uptime": {
    "days": 0, "hours": 0, "minutes": 6, "seconds": 18
  },
  "uptime_string": "0y 0mo 0w 0d 0h 6m 18s"
}
```

### 2. **Banking Login** (`POST /api/login`)
- ✅ **Success Response (200)**: Login berhasil, OTP dikirim
- ❌ **Error Response (400)**: PIN ATM tidak valid

```json
{
  "code": 200,
  "message": "Phone number will be registered. OTP sent for verification",
  "data": {
    "expires_in": 300,
    "is_new_user": true,
    "login_token": "47898ef5054fcaffd7ad4afc77899868b1800287e1169014ed7472fc41cb0315"
  }
}
```

### 3. **Banking Login Verification** (`POST /api/login/verify`)
- ✅ **Success Response (200)**: Verifikasi berhasil, token akses diberikan

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "def50200684c7d6ad2e9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "session_id": "session_123456789",
    "user": {
      "id": 1,
      "name": "John Doe Testing",
      "phone": "+62812345678901",
      "account_number": "1234567890123456",
      "balance": 0,
      "status": "active"
    }
  }
}
```

### 4. **User Profile** (`GET /api/profile`)
- ✅ **Success Response (200)**: Data profil user

```json
{
  "code": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe Testing",
    "phone": "+62812345678901",
    "mother_name": "Jane Doe Mother",
    "account_number": "1234567890123456",
    "balance": 250000,
    "status": "active",
    "avatar": null,
    "created_at": "2025-08-02T01:40:01.000Z",
    "updated_at": "2025-08-02T01:45:15.000Z"
  }
}
```

### 5. **Admin Dashboard** (`GET /api/admin/dashboard`)
- ✅ **Success Response (200)**: Statistik dashboard komprehensif

```json
{
  "code": 200,
  "message": "Dashboard data retrieved successfully",
  "data": {
    "total_users": 12,
    "total_admins": 3,
    "total_transactions": {
      "today": 5, "today_amount": 750000,
      "this_week": 23, "this_week_amount": 3250000,
      "this_month": 89, "this_month_amount": 12500000,
      "all_time": 1234, "all_time_amount": 123400000
    },
    "performance": {
      "last_7_days": [
        { "period": "Jul 27", "count": 3, "amount": 450000 }
      ],
      "last_30_days": [
        { "period": "Jul 04", "count": 2, "amount": 300000 }
      ]
    }
  }
}
```

### 6. **Topup Balance** (`POST /api/transactions/topup`)
- ✅ **Success Response (200)**: Topup berhasil

```json
{
  "code": 200,
  "message": "Topup successful",
  "data": {
    "transaction": {
      "id": 123,
      "user_id": 1,
      "type": "topup",
      "amount": 100000,
      "balance_before": 250000,
      "balance_after": 350000,
      "description": "Top up saldo untuk kebutuhan sehari-hari",
      "status": "completed",
      "created_at": "2025-08-02T01:50:15.000Z"
    },
    "user": {
      "id": 1,
      "name": "John Doe Testing",
      "account_number": "1234567890123456",
      "balance": 350000
    }
  }
}
```

## 🔄 Update Details

### File yang Diubah:
- `postman/MBankingCore-API.postman_collection.json`

### Perubahan:
1. **Added Example Responses**: Menambahkan contoh response sukses dan error
2. **Updated Collection Name**: Menambahkan "with Examples" di nama koleksi
3. **Updated Description**: Menyebutkan fitur example responses
4. **Real Response Data**: Menggunakan struktur response yang sebenarnya dari API

### Benefit:
- ✅ **Developer Experience**: Lebih mudah memahami struktur response
- ✅ **Documentation**: Contoh real response sebagai dokumentasi
- ✅ **Testing**: Dapat membandingkan actual vs expected response
- ✅ **Debugging**: Membantu troubleshooting saat response tidak sesuai

## 🚀 Cara Menggunakan

1. **Import Collection**: Import file `MBankingCore-API.postman_collection.json`
2. **Import Environment**: Import file `MBankingCore-API.postman_environment.json`
3. **View Examples**: Klik pada request → Tab "Examples" untuk melihat contoh response
4. **Test API**: Jalankan request dan bandingkan dengan contoh response

## 📝 Notes

- Example responses diambil dari actual API response
- Beberapa data sensitif (seperti token) telah di-sanitize
- Error responses menunjukkan berbagai skenario kegagalan
- Response structure mengikuti standar API MBankingCore

---

**Happy Testing!** 🎉
