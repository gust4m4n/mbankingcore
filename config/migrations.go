package config

import (
	"log"
	"mbankingcore/models"

	"golang.org/x/crypto/bcrypt"
)

// RunMigrations handles database setup for new project
func RunMigrations() error {
	log.Println("Setting up database for new project...")

	// Auto-migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.BankAccount{},
		&models.DeviceSession{},
		&models.OTPSession{},
		&models.Article{},
		&models.Onboarding{},
		&models.Photo{},
		&models.Config{},
		&models.Transaction{},
		&models.AuditLog{},
		&models.LoginAudit{},
	)
	if err != nil {
		log.Printf("Failed to auto-migrate models: %v", err)
		return err
	}
	log.Println("✅ Database tables created successfully")

	// Seed initial data
	if err := seedInitialData(); err != nil {
		log.Printf("Failed to seed initial data: %v", err)
		return err
	}

	log.Println("🚀 Database setup completed successfully!")
	return nil
}

// seedInitialData creates essential initial data for new project
func seedInitialData() error {
	log.Println("Seeding initial data...")

	// Seed initial admin users
	if err := seedInitialAdmins(); err != nil {
		return err
	}

	// Seed initial configuration values
	if err := seedInitialConfigs(); err != nil {
		return err
	}

	// Seed initial onboarding content
	if err := seedInitialOnboarding(); err != nil {
		return err
	}

	log.Println("✅ Initial data seeding completed")
	return nil
}

