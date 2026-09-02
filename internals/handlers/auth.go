package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hareshkhan01/PollyRoute/internals/service/auth"
)

type AuthHandlers struct {
	authService *auth.AuthService
}

func NewAuthHandlers(authService *auth.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

func (a *AuthHandlers) Login(c *gin.Context) {

}

func (a *AuthHandlers) Register(c *gin.Context) {

}

func (a *AuthHandlers) Refresh(c *gin.Context) {

}

func (a *AuthHandlers) Logout(c *gin.Context) {

}
