package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/request"
	"github.com/hareshkhan01/PollyRoute/internals/service/auth"
)

type AuthHandlers struct {
	authService auth.AuthService
	// jwtSecret   string
}

// func NewAuthHandlers(authService auth.AuthService, jwtSecret string) *AuthHandlers {
func NewAuthHandlers(authService auth.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		// jwtSecret:   jwtSecret,
	}
}

func (a *AuthHandlers) Login(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var oldUser request.OldUser

		if err := c.ShouldBindBodyWithJSON(&oldUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		}

		accessToken, refreshToken, err := a.authService.LoginUser(ctx, oldUser.Email, oldUser.RawPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func (a *AuthHandlers) Register(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser request.NewUser

		if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
			return
		}

		userId, err := a.authService.RegisterUser(ctx, newUser.Name, newUser.Email, newUser.RawPassword)

		if userId != "" && err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error!"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Registerd Successfully"})

	}

}

func (a *AuthHandlers) Refresh(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusForbidden, gin.H{"error": "currently unavailable"})

	}
}

func (a *AuthHandlers) Logout(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.JSON(http.StatusForbidden, gin.H{"error": "currently unavailable"})

	}
}
