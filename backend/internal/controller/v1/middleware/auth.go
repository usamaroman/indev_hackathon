package middleware

import (
	"net/http"
	"strings"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	UserIDCtx   = "userID"
	UserTypeCtx = "userType"
	HotelIDCtx  = "hotelID"
)

type AuthMiddleware struct {
	authService service.Auth
}

func bearerToken(ctx *gin.Context) (string, bool) {
	header := ctx.GetHeader("Authorization")
	parts := strings.Split(header, " ")

	if len(parts) != 2 {
		return "", false
	}

	if parts[0] != "Bearer" {
		return "", false
	}

	return parts[1], true
}

func NewAuthMiddleware(authService service.Auth) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) AuthOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims.UserType == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		c.Set(UserTypeCtx, claims.UserType)
		c.Set(UserIDCtx, claims.UserID)
		c.Set(HotelIDCtx, claims.HotelID)

		c.Next()
	}
}

func (m *AuthMiddleware) HotelsOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims.UserType != types.AdminUser.String() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no roots",
			})
			return
		}

		c.Set(UserTypeCtx, claims.UserType)
		c.Set(UserIDCtx, claims.UserID)
		c.Set(HotelIDCtx, claims.HotelID)

		c.Next()
	}
}

func (m *AuthMiddleware) CustomersOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims.UserType != types.Customer.String() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no roots",
			})
			return
		}

		c.Set(UserTypeCtx, claims.UserType)
		c.Set(UserIDCtx, claims.UserID)
		c.Set(HotelIDCtx, claims.HotelID)

		c.Next()
	}
}

func (m *AuthMiddleware) BusinessAndCustomersOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims.UserType != types.Customer.String() && claims.UserType != types.AdminUser.String() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no roots",
			})
			return
		}

		c.Set(UserTypeCtx, claims.UserType)
		c.Set(UserIDCtx, claims.UserID)
		c.Set(HotelIDCtx, claims.HotelID)

		c.Next()
	}
}

func (m *AuthMiddleware) PublicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.Set(UserTypeCtx, "customer")
			c.Next()
			return
		}

		claims, err := m.authService.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set(UserTypeCtx, claims.UserType)
		c.Set(UserIDCtx, claims.UserID)
		c.Set(HotelIDCtx, claims.HotelID)

		c.Next()
	}
}
