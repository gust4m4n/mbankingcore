# MBankingCore - Sistem Migrasi

Proyek ini menggunakan sistem migrasi terintegrasi yang secara otomatis mengatur database dan data awal ketika aplikasi dimulai.

## Cara Kerja

### Migrasi Otomatis saat Startup

Ketika Anda menjalankan aplikasi utama (`go run main.go`), sistem secara otomatis:

1. **Terhubung ke database**
2. **Menjalankan semua migrasi** (pembuatan tabel, update schema)
3. **Mengisi data awal** (konfigurasi default, konten sampel)
4. **Memulai API server**

### Tool Migrasi Manual

Untuk setup database tanpa memulai server, gunakan tool migrasi:

```bash
# Build tool migrasi
go build -o migrate ./cmd/migrate

# Jalankan migrasi saja
./migrate
```

## Komponen Migrasi

### 1. Auto-Migration

- Membuat semua tabel database berdasarkan struktur model
- Update tabel yang ada ketika model berubah
- Menangani: Users, DeviceSessions, Articles, Onboarding, Photos, Config

### 2. Migrasi Kustom

- **User Roles**: Memastikan semua user memiliki assignment role yang tepat
- **Konsistensi Data**: Memperbaiki masalah integritas data
- **Update Schema**: Menangani perubahan schema yang kompleks

### 3. Seeding Data Awal

#### Data Konfigurasi

- `app_version`: Info versi aplikasi (default: "1.0.0")
- `tnc`: Konten default Syarat dan Ketentuan
- `privacy-policy`: Konten default Kebijakan Privasi  
- `maintenance_mode`: Flag maintenance (default: "false")
- `max_upload_size`: Batas ukuran upload file (default: "10485760" = 10MB)

#### Konten Onboarding

- Slide welcome: "Welcome to MBankingCore"
- Fitur authentication: "Secure Authentication"
- Info manajemen API: "Easy API Management"
- Panduan getting started: "Get Started"

## Struktur Database

### Tabel yang Dibuat Otomatis

#### 1. Users
```sql
users (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    phone VARCHAR,
    provider VARCHAR DEFAULT 'email',
    role VARCHAR DEFAULT 'user',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

#### 2. Device Sessions
```sql
device_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    device_id VARCHAR NOT NULL,
    device_type VARCHAR,
    device_name VARCHAR,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    last_activity TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

