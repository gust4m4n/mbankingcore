package main

import (
	"log"
	"os"

	"mbankingcore/config"
	"mbankingcore/handlers"
	"mbankingcore/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Connect to database
	config.ConnectDatabase()

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(config.DB)
	articleHandler := handlers.NewArticleHandler(config.DB)
	photoHandler := handlers.NewPhotoHandler(config.DB)
	bankAccountHandler := handlers.NewBankAccountHandler(config.DB)

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (public)
		api.POST("/login", authHandler.BankingLogin)              // Banking Login Step 1 - Send OTP
		api.POST("/login/verify", authHandler.BankingLoginVerify) // Banking Login Step 2 - Verify OTP
		api.POST("/refresh", authHandler.RefreshToken)            // Refresh token

		// Public onboarding routes (remain public)
		api.GET("/onboardings", handlers.GetOnboardings)    // Get all onboardings (public)
		api.GET("/onboardings/:id", handlers.GetOnboarding) // Get onboarding by ID (public)

		// Public terms and conditions routes (config-based)
		api.GET("/terms-conditions", handlers.GetTermsConditions) // Get terms and conditions from config (public)

		// Admin terms and conditions management (config-based)
		api.POST("/terms-conditions", middleware.AdminMiddleware(), handlers.SetTermsConditions) // Set terms and conditions content (admin/owner only)

		// Public privacy policy routes (config-based)
		api.GET("/privacy-policy", handlers.GetPrivacyPolicy) // Get privacy policy from config (public)

		// Admin privacy policy management (config-based)
		api.POST("/privacy-policy", middleware.AdminMiddleware(), handlers.SetPrivacyPolicy) // Set privacy policy content (admin/owner only)

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile management
			protected.GET("/profile", authHandler.Profile)       // Get user profile
			protected.PUT("/profile", authHandler.UpdateProfile) // Update user profile
			protected.PUT("/change-pin", authHandler.ChangePIN)  // Change PIN ATM

			// Session management
			protected.GET("/sessions", authHandler.GetActiveSessions)         // Get active sessions
			protected.POST("/logout", authHandler.Logout)                     // Logout
			protected.POST("/logout-others", authHandler.LogoutOtherSessions) // Logout other sessions

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

			// Config management - get only (all authenticated users can read configs)
			protected.GET("/config/:key", handlers.GetConfig) // Get config value by key (authenticated users)
		} // Admin-only routes (require admin or owner role)
		admin := api.Group("/")
		admin.Use(middleware.AdminMiddleware())
		{
			// Article management (admin/owner only)
			admin.POST("/articles", articleHandler.CreateArticle) // Create article (admin/owner only)

			// Onboarding management (admin/owner only)
			admin.POST("/onboardings", handlers.CreateOnboarding)       // Create onboarding (admin/owner only)
			admin.PUT("/onboardings/:id", handlers.UpdateOnboarding)    // Update onboarding (admin/owner only)
			admin.DELETE("/onboardings/:id", handlers.DeleteOnboarding) // Delete onboarding (admin/owner only)

			// Photo management (admin/owner only)
			admin.POST("/photos", photoHandler.CreatePhoto) // Create photo (admin/owner only)

			// User management - basic operations (admin/owner only)
			admin.GET("/users", handlers.ListUsers)            // List all users (admin/owner only)
			admin.GET("/admin/users", handlers.ListAdminUsers) // List admin and owner users (admin/owner only)
			admin.GET("/users/:id", handlers.GetUserByID)      // Get user by ID (admin/owner only)
			admin.DELETE("/users/:id", handlers.DeleteUser)    // Delete user by ID (admin/owner only)

			// Config management (admin/owner only)
			admin.POST("/config", handlers.SetConfig)           // Set config value (admin/owner only)
			admin.GET("/admin/configs", handlers.GetAllConfigs) // Get all configs (admin/owner only)
			admin.DELETE("/config/:key", handlers.DeleteConfig) // Delete config by key (admin/owner only)
		}

		// Owner-only routes (require owner role)
		owner := api.Group("/")
		owner.Use(middleware.OwnerMiddleware())
		{
			// User management - role changes (owner only)
			owner.POST("/users", handlers.CreateUser)    // Create user (owner only - can set any role)
			owner.PUT("/users/:id", handlers.UpdateUser) // Update user by ID (owner only - can change roles)
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

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
