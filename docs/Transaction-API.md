# Transaction API Documentation

## Overview
API untuk mengelola transaksi balance (top-up dan withdraw) dengan riwayat transaksi.

## Authentication
Semua endpoint memerlukan Bearer token di header `Authorization`.

## Endpoints

### 1. Top-up Balance
**POST** `/api/transactions/topup`

Menambahkan balance ke akun user.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "amount": 100000
}
```

**Response Success (200):**
```json
{
  "code": 200,
  "message": "Top-up successful",
  "data": {
    "transaction_id": 1,
    "amount": 100000,
    "balance_before": 50000,
    "balance_after": 150000,
    "transaction_at": "2024-01-01T10:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "code": 400,
  "message": "Key: 'TopupRequest.Amount' Error:Field validation for 'Amount' failed on the 'min' tag"
}
```

---

### 2. Withdraw Balance
**POST** `/api/transactions/withdraw`

Mengurangi balance dari akun user.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "amount": 50000
}
```

**Response Success (200):**
```json
{
  "code": 200,
  "message": "Withdrawal successful",
  "data": {
    "transaction_id": 2,
    "amount": 50000,
    "balance_before": 150000,
    "balance_after": 100000,
    "transaction_at": "2024-01-01T10:05:00Z"
  }
}
```

**Response Error (400) - Insufficient Balance:**
```json
{
  "code": 400,
  "message": "Insufficient balance"
}
```

---

### 3. Transfer Balance
**POST** `/api/transactions/transfer`

Transfer balance ke user lain menggunakan nomor rekening.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "to_account_number": "1234567890",
  "amount": 75000,
  "description": "Transfer untuk bayar kos"
}
```

**Response Success (200):**
```json
{
  "code": 200,
  "message": "Transfer successful",
  "data": {
    "transaction_id": 3,
    "to_account_number": "1234567890",
    "to_account_name": "Jane Doe",
    "amount": 75000,
    "sender_balance_before": 100000,
    "sender_balance_after": 25000,
    "description": "Transfer untuk bayar kos",
    "transaction_at": "2024-01-01T10:10:00Z"
  }
}
```

**Response Error (400) - Insufficient Balance:**
```json
{
  "code": 400,
  "message": "Insufficient balance"
}
```

**Response Error (404) - Account Not Found:**
```json
{
  "code": 404,
  "message": "Recipient account number not found or inactive"
}
```

**Response Error (400) - Self Transfer:**
```json
{
  "code": 400,
  "message": "Cannot transfer to your own account"
}
```

---

### 4. Get User Transaction History

**GET** `/api/transactions/history`

Mendapatkan riwayat transaksi user.

**Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

- `page` (optional): Halaman yang diminta (default: 1)
- `limit` (optional): Jumlah data per halaman (default: 10)

**Response Success (200):**

```json
{
  "code": 200,
  "message": "User transactions retrieved successfully",
  "data": {
    "transactions": [
      {
        "id": 1,
        "type": "topup",
        "amount": 100000,
        "balance_before": 0,
        "balance_after": 100000,
        "description": "Top up via ATM",
        "status": "completed",
        "created_at": "2024-01-01T10:00:00Z"
      },
      {
        "id": 2,
        "type": "withdraw",
        "amount": 25000,
        "balance_before": 100000,
        "balance_after": 75000,
        "description": "Withdraw untuk belanja",
        "status": "completed",
        "created_at": "2024-01-01T10:05:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

---

### 5. Get All Transactions (Admin Only)

---

---

### 5. Get All Transactions (Admin Only)

**GET** `/api/admin/transactions`

Mendapatkan semua transaksi (admin only).

**Headers:**

```
Authorization: Bearer <admin_token>
```

**Query Parameters:**

- `page` (optional): Halaman yang diminta (default: 1)
- `limit` (optional): Jumlah data per halaman (default: 10)
- `user_id` (optional): Filter berdasarkan user ID

**Response Success (200):**

```json
{
  "code": 200,
  "message": "All transactions retrieved successfully",
  "data": {
    "transactions": [
      {
        "id": 1,
        "user_id": 1,
        "user_name": "John Doe",
        "type": "topup",
        "amount": 100000,
        "balance_before": 0,
        "balance_after": 100000,
        "description": "Top up via ATM",
        "status": "completed",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

## Transaction Types

- `topup`: Penambahan balance
- `withdraw`: Pengurangan balance  
- `transfer_out`: Transfer keluar
- `transfer_in`: Transfer masuk

## Transaction Status

- `completed`: Transaksi berhasil
- `failed`: Transaksi gagal
- `pending`: Transaksi menunggu

## Business Rules

### Top-up

- Amount minimal 1
- Tidak ada limit maksimal
- Balance akan bertambah sesuai amount

### Withdraw

- Amount minimal 1
- Tidak boleh melebihi balance yang tersedia
- Balance akan berkurang sesuai amount

### Transfer

- Amount minimal 1
- Tidak boleh melebihi balance yang tersedia
- Tidak boleh transfer ke rekening sendiri
- Recipient account harus aktif dan valid
- Akan membuat 2 record transaksi: transfer_out untuk sender, transfer_in untuk receiver

## Transaction Safety

- Menggunakan database transaction dengan row-level locking
- Validasi balance dilakukan sebelum update
- Rollback otomatis jika terjadi error

## Error Codes

- `400`: Bad Request (validation error, insufficient balance)
- `401`: Unauthorized
- `403`: Forbidden (admin only)
- `404`: Not Found (user/account not found)
- `500`: Internal Server Error

## Notes

- Semua amount menggunakan format integer (int64) dalam satuan terkecil mata uang
- Balance selalu dalam format integer (contoh: 100000 = Rp 100.000)
- Semua response menggunakan timezone UTC
- Transaction history diurutkan berdasarkan created_at descending (terbaru ke terlama)
