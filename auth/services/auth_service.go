package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-auth-backend/auth/config"
	"go-auth-backend/auth/models"
	"go-auth-backend/auth/repository"
)

type AuthService struct {
	Repo *repository.UserRepository
	Cfg  *config.Config
}

func NewAuthService(repo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{Repo: repo, Cfg: cfg}
}

func (s *AuthService) Register(email, password, name string) error {
	existing, _ := s.Repo.FindByEmail(email)
	if existing != nil {
		return errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashed),
	}

	return s.Repo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if user == nil {
		return "", errors.New("user not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.Cfg.JWTSecret))
}

func (s *AuthService) GetUserByToken(tokenStr string) (*models.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token data")
	}

	user, err := s.Repo.FindByID(sub)
	if err != nil {
		return nil, err
	}

	return user, nil
}