// seedInitialAdmins creates default admin users
func seedInitialAdmins() error {
	log.Println("Seeding initial admin users...")

	// Check if admin users already exist
	var count int64
	DB.Model(&models.Admin{}).Count(&count)

	if count > 0 {
		log.Println("✅ Admin users already exist")
		return nil
	}

	// Create admin user 1: admin@mbankingcore.com
	hashedPassword1, err := bcrypt.GenerateFromPassword([]byte("Admin123?"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return err
	}

	admin1 := models.Admin{
		Name:     "Admin User",
		Email:    "admin@mbankingcore.com",
		Password: string(hashedPassword1),
		Role:     "admin",
		Status:   1, // active
	}

	if err := DB.Create(&admin1).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return err
	}

	// Create super admin user 2: super@mbankingcore.com
	hashedPassword2, err := bcrypt.GenerateFromPassword([]byte("Super123?"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash super admin password: %v", err)
		return err
	}

	admin2 := models.Admin{
		Name:     "Super Admin",
		Email:    "super@mbankingcore.com",
		Password: string(hashedPassword2),
		Role:     "super",
		Status:   1, // active
	}

	if err := DB.Create(&admin2).Error; err != nil {
		log.Printf("Failed to create super admin: %v", err)
		return err
	}

	log.Println("✅ Created essential admin users:")
	log.Println("   - admin@mbankingcore.com (role: admin)")
	log.Println("   - super@mbankingcore.com (role: super)")
	return nil
}

// seedInitialConfigs creates default configuration values
func seedInitialConfigs() error {
	log.Println("Seeding initial configuration values...")

	// Check if configs already exist
	var count int64
	DB.Model(&models.Config{}).Count(&count)

	if count > 0 {
		log.Println("✅ Configuration values already exist")
		return nil
	}

	initialConfigs := []models.Config{
		{Key: "app_name", Value: "MBankingCore"},
		{Key: "app_version", Value: "0.9"},
		{Key: "terms_conditions", Value: getTermsConditionsContent()},
		{Key: "privacy_policy", Value: getPrivacyPolicyContent()},
		{Key: "contact_email", Value: "support@mbankingcore.com"},
		{Key: "contact_phone", Value: "+62-21-12345678"},
		{Key: "maintenance_mode", Value: "false"},
		{Key: "max_sessions_per_user", Value: "5"},
	}

	for _, config := range initialConfigs {
		if err := DB.Create(&config).Error; err != nil {
			log.Printf("Failed to create config %s: %v", config.Key, err)
			return err
		}
	}

	log.Printf("✅ Created %d initial configuration values", len(initialConfigs))
	return nil
}

// seedInitialOnboarding creates default onboarding slides
func seedInitialOnboarding() error {
	log.Println("Seeding initial onboarding content...")

	// Check if onboarding content already exists
	var count int64
	DB.Model(&models.Onboarding{}).Count(&count)

	if count > 0 {
		log.Println("✅ Onboarding content already exists")
		return nil
	}

	initialOnboarding := []models.Onboarding{
		{
			Title:       "Welcome to MBankingCore",
			Description: "Your secure and reliable mobile banking solution",
			Image:       "https://example.com/welcome.png",
			IsActive:    true,
		},
		{
			Title:       "Secure Banking",
			Description: "Bank safely with advanced security features",
			Image:       "https://example.com/security.png",
			IsActive:    true,
		},
		{
			Title:       "Easy Transactions",
			Description: "Send money and pay bills with just a few taps",
			Image:       "https://example.com/transactions.png",
			IsActive:    true,
		},
		{
			Title:       "24/7 Support",
			Description: "Get help whenever you need it",
			Image:       "https://example.com/support.png",
			IsActive:    true,
		},
	}

	for _, onboarding := range initialOnboarding {
		if err := DB.Create(&onboarding).Error; err != nil {
			log.Printf("Failed to create onboarding slide: %v", err)
			return err
		}
	}

	log.Printf("✅ CreatseedInitialUsersed %d initial onboarding slides", len(initialOnboarding))
	return nil
}

// getTermsConditionsContent returns comprehensive Terms & Conditions content for Indonesian banking app
func getTermsConditionsContent() string {
	return `<h1>SYARAT DAN KETENTUAN PENGGUNAAN MBANKINGCORE</h1>

<p><b>Terakhir diperbarui: 31 Juli 2025</b></p>

<h2>1. PENERIMAAN SYARAT DAN KETENTUAN</h2>
<p>Dengan mengunduh, menginstal, mengakses atau menggunakan aplikasi MBankingCore ("<b>Aplikasi</b>"), Anda menyetujui untuk terikat oleh Syarat dan Ketentuan ini ("<b>S&K</b>"). Jika Anda tidak menyetujui S&K ini, harap tidak menggunakan Aplikasi.</p>

<h2>2. DEFINISI</h2>
<ul>
<li><b>"Aplikasi"</b> adalah aplikasi mobile banking MBankingCore</li>
<li><b>"Layanan"</b> adalah semua fitur dan layanan yang disediakan melalui Aplikasi</li>
<li><b>"Pengguna"</b> adalah individu yang telah terdaftar dan menggunakan Layanan</li>
<li><b>"Kami"</b> adalah penyedia Aplikasi MBankingCore</li>
<li><b>"Rekening"</b> adalah rekening bank yang terhubung dengan Aplikasi</li>
</ul>

<h2>3. PERSYARATAN PENGGUNAAN</h2>
<h3>3.1 Kelayakan</h3>
<ul>
<li>Anda harus berusia minimal 17 tahun</li>
<li>Memiliki rekening bank yang valid di Indonesia</li>
<li>Memiliki nomor telepon aktif yang terdaftar di bank</li>
<li>Menyediakan informasi yang akurat dan lengkap</li>
</ul>

<h3>3.2 Registrasi Akun</h3>
<ul>
<li>Satu nomor telepon hanya dapat digunakan untuk satu akun</li>
<li>Anda bertanggung jawab menjaga kerahasiaan PIN dan password</li>
<li>Segera laporkan jika terjadi penyalahgunaan akun</li>
</ul>

<h2>4. LAYANAN YANG TERSEDIA</h2>
<h3>4.1 Informasi Rekening</h3>
<ul>
<li>Cek saldo</li>
<li>Riwayat transaksi</li>
<li>Informasi rekening</li>
</ul>

<h3>4.2 Transfer Dana</h3>
<ul>
<li>Transfer antar bank</li>
<li>Transfer antar pengguna MBankingCore</li>
<li>Transfer terjadwal</li>
</ul>

<h3>4.3 Pembayaran</h3>
<ul>
<li>Pembayaran tagihan (listrik, air, telepon)</li>
<li>Pembayaran merchant</li>
<li>Top-up e-wallet</li>
</ul>

<h2>5. KEAMANAN DAN PERLINDUNGAN</h2>
<h3>5.1 Kewajiban Pengguna</h3>
<ul>
<li>Menjaga kerahasiaan PIN/password</li>
<li>Menggunakan koneksi internet yang aman</li>
<li>Melakukan logout setelah selesai menggunakan</li>
<li>Tidak berbagi informasi akun dengan pihak lain</li>
</ul>

<h3>5.2 Sistem Keamanan</h3>
<ul>
<li>Enkripsi data end-to-end</li>
<li>Otentikasi dua faktor (2FA)</li>
<li>Monitoring transaksi real-time</li>
<li>Notifikasi setiap transaksi</li>
</ul>

<h2>6. BATASAN DAN LARANGAN</h2>
<h3>6.1 Larangan Penggunaan</h3>
<ul>
<li>Menggunakan Aplikasi untuk kegiatan ilegal</li>
<li>Melakukan transaksi fiktif atau penipuan</li>
<li>Mengganggu sistem atau server Aplikasi</li>
<li>Menyalahgunakan fitur atau layanan</li>
</ul>

<h3>6.2 Batas Transaksi</h3>
<ul>
<li>Transfer harian: <b>Rp 25.000.000</b></li>
<li>Transfer per transaksi: <b>Rp 5.000.000</b></li>
<li>Pembayaran harian: <b>Rp 10.000.000</b></li>
<li>Batas dapat disesuaikan sesuai profil risiko</li>
</ul>

<h2>7. BIAYA DAN TARIF</h2>
<h3>7.1 Biaya Layanan</h3>
<ul>
<li>Transfer antar bank: <b>Rp 6.500</b></li>
<li>Transfer antar pengguna MBankingCore: <b>GRATIS</b></li>
<li>Cek saldo dan mutasi: <b>GRATIS</b></li>
<li>Pembayaran tagihan: <b>Rp 2.500</b></li>
</ul>

<h3>7.2 Perubahan Tarif</h3>
<ul>
<li>Kami berhak mengubah tarif dengan pemberitahuan 30 hari sebelumnya</li>
<li>Perubahan akan diberitahukan melalui Aplikasi atau email</li>
</ul>

<h2>8. PRIVASI DAN PERLINDUNGAN DATA</h2>
<h3>8.1 Pengumpulan Data</h3>
<ul>
<li>Data pribadi dikumpulkan sesuai keperluan layanan</li>
<li>Data transaksi disimpan untuk audit dan compliance</li>
<li>Lokasi perangkat untuk keamanan tambahan</li>
</ul>

<h3>8.2 Penggunaan Data</h3>
<ul>
<li>Memproses transaksi dan layanan</li>
<li>Analisis risiko dan fraud detection</li>
<li>Peningkatan layanan dan fitur</li>
<li>Compliance dengan regulasi</li>
</ul>

<h2>9. TANGGUNG JAWAB DAN GANTI RUGI</h2>
<h3>9.1 Batasan Tanggung Jawab</h3>
<ul>
<li>Tidak bertanggung jawab atas kerugian akibat kelalaian pengguna</li>
<li>Tidak bertanggung jawab atas gangguan jaringan atau sistem bank</li>
<li>Tanggung jawab terbatas pada jumlah transaksi yang bermasalah</li>
</ul>

<h3>9.2 Force Majeure</h3>
<ul>
<li>Tidak bertanggung jawab atas kejadian di luar kendali</li>
<li>Termasuk bencana alam, perang, atau gangguan pemerintah</li>
</ul>

<h2>10. PENANGGUHAN DAN PENGHENTIAN</h2>
<h3>10.1 Penangguhan Akun</h3>
<ul>
<li>Akun dapat ditangguhkan jika melanggar S&K</li>
<li>Penangguhan karena aktivitas mencurigakan</li>
<li>Pemberitahuan akan diberikan jika memungkinkan</li>
</ul>

<h3>10.2 Penghentian Layanan</h3>
<ul>
<li>Pengguna dapat menghentikan layanan kapan saja</li>
<li>Kami dapat menghentikan layanan dengan pemberitahuan 30 hari</li>
<li>Saldo akan dikembalikan sesuai prosedur bank</li>
</ul>

<h2>11. PERUBAHAN SYARAT DAN KETENTUAN</h2>
<ul>
<li>S&K dapat diubah sewaktu-waktu</li>
<li>Perubahan akan diberitahukan melalui Aplikasi</li>
<li>Penggunaan Aplikasi setelah perubahan dianggap sebagai persetujuan</li>
<li>Versi terbaru selalu tersedia di dalam Aplikasi</li>
</ul>

<h2>12. PENYELESAIAN SENGKETA</h2>
<h3>12.1 Hukum yang Berlaku</h3>
<ul>
<li>S&K ini tunduk pada hukum Republik Indonesia</li>
<li>Penyelesaian sengketa melalui pengadilan di Jakarta</li>
</ul>

<h3>12.2 Mediasi</h3>
<ul>
<li>Upaya penyelesaian secara kekeluargaan terlebih dahulu</li>
<li>Mediasi melalui Otoritas Jasa Keuangan (OJK) jika diperlukan</li>
</ul>

<h2>13. KONTAK DAN BANTUAN</h2>
<h3>Customer Service:</h3>
<ul>
<li>Email: <b>support@mbankingcore.com</b></li>
<li>Telepon: <b>1500-888 (24/7)</b></li>
<li>WhatsApp: <b>+62-812-3456-7890</b></li>
<li>Live Chat: Tersedia di dalam Aplikasi</li>
</ul>

<h3>Jam Operasional:</h3>
<ul>
<li>Senin - Jumat: 06.00 - 22.00 WIB</li>
<li>Sabtu - Minggu: 08.00 - 20.00 WIB</li>
<li>Emergency Support: 24/7</li>
</ul>

<hr>
<p><i>Dengan menggunakan Aplikasi MBankingCore, Anda menyatakan telah membaca, memahami, dan menyetujui seluruh Syarat dan Ketentuan di atas.</i></p>`
}

// getPrivacyPolicyContent returns comprehensive Privacy Policy content for Indonesian banking app
func getPrivacyPolicyContent() string {
	return `<h1>KEBIJAKAN PRIVASI MBANKINGCORE</h1>

<p><b>Terakhir diperbarui: 31 Juli 2025</b></p>

<h2>1. PENDAHULUAN</h2>
<p>MBankingCore ("<b>kami</b>", "<b>Aplikasi</b>") berkomitmen untuk melindungi privasi dan keamanan data pribadi Anda. Kebijakan Privasi ini menjelaskan bagaimana kami mengumpulkan, menggunakan, menyimpan, dan melindungi informasi pribadi Anda saat menggunakan layanan kami.</p>

<h2>2. INFORMASI YANG KAMI KUMPULKAN</h2>
<h3>2.1 Data Identitas Pribadi</h3>
<ul>
<li><b>Informasi Dasar:</b> Nama lengkap, tanggal lahir, nomor KTP/NIK</li>
<li><b>Kontak:</b> Nomor telepon, alamat email, alamat rumah</li>
<li><b>Foto:</b> Foto selfie, foto KTP untuk verifikasi identitas</li>
<li><b>Biometrik:</b> Sidik jari, face recognition (jika diaktifkan)</li>
</ul>

<h3>2.2 Data Keuangan</h3>
<ul>
<li><b>Informasi Rekening:</b> Nomor rekening, nama bank, saldo</li>
<li><b>Riwayat Transaksi:</b> Transfer, pembayaran, top-up, withdrawals</li>
<li><b>Pola Penggunaan:</b> Frekuensi transaksi, merchant favorit</li>
<li><b>Data Kredit:</b> Riwayat kredit untuk scoring (jika tersedia)</li>
</ul>

<h3>2.3 Data Teknis</h3>
<ul>
<li><b>Informasi Perangkat:</b> Model, OS, versi aplikasi, device ID</li>
<li><b>Lokasi:</b> GPS location untuk keamanan transaksi</li>
<li><b>Log Aktivitas:</b> Waktu login, IP address, aktivitas dalam aplikasi</li>
<li><b>Cookies:</b> Preferensi pengguna, session management</li>
</ul>

<h3>2.4 Data Komunikasi</h3>
<ul>
<li><b>Customer Service:</b> Rekaman chat, email, telepon</li>
<li><b>Notifikasi:</b> Preferensi push notification, SMS</li>
<li><b>Feedback:</b> Rating, review, saran perbaikan</li>
</ul>

<h2>3. CARA KAMI MENGUMPULKAN DATA</h2>
<h3>3.1 Langsung dari Anda</h3>
<ul>
<li>Saat registrasi akun baru</li>
<li>Pengisian profil dan verifikasi KYC</li>
<li>Melakukan transaksi atau menggunakan fitur</li>
<li>Menghubungi customer service</li>
</ul>

<h3>3.2 Otomatis dari Aplikasi</h3>
<ul>
<li>Log aktivitas dan penggunaan aplikasi</li>
<li>Data lokasi (dengan izin)</li>
<li>Informasi perangkat dan jaringan</li>
<li>Cookies dan teknologi tracking</li>
</ul>

<h3>3.3 Dari Pihak Ketiga</h3>
<ul>
<li><b>Bank Partner:</b> Informasi rekening dan transaksi</li>
<li><b>Credit Bureau:</b> Riwayat kredit dan scoring</li>
<li><b>Anti-Fraud Provider:</b> Verifikasi identitas dan risk assessment</li>
<li><b>Analytics Provider:</b> Data agregat untuk improvement</li>
</ul>

<h2>4. TUJUAN PENGGUNAAN DATA</h2>
<h3>4.1 Penyediaan Layanan</h3>
<ul>
<li><b>Verifikasi Identitas:</b> KYC, AML compliance, fraud prevention</li>
<li><b>Memproses Transaksi:</b> Transfer, pembayaran, top-up</li>
<li><b>Customer Support:</b> Bantuan teknis dan layanan pelanggan</li>
<li><b>Personalisasi:</b> Rekomendasi produk dan fitur yang relevan</li>
</ul>

<h3>4.2 Keamanan dan Kepatuhan</h3>
<ul>
<li><b>Fraud Detection:</b> Monitoring transaksi mencurigakan</li>
<li><b>Risk Management:</b> Penilaian risiko dan credit scoring</li>
<li><b>Regulatory Compliance:</b> Pelaporan ke Bank Indonesia, OJK</li>
<li><b>Audit Trail:</b> Jejak audit untuk investigasi</li>
</ul>

<h3>4.3 Peningkatan Layanan</h3>
<ul>
<li><b>Analytics:</b> Analisis penggunaan untuk improvement</li>
<li><b>A/B Testing:</b> Testing fitur baru untuk user experience</li>
<li><b>Machine Learning:</b> AI untuk fraud detection dan personalization</li>
<li><b>Research:</b> Riset pasar untuk pengembangan produk</li>
</ul>

<h2>5. BERBAGI DATA DENGAN PIHAK KETIGA</h2>
<h3>5.1 Bank dan Financial Institution</h3>
<ul>
<li><b>Bank Partner:</b> Untuk memproses transaksi perbankan</li>
<li><b>Payment Gateway:</b> Pembayaran merchant dan e-commerce</li>
<li><b>E-wallet Provider:</b> Top-up dan transfer antar e-wallet</li>
<li><b>Credit Bureau:</b> Credit scoring dan risk assessment</li>
</ul>

<h3>5.2 Technology Provider</h3>
<ul>
<li><b>Cloud Provider:</b> AWS, Google Cloud untuk data storage</li>
<li><b>Security Provider:</b> Anti-fraud, cybersecurity services</li>
<li><b>Analytics Provider:</b> Google Analytics, Firebase</li>
<li><b>Communication:</b> SMS gateway, email provider, push notification</li>
</ul>

<h2>6. KEAMANAN DATA</h2>
<h3>6.1 Enkripsi</h3>
<ul>
<li><b>Data at Rest:</b> AES-256 encryption untuk data storage</li>
<li><b>Data in Transit:</b> TLS 1.3 untuk komunikasi</li>
<li><b>Database:</b> Field-level encryption untuk data sensitif</li>
<li><b>Backup:</b> Encrypted backup dengan secure key management</li>
</ul>

<h3>6.2 Access Control</h3>
<ul>
<li><b>Role-based Access:</b> Akses terbatas sesuai job function</li>
<li><b>Multi-factor Authentication:</b> 2FA untuk admin access</li>
<li><b>Audit Log:</b> Complete logging untuk semua akses data</li>
<li><b>Privileged Access Management:</b> Secure admin access</li>
</ul>

<h2>7. HAK-HAK ANDA</h2>
<h3>7.1 Hak Akses</h3>
<ul>
<li><b>Data Portability:</b> Ekspor data dalam format standar</li>
<li><b>Data Transparency:</b> Informasi lengkap data yang dimiliki</li>
<li><b>Processing Activity:</b> Detail bagaimana data digunakan</li>
<li><b>Third Party Sharing:</b> List pihak ketiga yang menerima data</li>
</ul>

<h3>7.2 Hak Koreksi</h3>
<ul>
<li><b>Update Profile:</b> Self-service untuk update data pribadi</li>
<li><b>Data Correction:</b> Request koreksi data yang tidak akurat</li>
<li><b>Verification Process:</b> Proses verifikasi untuk data sensitif</li>
<li><b>Notification:</b> Pemberitahuan perubahan ke pihak ketiga</li>
</ul>

<h3>7.3 Hak Penghapusan</h3>
<ul>
<li><b>Account Deletion:</b> Penghapusan akun dan data terkait</li>
<li><b>Selective Deletion:</b> Penghapusan data spesifik</li>
<li><b>Retention Override:</b> Penghapusan sebelum periode retensi</li>
<li><b>Legal Basis:</b> Consideration terhadap kewajiban hukum</li>
</ul>

<h2>8. RETENSI DATA</h2>
<h3>8.1 Periode Penyimpanan</h3>
<ul>
<li><b>Data Transaksi:</b> 10 tahun (sesuai regulasi BI)</li>
<li><b>Data Identitas:</b> Selama akun aktif + 5 tahun</li>
<li><b>Log Komunikasi:</b> 3 tahun untuk audit purpose</li>
<li><b>Analytics Data:</b> 2 tahun dalam bentuk agregat</li>
</ul>

<h3>8.2 Penghapusan Data</h3>
<ul>
<li><b>Account Closure:</b> Data dihapus setelah periode retensi</li>
<li><b>Right to be Forgotten:</b> Penghapusan atas permintaan</li>
<li><b>Secure Deletion:</b> Cryptographic erasure dan overwriting</li>
<li><b>Certificate of Destruction:</b> Bukti penghapusan data</li>
</ul>

<h2>9. COOKIES DAN TRACKING</h2>
<h3>9.1 Jenis Cookies</h3>
<ul>
<li><b>Essential Cookies:</b> Untuk fungsi dasar aplikasi</li>
<li><b>Performance Cookies:</b> Analytics dan monitoring</li>
<li><b>Functional Cookies:</b> Preferensi dan personalization</li>
<li><b>Marketing Cookies:</b> Targeted advertising</li>
</ul>

<h3>9.2 Cookie Management</h3>
<ul>
<li><b>Cookie Settings:</b> Control di aplikasi settings</li>
<li><b>Browser Settings:</b> Disable cookies di browser</li>
<li><b>Third Party Cookies:</b> Opt-out dari advertising cookies</li>
<li><b>Cookie Policy:</b> Detail lengkap di cookie policy page</li>
</ul>

<h2>10. ANAK DI BAWAH UMUR</h2>
<ul>
<li>Layanan tidak ditujukan untuk anak di bawah 17 tahun</li>
<li>Verifikasi usia saat registrasi</li>
<li>Immediate deletion jika ditemukan data anak</li>
<li>Parental consent untuk usia 17-21 tahun</li>
</ul>

<h2>11. PERUBAHAN KEBIJAKAN PRIVASI</h2>
<h3>11.1 Notifikasi Perubahan</h3>
<ul>
<li><b>Material Changes:</b> Email notification 30 hari sebelumnya</li>
<li><b>Minor Updates:</b> In-app notification</li>
<li><b>Version History:</b> Archive versi sebelumnya</li>
<li><b>Continued Use:</b> Deemed acceptance jika tetap menggunakan</li>
</ul>

<h2>12. KONTAK DATA PROTECTION</h2>
<h3>12.1 Data Protection Officer</h3>
<ul>
<li><b>Email:</b> privacy@mbankingcore.com</li>
<li><b>Telepon:</b> +62-21-5000-1234</li>
<li><b>Alamat:</b> Menara MBankingCore, Jakarta 12345</li>
<li><b>Response Time:</b> Maksimal 30 hari untuk complex request</li>
</ul>

<h3>12.2 Complaint Process</h3>
<ul>
<li><b>Internal Complaint:</b> Melalui customer service</li>
<li><b>Regulator Complaint:</b> Ke Kementerian Kominfo</li>
<li><b>International:</b> GDPR representative untuk EU residents</li>
<li><b>Escalation:</b> Clear escalation process untuk unresolved issues</li>
</ul>

<h2>13. KETENTUAN KHUSUS</h2>
<h3>13.1 Emerging Technology</h3>
<ul>
<li><b>AI/ML Ethics:</b> Responsible AI untuk decision making</li>
<li><b>Blockchain:</b> Privacy consideration untuk DLT</li>
<li><b>IoT Integration:</b> Security untuk connected devices</li>
<li><b>Quantum Computing:</b> Quantum-safe cryptography preparation</li>
</ul>

<hr>
<p><b>Efektif sejak:</b> 31 Juli 2025<br>
<b>Versi:</b> 2.1<br>
<b>Bahasa:</b> Bahasa Indonesia (versi resmi), English (reference)</p>

<p><i>Dengan menggunakan layanan MBankingCore, Anda menyatakan telah membaca, memahami, dan menyetujui Kebijakan Privasi ini.</i></p>`
}
