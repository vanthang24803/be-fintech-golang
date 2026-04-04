package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/docs"
	"github.com/maynguyen24/sever/internal/middleware"
	"github.com/maynguyen24/sever/internal/router"
	"github.com/maynguyen24/sever/pkg/logger"
	"github.com/maynguyen24/sever/pkg/response"
	"go.uber.org/zap"
)

// customErrorHandler acts as an interceptor for all uncaught HTTP handler errors
var customErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log the error centrally
	logger.Log.Error("Request Failed", zap.Error(err))

	// Map HTTP code to business logic error code (e.g. 500 -> 5000)
	businessCode := code * 10
	return response.Error(c, code, businessCode, message, err.Error())
}

type Server struct {
	App  *fiber.App
	port string
	cfg  *configs.Config
}

type Option func(*Server)

// WithPort sets the server port dynamically
func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithConfig injects the centralized configuration
func WithConfig(cfg *configs.Config) Option {
	return func(s *Server) {
		s.cfg = cfg
	}
}

// NewServer configures a new Server instance overriding defaults with options
func NewServer(opts ...Option) *Server {
	s := &Server{
		port: "8386", // default
	}

	for _, opt := range opts {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Setup layout middlewares
	app.Use(middleware.RequestLogger())

	// Public API reference and raw OpenAPI document
	docs.RegisterRoutes(app)

	// Init custom routes
	router.SetupRoutes(app, s.cfg)

	app.Get("/", func(c *fiber.Ctx) error {
		logger.Log.Info("Testing API")
		return response.Success(c, 2000, "Hello World from server!", nil)
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "This is a simulated validation error")
	})

	s.App = app
	return s
}

// Start launches the Fiber server
func (s *Server) Start() error {
	return s.App.Listen(":" + s.port)
}
