package application

import (
	"errors"
	"fmt"

	"github.com/narinsak-u/backend/internal/auth/domain"
)

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrUserNotFound           = errors.New("user not found")
	ErrPasswordHashing        = errors.New("failed to hash password")
	ErrTokenGeneration        = errors.New("failed to generate token")
)

type CurrentUserRepository interface {
	FindByID(uint) (*domain.User, error)
}

type UserRepository interface {
	CurrentUserRepository
	Create(*domain.User) error
	FindByEmail(string) (*domain.User, error)
}
type Passwords interface {
	Hash(string) (string, error)
	Matches(string, string) bool
}
type Tokens interface {
	Generate(uint, string) (string, error)
}

type Service struct {
	repo UserRepository
	passwords Passwords
	tokens Tokens
}

func NewService(repo UserRepository, passwords Passwords, tokens Tokens) *Service {
	return &Service{repo: repo, passwords: passwords, tokens: tokens}
}

func (s *Service) Register(name, email, password string) (*domain.User, string, error) {
	hash, err := s.passwords.Hash(password)
	if err != nil { return nil, "", fmt.Errorf("%w: %v", ErrPasswordHashing, err) }
	user := &domain.User{Name:name, Email:email, Password:hash, Role:domain.RoleCustomer}
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, ErrEmailAlreadyRegistered) { return nil, "", ErrEmailAlreadyRegistered }
		return nil, "", fmt.Errorf("create user: %w", err)
	}
	token, err := s.tokens.Generate(user.ID, user.Role)
	if err != nil { return nil, "", fmt.Errorf("%w: %v", ErrTokenGeneration, err) }
	return user, token, nil
}

func (s *Service) Login(email, password string) (*domain.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { return nil, "", ErrInvalidCredentials }
		return nil, "", fmt.Errorf("find user by email: %w", err)
	}
	if !s.passwords.Matches(password, user.Password) { return nil, "", ErrInvalidCredentials }
	token, err := s.tokens.Generate(user.ID, user.Role)
	if err != nil { return nil, "", fmt.Errorf("%w: %v", ErrTokenGeneration, err) }
	return user, token, nil
}

func (s *Service) Profile(id uint) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if errors.Is(err, ErrUserNotFound) { return nil, ErrUserNotFound }
	if err != nil { return nil, fmt.Errorf("find user by id: %w", err) }
	return user, nil
}
