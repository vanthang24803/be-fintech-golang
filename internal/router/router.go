package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/handler"
	"github.com/maynguyen24/sever/internal/middleware"
	"github.com/maynguyen24/sever/internal/repository"
	"github.com/maynguyen24/sever/internal/service"
	"github.com/maynguyen24/sever/pkg/i18n"
	"github.com/maynguyen24/sever/pkg/logger"
	"github.com/maynguyen24/sever/pkg/mailer"
	"github.com/maynguyen24/sever/pkg/push"
	"github.com/maynguyen24/sever/pkg/queue"
	"github.com/maynguyen24/sever/pkg/upload"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	)

	// SetupRoutes wires up all dependencies and mounts the API endpoints
	func SetupRoutes(app *fiber.App, cfg *configs.Config, db *sqlx.DB) {
	// 0. Initialize I18n
	_ = i18n.LoadLocales("locales")

	// 1. Global Middlewares
	app.Use(middleware.I18nMiddleware())

	// Initialize Mailer
	var emailer mailer.Mailer
	if cfg.SMTPUser != "" && cfg.SMTPPass != "" {
		emailer = mailer.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	} else {
		logger.Log.Warn("SMTP credentials missing, using MockMailer")
		emailer = mailer.NewMockMailer()
	}

	// Initialize MinIO Uploader

	uploader, err := upload.NewMinioUploader(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucketName,
		cfg.MinioUseSSL,
	)
	if err != nil {
		logger.Log.Error("Failed to initialize MinIO uploader", zap.Error(err))
	}

	// 2. Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	sourceRepo := repository.NewSourcePaymentRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	fundRepo := repository.NewFundRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	reportRepo := repository.NewReportRepository(db)
	goalRepo := repository.NewSavingsGoalRepository(db)

	// 2. Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, tokenRepo, emailer, cfg)
	sourceService := service.NewSourcePaymentService(sourceRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize Push Client
	var pushClient push.PushClient
	firebaseClient, err := push.NewFirebaseClient(cfg.FirebaseCredentialsPath)
	if err != nil {
		// Fallback to Mock in dev/missing credentials
		pushClient = &push.MockPushClient{}
	} else {
		pushClient = firebaseClient
	}
	// Initialize Queue Client
	queueClient := queue.NewClient(cfg.RedisAddr, cfg.RedisPassword)

	notificationService := service.NewNotificationService(notificationRepo, deviceRepo, pushClient, queueClient)
	transactionService := service.NewTransactionService(transactionRepo, budgetRepo, notificationService, queueClient)
	fundService := service.NewFundService(fundRepo)
	deviceService := service.NewDeviceService(deviceRepo)
	budgetService := service.NewBudgetService(budgetRepo)
	reportService := service.NewReportService(reportRepo)
	goalService := service.NewSavingsGoalService(goalRepo, fundRepo, notificationService, queueClient)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	fidoService, _ := service.NewFIDOService(deviceRepo, redisClient, cfg)

	// 3. Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	sourceHandler := handler.NewSourcePaymentHandler(sourceService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	fundHandler := handler.NewFundHandler(fundService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	reportHandler := handler.NewReportHandler(reportService)
	goalHandler := handler.NewSavingsGoalHandler(goalService)
	fidoHandler := handler.NewFIDOHandler(fidoService)
	uploadHandler := handler.NewUploadHandler(uploader)

	// 4. API Routes Group
	api := app.Group("/api/v1")

	// Auth Routes (public) — all POST
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// Google OAuth
	auth.Post("/google/url", authHandler.GetGoogleAuthURL)
	auth.Post("/google/callback", authHandler.GoogleCallback)

	// FIDO2 Biometric Step-up Auth (public)
	auth.Post("/biometric/begin", fidoHandler.BeginAuthentication)
	auth.Post("/biometric/finish", fidoHandler.FinishAuthentication)

	// Protected Routes (require valid JWT)
	protected := api.Group("/", middleware.AuthRequired(cfg))

	// Upload — POST
	protected.Post("/upload", uploadHandler.Upload)

	// Profile (Update requires FIDO)
	protected.Post("/profile/me", userHandler.GetProfile)
	protected.Post("/profile/update", middleware.FIDOMiddleware(), userHandler.UpdateProfile)

	// Source Payment — all POST
	sources := protected.Group("/sources")
	sources.Post("/create", sourceHandler.Create)
	sources.Post("/list", sourceHandler.GetAll)
	sources.Post("/detail/:id", sourceHandler.GetByID)
	sources.Post("/update/:id", sourceHandler.Update)
	sources.Post("/delete/:id", sourceHandler.Delete)

	// Category — all POST
	categories := protected.Group("/categories")
	categories.Post("/create", categoryHandler.Create)
	categories.Post("/list", categoryHandler.GetAll)
	categories.Post("/detail/:id", categoryHandler.GetByID)
	categories.Post("/update/:id", categoryHandler.Update)
	categories.Post("/delete/:id", categoryHandler.Delete)

	// Transaction — all POST (body params replace query params for filters)
	transactions := protected.Group("/transactions")
	transactions.Post("/create", transactionHandler.Create)
	transactions.Post("/list", transactionHandler.GetAll)
	transactions.Post("/detail/:id", transactionHandler.GetByID)
	transactions.Post("/update/:id", transactionHandler.Update)
	transactions.Post("/delete/:id", transactionHandler.Delete)

	// Fund — all POST
	funds := protected.Group("/funds")
	funds.Post("/create", fundHandler.Create)
	funds.Post("/list", fundHandler.GetAll)
	funds.Post("/detail/:id", fundHandler.GetByID)
	funds.Post("/update/:id", fundHandler.Update)
	funds.Post("/delete/:id", fundHandler.Delete)
	funds.Post("/deposit/:id", fundHandler.Deposit)
	funds.Post("/withdraw/:id", middleware.FIDOMiddleware(), fundHandler.Withdraw)

	// Notification — all POST
	notifications := protected.Group("/notifications")
	notifications.Post("/list", notificationHandler.List)
	notifications.Post("/unread-count", notificationHandler.UnreadCount)
	notifications.Post("/mark-read", notificationHandler.MarkRead)
	notifications.Post("/delete/:id", notificationHandler.Delete)

	// Device — all POST
	devices := protected.Group("/devices")
	devices.Post("/register", deviceHandler.Register)
	devices.Post("/list", deviceHandler.List)
	devices.Post("/delete/:id", middleware.FIDOMiddleware(), deviceHandler.Delete)
	devices.Post("/biometric/enroll/:id/begin", fidoHandler.BeginEnrollment)
	devices.Post("/biometric/enroll/:id/finish", fidoHandler.FinishEnrollment)

	// Budget — all POST
	budgets := protected.Group("/budgets")
	budgets.Post("/create", budgetHandler.Create)
	budgets.Post("/list", budgetHandler.List)
	budgets.Post("/detail/:id", budgetHandler.GetDetail)
	budgets.Post("/update/:id", budgetHandler.Update)
	budgets.Post("/delete/:id", middleware.FIDOMiddleware(), budgetHandler.Delete)

	// Report — all POST
	reports := protected.Group("/reports")
	reports.Post("/category-summary", reportHandler.GetCategorySummary)
	reports.Post("/monthly-trend", reportHandler.GetMonthlyTrend)

	// Savings Goal — all POST
	goals := protected.Group("/goals")
	goals.Post("/create", goalHandler.Create)
	goals.Post("/list", goalHandler.List)
	goals.Post("/detail/:id", goalHandler.GetDetail)
	goals.Post("/contribute", goalHandler.Contribute)
	goals.Post("/withdraw", middleware.FIDOMiddleware(), goalHandler.Withdraw)
}
