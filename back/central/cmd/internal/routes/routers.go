//go:build legacy

package routes

import (
	"central/internal/app/usecaseauth"
	"central/internal/domain/ports"
	"central/internal/infra/primary/http2/middleware"
	"central/internal/infra/primary/http2/routes"
	"central/internal/pkg/env"
	"central/internal/pkg/log"
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	server      *http.Server
	router      *gin.Engine
	logger      log.ILogger
	handlers    *routes.Handlers
	listener    net.Listener
	env         env.IConfig
	jwtService  ports.IJWTService
	authUseCase usecaseauth.IUseCaseAuth
}

func New(
	address string,
	logger log.ILogger,
	handlers *routes.Handlers,
	env env.IConfig,
	jwtService ports.IJWTService,
	authUseCase usecaseauth.IUseCaseAuth,
) (*HTTPServer, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// Configurar modo de Gin
	configureGinMode(env)

	// Configurar logger personalizado
	ginLogger := routes.NewGinLogger(logger)
	gin.DefaultWriter = ginLogger

	// Crear router con middleware
	router := createRouter(logger, env)

	httpServer := &HTTPServer{
		router:      router,
		logger:      logger,
		handlers:    handlers,
		listener:    lis,
		jwtService:  jwtService,
		authUseCase: authUseCase,
		server: &http.Server{
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		env: env,
	}

	return httpServer, nil
}

// configureGinMode configura el modo de Gin según el entorno
func configureGinMode(env env.IConfig) {
	if env.Get("LOG_LEVEL") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// createRouter crea y configura el router de Gin
func createRouter(logger log.ILogger, env env.IConfig) *gin.Engine {
	router := gin.New()

	// Middleware global
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.LoggingMiddleware(logger, env))

	return router
}

func (s *HTTPServer) Routers() {

	basicRoutes := routes.NewBasicRoutes(s.router, s.logger, s.env)
	basicRoutes.Setup()
	apiRoutes := routes.NewAPIRoutes(s.router, s.handlers, s.jwtService, s.authUseCase, s.logger)
	apiRoutes.Setup()
}

func (s *HTTPServer) Start() error {
	return s.server.Serve(s.listener)
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
