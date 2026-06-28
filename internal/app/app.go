package app

import (
	"github.com/gin-gonic/gin"
	"github.com/product-management-server/config"
	"github.com/product-management-server/internal/handler"
	"github.com/product-management-server/internal/infrastructure"
	"github.com/product-management-server/internal/repository"
	"github.com/product-management-server/internal/routes"
	"github.com/product-management-server/internal/service"
)

// App holds all wired dependencies and the HTTP server.
type App struct {
	Router *gin.Engine
	cfg    *config.Config
}

// New wires all dependencies (DB → Repository → Service → Handler → Router)
// and returns a ready-to-run App.
func New(cfg *config.Config) (*App, error) {
	// Initialize database connection
	db, err := infrastructure.NewDatabaseConnection(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Wire layers
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	// Setup Gin router with default middleware (logger + recovery)
	router := gin.Default()
	routes.Register(router, h)

	return &App{
		Router: router,
		cfg:    cfg,
	}, nil
}

// Run starts the HTTP server on the configured port.
func (a *App) Run() error {
	return a.Router.Run(a.cfg.Port)
}