#### 3. Articles
```sql
articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER REFERENCES users(id),
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

#### 4. Onboarding
```sql
onboarding (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    image VARCHAR,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

#### 5. Photos
```sql
photos (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    filename VARCHAR NOT NULL,
    original_name VARCHAR,
    file_size INTEGER,
    mime_type VARCHAR,
    path VARCHAR NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

#### 6. Config
```sql
config (
    id SERIAL PRIMARY KEY,
    key VARCHAR UNIQUE NOT NULL,
    value TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

## Update Terbaru (Juli 2025)

### Perubahan Project
- **Nama Project**: Berubah dari `mbxbackend` menjadi `mbankingcore`
- **Branding**: Menggunakan `MBankingCore` untuk dokumentasi dan label
- **Dokumentasi**: Semua file .md dipindah ke root directory
- **Lokalisasi**: Semua dokumentasi ditranslasi ke bahasa Indonesia

### Fitur Baru
- **Multi-Platform Authentication**: Support untuk berbagai device dengan session management
- **Device Session Management**: Tracking device individual dengan logout selektif
- **User Management**: Role-based access control (admin/user)
- **Article Management**: CRUD operations untuk artikel dengan author tracking
- **Photo Gallery**: Sistem upload dan manajemen foto dengan metadata
- **Configuration Management**: Dynamic app configuration melalui database

### Perubahan Database
- **Removed user_agent field**: Field `user_agent` telah dihapus dari model `DeviceSession` dan tabel `device_sessions` untuk menyederhanakan struktur data

### Security Enhancements
- **Double-Layer Password Protection**: SHA256 (client) + bcrypt (server)
- **JWT Token Strategy**: Access token (15 menit) + Refresh token (7 hari)
- **Device-Specific Tokens**: Token terikat dengan device_id untuk keamanan extra
- **Auto Session Invalidation**: Password change membatalkan semua sessions

## Konfigurasi Database

Sistem menggunakan environment variables berikut:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=mbankingcore
DB_SSLMODE=disable
```

## Keamanan Migrasi

### Setup Pertama Kali
- ✅ **Aman**: Membuat semua tabel dan data awal
- ✅ **Idempotent**: Dapat dijalankan berkali-kali dengan aman
- ✅ **Non-destructive**: Tidak pernah menghapus data yang ada

### Update
- ✅ **Otomatis**: Berjalan di setiap start aplikasi
- ✅ **Pintar**: Hanya membuat data yang hilang, melewati yang sudah ada
- ✅ **Logged**: Semua aksi dicatat untuk debugging

## Struktur File

```
config/
├── database.go    # Koneksi database dan trigger migrasi
└── migrations.go  # Logic migrasi dan seeding data

cmd/
└── migrate/
    └── main.go    # Tool migrasi standalone

models/            # Model database (auto-migrated)
├── user.go
├── onboarding.go
├── config.go
└── ...
```

## Contoh Penggunaan

### Setup Development
```bash
# 1. Setup environment
cp .env.example .env
# Edit .env dengan kredensial database Anda

# 2. Start aplikasi (migrasi berjalan otomatis)
go run main.go
```

### Deployment Production
```bash
# 1. Jalankan migrasi terlebih dahulu
go build -o migrate ./cmd/migrate
./migrate

# 2. Start aplikasi
go build -o mbankingcore .
./mbankingcore
```

### Reset Database (Development)
```bash
# Drop database, buat ulang, dan jalankan migrasi
dropdb mbankingcore
createdb mbankingcore
go run main.go
```

## Troubleshooting

### Migrasi Gagal
- Cek kredensial database di `.env`
- Pastikan PostgreSQL sedang berjalan
- Verifikasi database ada dan dapat diakses

### Data Duplikat
- Sistem bersifat idempotent - aman untuk dijalankan ulang
- Data yang ada dipertahankan
- Hanya data yang hilang yang dibuat

### Masalah Schema
- GORM AutoMigrate menangani sebagian besar perubahan schema
- Untuk perubahan kompleks, tambahkan fungsi migrasi kustom
- Backup database sebelum update besar

## Menambah Migrasi Baru

Untuk menambah logic migrasi baru:

1. **Edit `config/migrations.go`**
2. **Tambahkan function ke `runCustomMigrations()`**
3. **Test dengan migration tool**: `./migrate`
4. **Deploy**: Migrasi berjalan otomatis saat startup

Contoh:
```go
func runCustomMigrations() error {
    // Migrasi yang ada...
    
    // Tambah migrasi baru
    if err := migrateNewFeature(); err != nil {
        return err
    }
    
    return nil
}

func migrateNewFeature() error {
    // Logic migrasi Anda di sini
    return nil
}
```

## Keuntungan

✅ **Zero-Config**: Bekerja langsung tanpa konfigurasi  
✅ **Developer Friendly**: Tidak perlu script SQL manual  
✅ **Production Ready**: Aman untuk deployment production  
✅ **Maintainable**: Semua logic migrasi di satu tempat  
✅ **Fleksibel**: Mudah menambah migrasi baru  
✅ **Logged**: Visibilitas lengkap proses migrasi

## Monitoring & Log Migrasi

### Log Output Contoh
```
Starting database migrations...
✅ Auto-migration completed successfully
Running custom migrations...
Migrating user roles...
✅ All users already have proper roles
✅ Custom migrations completed
Seeding initial data...
Seeding initial configurations...
✅ Config already exists: app_version
✅ Config already exists: tnc
✅ Config already exists: privacy-policy
✅ Config already exists: maintenance_mode
✅ Config already exists: max_upload_size
Seeding initial onboarding content...
✅ Onboarding content already exists
✅ Initial data seeding completed
🚀 All migrations and initial setup completed successfully!
```

### Performance Tips

- **Database Connection**: Gunakan connection pooling untuk performa optimal
- **Migration Time**: Migrasi biasanya selesai dalam < 5 detik untuk database kosong
- **Memory Usage**: Sistem menggunakan memory minimal untuk migrasi
- **Concurrent Access**: Aman untuk multiple instances (menggunakan database locks)

### Best Practices

1. **Backup Database**: Selalu backup sebelum menjalankan migrasi di production
2. **Test Environment**: Test migrasi di development environment terlebih dahulu
3. **Monitor Logs**: Pantau log migrasi untuk memastikan semua berjalan dengan baik
4. **Version Control**: Track perubahan migration logic di version control
5. **Rollback Plan**: Siapkan rencana rollback untuk perubahan schema besar

---

**💡 Tips**: Untuk development, Anda bisa menjalankan `go run main.go` dan migrasi akan berjalan otomatis. Untuk production, gunakan tool migrasi terpisah terlebih dahulu untuk memastikan database siap sebelum menjalankan aplikasi.
