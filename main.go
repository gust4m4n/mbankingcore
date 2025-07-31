package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"mbankingcore/config"
	"mbankingcore/handlers"
	"mbankingcore/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// generateSelfSignedCert creates a self-signed certificate for development
func generateSelfSignedCert(certFile, keyFile string) error {
	// Check if certificate files already exist
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			log.Println("✅ SSL certificate files already exist")
			return nil
		}
	}

	log.Println("🔐 Generating self-signed SSL certificate for development...")

	// Create certificates directory if it doesn't exist
	certDir := filepath.Dir(certFile)
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return err
	}

	// Generate certificate using openssl command
	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096",
		"-keyout", keyFile, "-out", certFile, "-days", "365", "-nodes",
		"-subj", "/C=ID/ST=Jakarta/L=Jakarta/O=MBankingCore/OU=Development/CN=localhost")

	if err := cmd.Run(); err != nil {
		// If openssl is not available, create a basic message
		log.Println("⚠️  OpenSSL not found. Please install OpenSSL or provide your own certificates")
		log.Printf("📁 Place your certificate files at: %s and %s", certFile, keyFile)
		log.Println("💡 To install OpenSSL on macOS: brew install openssl")
		return err
	}

	log.Println("✅ Self-signed SSL certificate generated successfully")
	return nil
}

