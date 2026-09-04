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

	user, err := a.userRepository.FindByEmail(ctx, email)

	if user != nil {
		return user.ID, fmt.Errorf("user's email already exist")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("Failed to create password hash: %w", err)
	}
	userId, err := a.userRepository.CreateUser(ctx, name, email, string(hashedBytes))
	if err != nil {
		return "", fmt.Errorf("User Registration Failed: %w", err)
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
		return "", "", fmt.Errorf("Can not generate access token: %w", err)
	}

	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("Can not generate refresh token: %w", err)
	}

	refresTokeExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	err = a.userRepository.UpdateRefreshToken(ctx, user.ID, &refreshToken, &refresTokeExpiresAt)

	if err != nil {
		return "", "", err
	}

	return accessToke, refreshToken, nil

}

func (a *authService) LogOutUser(ctx context.Context, userId string) error {
	err := a.userRepository.UpdateRefreshToken(ctx, userId, nil, nil)

	if err != nil {
		return fmt.Errorf("Logout Failed: %w", err)
	}

	return nil
}

func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	user, err := a.userRepository.FindByRefreshToken(ctx, refreshToken)

	if err != nil {
		log.Println("Failed fetch user by refresh token: ", err)
		return "", fmt.Errorf("Unable to create new token.")
	}

	if time.Now().Compare(*user.RefreshTokenExpiresAt) >= 1 {
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
