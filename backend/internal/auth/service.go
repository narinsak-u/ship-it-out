package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

var (
	ErrInvalidCredentials       = errors.New("invalid email or password")
	ErrEmailAlreadyRegistered   = errors.New("email already registered")
	ErrUserNotFound             = errors.New("user not found")
	ErrInvalidTokenClaims       = errors.New("invalid token claims")
	ErrPasswordHashing          = errors.New("failed to hash password")
	ErrTokenGeneration          = errors.New("failed to generate token")
)

type Service struct {
	repo   Repository
	secret []byte
	ttl    time.Duration
}

func NewService(repo Repository, secret string, ttl time.Duration) *Service {
	return &Service{repo: repo, secret: []byte(secret), ttl: ttl}
}

func (s *Service) Register(name, email, password string) (*models.User, Token, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, Token{}, fmt.Errorf("%w: %v", ErrPasswordHashing, err)
	}
	user := &models.User{Name: name, Email: email, Password: hash, Role: "customer"}
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, ErrEmailAlreadyRegistered) {
			return nil, Token{}, ErrEmailAlreadyRegistered
		}
		return nil, Token{}, fmt.Errorf("create user: %w", err)
	}
	token, err := NewToken(user.ID, user.Role, s.secret, s.ttl)
	if err != nil {
		return nil, Token{}, fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}
	return user, token, nil
}

func (s *Service) Login(email, password string) (*models.User, Token, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, Token{}, ErrInvalidCredentials
		}
		return nil, Token{}, fmt.Errorf("find user by email: %w", err)
	}
	if !utils.CheckPassword(password, user.Password) {
		return nil, Token{}, ErrInvalidCredentials
	}
	token, err := NewToken(user.ID, user.Role, s.secret, s.ttl)
	if err != nil {
		return nil, Token{}, fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}
	return user, token, nil
}

func (s *Service) Profile(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if errors.Is(err, ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}


// Token keeps the signed token string out of auth workflow and verification logic.
type Token struct{ value string }

func NewToken(userID uint, role string, secret []byte, ttl time.Duration) (Token, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	value, err := token.SignedString(secret)
	if err != nil {
		return Token{}, err
	}
	return Token{value: value}, nil
}

func (t Token) String() string { return t.value }

type TokenClaims struct {
	UserID uint
	Role   string
}

func ParseToken(value string, secret []byte) (TokenClaims, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return TokenClaims{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenClaims{}, ErrInvalidTokenClaims
	}
	userID, ok := claims["user_id"].(float64)
	if !ok || userID <= 0 {
		return TokenClaims{}, ErrInvalidTokenClaims
	}
	role, ok := claims["role"].(string)
	if !ok {
		return TokenClaims{}, ErrInvalidTokenClaims
	}
	return TokenClaims{UserID: uint(userID), Role: role}, nil
}
