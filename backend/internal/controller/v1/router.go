package v1

import (
	"log/slog"
	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/usamaroman/demo_indev_hackathon/backend/docs"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/middleware"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(log *slog.Logger, router *gin.Engine, services *service.Services) {
	router.Use(middleware.CORS(log))
	router.Use(middleware.Log(log))
	router.Use(middleware.Metrics())

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	authMiddleware := middleware.NewAuthMiddleware(services.Auth)

	authGroup := router.Group("/auth")
	{
		newAuthRoutes(log, authGroup, services.Auth, services.User)
	}

	v1 := router.Group("/v1")

	newHotelRoutes(log, v1.Group("/hotel"), services.Hotel, authMiddleware)
}
