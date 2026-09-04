package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hareshkhan01/PollyRoute/internals/repository"
	"github.com/hareshkhan01/PollyRoute/internals/service/security"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, name string, email string, rawPassword string) (string, error)
	LoginUser(ctx context.Context, email string, rawPassword string) (string, string, error)
	LogOutUser(ctx context.Context, email string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type authService struct {
	userRepository repository.UserRepository
	jwtSecret      string
}

func NewAuthService(userRepository repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (a *authService) RegisterUser(ctx context.Context, name string, email string, rawPassword string) (string, error) {

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		log.Println("Failed to create password hash: %w", err)
		return "", fmt.Errorf("Failed to Register User")
	}
	userId, err := a.userRepository.CreateUser(ctx, name, email, string(hashedBytes))
	if err != nil {
		log.Println("User Registration Failed: ", err)
		return "", fmt.Errorf("User's email already exist!")
	}
	return userId, nil
}
func (a *authService) LoginUser(ctx context.Context, email string, rawPassword string) (string, string, error) {
	user, err := a.userRepository.FindByEmail(ctx, email)
	if err != nil {
		log.Println("No User Found: ", err)
		return "", "", fmt.Errorf("No User Found!")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rawPassword))
	if err != nil {
		log.Println("Incorrect Password:", err)
		return "", "", fmt.Errorf("Incorrect Password")
	}

	accessToke, err := security.GenerateJWtToken(user.ID, 30*time.Minute, a.jwtSecret)

	if err != nil {
		log.Println("Can not generate access token: %w", err)
		return "", "", fmt.Errorf("Failed to login")
	}

	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		log.Println("Can not generate refresh token: %w", err)
		return "", "", fmt.Errorf("Failed to login")
	}

	refresTokeExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	err = a.userRepository.UpdateRefreshToken(ctx, user.ID, &refreshToken, &refresTokeExpiresAt)

	if err != nil {
		log.Println("Failed to update refresh token inside db: ", err)
		return "", "", fmt.Errorf("Login Failed")
	}

	return accessToke, refreshToken, nil

}

func (a *authService) LogOutUser(ctx context.Context, userId string) error {
	err := a.userRepository.UpdateRefreshToken(ctx, userId, nil, nil)

	if err != nil {
		log.Println("Logout: Failed to set refresh token nil: ", err)
		return fmt.Errorf("Logout Failed")
	}

	return nil
}

func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	user, err := a.userRepository.FindByRefreshToken(ctx, refreshToken)

	if err != nil {
		log.Println("Failed fetch user by refresh token: ", err)
		return "", fmt.Errorf("Unable to create new token.")
	}

	if time.Now().After(*user.RefreshTokenExpiresAt) {
		log.Println("Refresh token expire")
		return "", fmt.Errorf("Session expired, Plese Login Again!")
	}

	accessToken, err := security.GenerateJWtToken(user.ID, 30*time.Minute, a.jwtSecret)

	if err != nil {
		log.Println("Failed to create access toke!")
		return "", fmt.Errorf("Falied: Please try again later!")
	}

	return accessToken, nil

}
