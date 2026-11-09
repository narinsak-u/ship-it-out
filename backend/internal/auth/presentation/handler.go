package presentation

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/auth/application"
	"github.com/narinsak-u/backend/internal/auth/domain"
	"github.com/narinsak-u/backend/internal/platform/config"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Workflows interface {
	Register(string,string,string) (*domain.User,string,error)
	Login(string,string) (*domain.User,string,error)
	Profile(uint) (*domain.User,error)
}
type Handler struct { service Workflows }
func NewHandler(service Workflows) *Handler { return &Handler{service:service} }
type registerRequest struct { Name string `json:"name"`; Email string `json:"email"`; Password string `json:"password"`; Role string `json:"role"` }
type loginRequest struct { Email string `json:"email"`; Password string `json:"password"` }
type safeUser struct { ID uint `json:"id"`; Name string `json:"name"`; Email string `json:"email"`; Role string `json:"role"`; CreatedAt time.Time `json:"created_at"` }
func response(user *domain.User) safeUser { return safeUser{user.ID,user.Name,user.Email,user.Role,user.CreatedAt} }
func cookie(c *fiber.Ctx, value string, maxAge int) {
	c.Cookie(&fiber.Cookie{Name:"jwt", Value:value, Path:"/", HTTPOnly:true, SameSite:"None", Secure:true, MaxAge:maxAge})
}
func (h *Handler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err:=c.BodyParser(&req); err!=nil { return utils.Error(c,400,"invalid request body") }
	if req.Name=="" || req.Email=="" || req.Password=="" { return utils.Error(c,400,"name, email, and password are required") }
	user,token,err:=h.service.Register(req.Name,req.Email,req.Password)
	if errors.Is(err,application.ErrEmailAlreadyRegistered) { return utils.Error(c,409,"email already registered") }
	if errors.Is(err,application.ErrPasswordHashing) { return utils.Error(c,500,"failed to hash password") }
	if errors.Is(err,application.ErrTokenGeneration) { return utils.Error(c,500,"failed to generate token") }
	if err!=nil { return utils.Error(c,500,"failed to register user") }
	cookie(c,token,int(config.App.JWTTTL.Seconds()))
	return utils.Success(c,fiber.Map{"user":response(user)})
}
func (h *Handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err:=c.BodyParser(&req); err!=nil { return utils.Error(c,400,"invalid request body") }
	user,token,err:=h.service.Login(req.Email,req.Password)
	if errors.Is(err,application.ErrInvalidCredentials) { return utils.Error(c,401,"invalid email or password") }
	if errors.Is(err,application.ErrTokenGeneration) { return utils.Error(c,500,"failed to generate token") }
	if err!=nil { return utils.Error(c,500,"failed to authenticate") }
	cookie(c,token,int(config.App.JWTTTL.Seconds()))
	return utils.Success(c,fiber.Map{"token":token,"user":response(user)})
}
func (h *Handler) Me(c *fiber.Ctx) error {
	id,ok:=c.Locals("user_id").(uint); if !ok { return utils.Error(c,401,"not authenticated") }
	user,err:=h.service.Profile(id)
	if errors.Is(err,application.ErrUserNotFound) { return utils.Error(c,404,"user not found") }
	if err!=nil { return utils.Error(c,500,"failed to load user profile") }
	return utils.Success(c,response(user))
}
func (h *Handler) Logout(c *fiber.Ctx) error { cookie(c,"",0); return utils.Success(c,fiber.Map{"message":"logged out"}) }
