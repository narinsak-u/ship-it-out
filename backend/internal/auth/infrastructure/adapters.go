package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/narinsak-u/backend/internal/auth/application"
	"github.com/narinsak-u/backend/internal/auth/domain"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
	"gorm.io/gorm"
)

type GormRepository struct { db *gorm.DB }
func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db:db} }
func (r *GormRepository) Create(user *domain.User) error {
	row := models.User{Name:user.Name, Email:user.Email, Password:user.Password, Role:user.Role}
	if err := r.db.Create(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) { return application.ErrEmailAlreadyRegistered }
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { return application.ErrEmailAlreadyRegistered }
		return err
	}
	user.ID, user.CreatedAt = row.ID, row.CreatedAt
	return nil
}
func (r *GormRepository) FindByEmail(email string) (*domain.User, error) {
	var row models.User
	if err := r.db.Where("email = ?", email).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, application.ErrUserNotFound }
		return nil, err
	}
	return fromModel(row), nil
}
func (r *GormRepository) FindByID(id uint) (*domain.User, error) {
	var row models.User
	if err := r.db.First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, application.ErrUserNotFound }
		return nil, err
	}
	return fromModel(row), nil
}
func fromModel(row models.User) *domain.User {
	return &domain.User{ID:row.ID, Name:row.Name, Email:row.Email, Password:row.Password, Role:row.Role, CreatedAt:row.CreatedAt}
}
type PasswordAdapter struct{}
func (PasswordAdapter) Hash(value string) (string, error) { return utils.HashPassword(value) }
func (PasswordAdapter) Matches(value, hash string) bool { return utils.CheckPassword(value, hash) }
type JWTAdapter struct { secret []byte; ttl time.Duration }
func NewJWTAdapter(secret string, ttl time.Duration) *JWTAdapter { return &JWTAdapter{secret:[]byte(secret), ttl:ttl} }
func (a *JWTAdapter) Generate(userID uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id":userID, "role":role, "exp":time.Now().Add(a.ttl).Unix()})
	return token.SignedString(a.secret)
}