// startHTTPSServer starts the server with HTTPS support
func startHTTPSServer(router *gin.Engine, address, certFile, keyFile string) error {
	// Create TLS configuration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Create HTTP server with TLS
	server := &http.Server{
		Addr:      address,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	log.Printf("🔒 Starting HTTPS server on %s", address)
	log.Printf("🔐 Using certificate: %s", certFile)
	log.Printf("🔑 Using private key: %s", keyFile)

	return server.ListenAndServeTLS(certFile, keyFile)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Connect to database
	config.ConnectDatabase()

	// Initialize Gin router
	router := gin.Default()

	// Add comprehensive CORS middleware for all APIs
	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Log comprehensive CORS request details for debugging
		log.Printf("CORS Request - Method: %s, Origin: %s, Path: %s, User-Agent: %s",
			c.Request.Method, origin, c.Request.URL.Path, c.GetHeader("User-Agent"))

		// Set CORS headers FIRST before any other processing
		// Always set CORS headers regardless of method or origin
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Device-ID, X-App-Version, Accept, Accept-Language, Content-Language, DNT, User-Agent, Keep-Alive, Request-Id, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization, X-Total-Count")

		// Handle preflight OPTIONS requests immediately
		if c.Request.Method == "OPTIONS" {
			log.Printf("CORS Preflight - Origin: %s, Request-Headers: %s, Request-Method: %s",
				origin, c.GetHeader("Access-Control-Request-Headers"), c.GetHeader("Access-Control-Request-Method"))

			// Respond immediately to preflight
			c.Status(204)
			return
		}

		// Continue to next middleware
		c.Next()
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(config.DB)
	articleHandler := handlers.NewArticleHandler(config.DB)
	photoHandler := handlers.NewPhotoHandler(config.DB)
	bankAccountHandler := handlers.NewBankAccountHandler(config.DB)
	adminHandler := handlers.NewAdminHandler(config.DB)
	transactionHandler := handlers.NewTransactionHandler(config.DB)
	auditHandler := handlers.NewAuditHandler()

	// API routes
	api := router.Group("/api")

	// Add audit logging middleware
	api.Use(middleware.AuditLogMiddleware())

	{
		// Authentication routes (public)
		api.POST("/login", middleware.AuditLoginMiddleware(), authHandler.BankingLogin)              // Banking Login Step 1 - Send OTP
		api.POST("/login/verify", middleware.AuditLoginMiddleware(), authHandler.BankingLoginVerify) // Banking Login Step 2 - Verify OTP
		api.POST("/refresh", authHandler.RefreshToken)                                               // Refresh token

		// Public onboarding routes (remain public)
		api.GET("/onboardings", handlers.GetOnboardings)    // Get all onboardings (public)
		api.GET("/onboardings/:id", handlers.GetOnboarding) // Get onboarding by ID (public)

		// Public terms and conditions routes (config-based)
		api.GET("/terms-conditions", handlers.GetTermsConditions) // Get terms and conditions from config (public)

		// Terms and conditions management (authenticated users)
		api.POST("/terms-conditions", middleware.AuthMiddleware(), handlers.SetTermsConditions) // Set terms and conditions content (authenticated users)

		// Public privacy policy routes (config-based)
		api.GET("/privacy-policy", handlers.GetPrivacyPolicy) // Get privacy policy from config (public)

		// Privacy policy management (authenticated users)
		api.POST("/privacy-policy", middleware.AuthMiddleware(), handlers.SetPrivacyPolicy) // Set privacy policy content (authenticated users)

		// Admin authentication routes (public)
		admin := api.Group("/admin")
		{
			admin.POST("/login", middleware.AuditLoginMiddleware(), adminHandler.AdminLogin)   // Admin login
			admin.POST("/logout", middleware.AuditLoginMiddleware(), adminHandler.AdminLogout) // Admin logout

			// Admin protected routes
			adminProtected := admin.Group("/")
			adminProtected.Use(middleware.AdminAuthMiddleware())
			{
				// Dashboard (admin only)
				adminProtected.GET("/dashboard", adminHandler.GetDashboard) // Get dashboard statistics

				// Admin management
				adminProtected.GET("/admins", adminHandler.GetAdmins)          // Get all admins
				adminProtected.GET("/admins/:id", adminHandler.GetAdminByID)   // Get admin by ID
				adminProtected.POST("/admins", adminHandler.CreateAdmin)       // Create new admin
				adminProtected.PUT("/admins/:id", adminHandler.UpdateAdmin)    // Update admin
				adminProtected.DELETE("/admins/:id", adminHandler.DeleteAdmin) // Delete admin

				// Transaction monitoring (admin only)
				adminProtected.GET("/transactions", transactionHandler.GetAllTransactions) // Get all transactions for monitoring
				adminProtected.POST("/transactions/reversal", transactionHandler.Reversal) // Reverse a transaction

				// Audit trails (admin only)
				adminProtected.GET("/audit-logs", auditHandler.GetAuditLogs)        // Get audit logs with filtering
				adminProtected.GET("/login-audits", auditHandler.GetLoginAuditLogs) // Get login audit logs with filtering

				// Config management (admin only)
				adminProtected.POST("/config", handlers.SetConfig)           // Set config value (admin only)
				adminProtected.GET("/configs", handlers.GetAllConfigs)       // Get all configs (admin only)
				adminProtected.GET("/config/:key", handlers.GetConfig)       // Get config value by key (admin only)
				adminProtected.DELETE("/config/:key", handlers.DeleteConfig) // Delete config by key (admin only)
			}
		} // Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile management
			protected.GET("/profile", authHandler.Profile)       // Get user profile
			protected.PUT("/profile", authHandler.UpdateProfile) // Update user profile
			protected.PUT("/change-pin", authHandler.ChangePIN)  // Change PIN ATM

			// Session management
			protected.GET("/sessions", authHandler.GetActiveSessions)                                            // Get active sessions
			protected.POST("/logout", middleware.AuditLoginMiddleware(), authHandler.Logout)                     // Logout
			protected.POST("/logout-others", middleware.AuditLoginMiddleware(), authHandler.LogoutOtherSessions) // Logout other sessions

			// Article management (all operations require authentication)
			protected.GET("/articles", articleHandler.GetArticles)          // Get all articles (protected)
			protected.GET("/articles/:id", articleHandler.GetArticleByID)   // Get article by ID (protected)
			protected.PUT("/articles/:id", articleHandler.UpdateArticle)    // Update article (user can update own)
			protected.DELETE("/articles/:id", articleHandler.DeleteArticle) // Delete article (user can delete own)
			protected.GET("/my-articles", articleHandler.GetMyArticles)     // Get my articles

			// Photo management (all operations require authentication)
			protected.GET("/photos", photoHandler.GetPhotos)          // Get all photos (protected)
			protected.GET("/photos/:id", photoHandler.GetPhotoByID)   // Get photo by ID (protected)
			protected.PUT("/photos/:id", photoHandler.UpdatePhoto)    // Update photo (user can update own)
			protected.DELETE("/photos/:id", photoHandler.DeletePhoto) // Delete photo (user can delete own)

			// Bank account management (authenticated users)
			protected.GET("/bank-accounts", bankAccountHandler.GetBankAccounts)               // Get user's bank accounts
			protected.POST("/bank-accounts", bankAccountHandler.CreateBankAccount)            // Create new bank account
			protected.PUT("/bank-accounts/:id", bankAccountHandler.UpdateBankAccount)         // Update bank account
			protected.DELETE("/bank-accounts/:id", bankAccountHandler.DeleteBankAccount)      // Delete bank account
			protected.PUT("/bank-accounts/:id/primary", bankAccountHandler.SetPrimaryAccount) // Set primary account

			// Article management (authenticated users)
			protected.POST("/articles", articleHandler.CreateArticle) // Create article (authenticated users)

			// Onboarding management (authenticated users)
			protected.POST("/onboardings", handlers.CreateOnboarding)       // Create onboarding (authenticated users)
			protected.PUT("/onboardings/:id", handlers.UpdateOnboarding)    // Update onboarding (authenticated users)
			protected.DELETE("/onboardings/:id", handlers.DeleteOnboarding) // Delete onboarding (authenticated users)

			// Photo management (authenticated users)
			protected.POST("/photos", photoHandler.CreatePhoto) // Create photo (authenticated users)

			// User management - basic operations (authenticated users)
			protected.GET("/users", handlers.ListUsers)         // List all users (authenticated users)
			protected.GET("/users/:id", handlers.GetUserByID)   // Get user by ID (authenticated users)
			protected.DELETE("/users/:id", handlers.DeleteUser) // Delete user by ID (authenticated users)

			// Transaction management (authenticated users)
			protected.POST("/transactions/topup", transactionHandler.Topup)                // Top up balance
			protected.POST("/transactions/withdraw", transactionHandler.Withdraw)          // Withdraw balance
			protected.POST("/transactions/transfer", transactionHandler.Transfer)          // Transfer balance to other user
			protected.GET("/transactions/history", transactionHandler.GetUserTransactions) // Get user transaction history
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "MBankingCore API is running",
			"data": gin.H{
				"status": "ok",
			},
		})
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get host from environment or default to all interfaces for Flutter compatibility
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0" // Bind to all interfaces to allow Flutter emulator/device access
	}

	// Check if HTTPS is enabled
	enableHTTPS := os.Getenv("ENABLE_HTTPS")
	httpsPort := os.Getenv("HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = "8443"
	}

	address := host + ":" + port
	httpsAddress := host + ":" + httpsPort

	// Certificate file paths
	certDir := os.Getenv("CERT_DIR")
	if certDir == "" {
		certDir = "./certs"
	}
	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	if enableHTTPS == "true" || enableHTTPS == "1" {
		// HTTPS Mode
		log.Println("🔒 HTTPS mode enabled")

		// Generate self-signed certificate if not exists
		if err := generateSelfSignedCert(certFile, keyFile); err != nil {
			log.Printf("❌ Failed to generate SSL certificate: %v", err)
			log.Println("🔄 Falling back to HTTP mode...")
			enableHTTPS = "false"
		} else {
			log.Printf("🚀 Starting HTTPS server on %s", httpsAddress)
			log.Printf("🔐 Health check available at: https://%s/health", httpsAddress)
			log.Printf("🔒 API base URL: https://%s/api", httpsAddress)

			// Start HTTPS server
			if err := startHTTPSServer(router, httpsAddress, certFile, keyFile); err != nil {
				log.Fatal("❌ Failed to start HTTPS server:", err)
			}
		}
	}

	if enableHTTPS != "true" && enableHTTPS != "1" {
		// HTTP Mode (default)
		log.Println("🌐 HTTP mode enabled")
		log.Printf("🚀 Starting HTTP server on %s", address)
		log.Printf("📋 Health check available at: http://%s/health", address)
		log.Printf("🔗 API base URL: http://%s/api", address)

		if err := router.Run(address); err != nil {
			log.Fatal("❌ Failed to start HTTP server:", err)
		}
	}
}
