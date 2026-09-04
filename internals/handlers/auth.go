package handlers

import (
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

func (a *AuthHandlers) Login(c *gin.Context) {
	var oldUser request.OldUser

	if err := c.ShouldBindBodyWithJSON(&oldUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
	}

	accessToken, refreshToken, err := a.authService.LoginUser(c.Request.Context(), oldUser.Email, oldUser.RawPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})

}

func (a *AuthHandlers) Register(c *gin.Context) {
	var newUser request.NewUser

	if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	userId, err := a.authService.RegisterUser(c.Request.Context(), newUser.Name, newUser.Email, newUser.RawPassword)

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

func (a *AuthHandlers) Refresh(c *gin.Context) {
	// c.JSON(http.StatusForbidden, gin.H{"error": "currently unavailable"})
	var refresh struct {
		Token string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&refresh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Login Again"})
		return
	}

	accessToken, err := a.authService.RefreshToken(c.Request.Context(), refresh.Token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"accessToken": accessToken})

}

func (a *AuthHandlers) Logout(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		log.Println("Logout:", userId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "logout request failed"})
		return
	}
	err := a.authService.LogOutUser(c.Request.Context(), string(userId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Logout"})
}
