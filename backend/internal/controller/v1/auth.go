package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/request"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/controller/v1/response"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/service"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authRoutes struct {
	log *slog.Logger

	authService service.Auth
	userService service.User
}

func newAuthRoutes(log *slog.Logger, g *gin.RouterGroup, authService service.Auth, userService service.User) {
	log = log.With(slog.String("component", "auth routes"))

	r := &authRoutes{
		log:         log,
		authService: authService,
		userService: userService,
	}

	g.POST("/login", r.login)
	g.POST("/refresh", r.refresh)
}

// @Summary Логин
// @Description Логин
// @Tags аутентификация
// @Accept json
// @Produce json
// @Param input body request.Login true "Тело запроса"
// @Success 201 {object} response.Login
// @Router /auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req request.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := validator.New().Struct(req); err != nil {
		r.log.Error("failed to validate request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	user, err := r.userService.GetByLogin(c, req.Login)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			r.log.Error("no such user", logger.Error(err), slog.String("login", req.Login))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		r.log.Error("failed to find user by login", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := r.authService.GenerateTokens(c, user.ID)
	if err != nil {
		if errors.Is(err, service.ErrWrongPassword) || errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		r.log.Error("failed to generate user token", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.Login{
		User:         response.BuildUser(*user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary Регенерация токена
// @Description Регенерация токена
// @Tags аутентификация
// @Accept json
// @Produce json
// @Param input body request.Refresh true "Тело запроса"
// @Success 200 {object} response.Refresh
// @Router /auth/refresh [post]
func (r *authRoutes) refresh(c *gin.Context) {
	var req request.Refresh

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := validator.New().Struct(req); err != nil {
		r.log.Error("failed to validate request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	refreshClaims, err := r.authService.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		r.log.Error("failed to parse refresh token", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := r.userService.GetByID(c, refreshClaims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			r.log.Error("failed to find user by login", logger.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		}

		r.log.Error("failed to find user by login", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := r.authService.GenerateTokens(c, user.ID)
	if err != nil {
		if errors.Is(err, service.ErrWrongPassword) || errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		r.log.Error("failed to generate user token", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.Refresh{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
