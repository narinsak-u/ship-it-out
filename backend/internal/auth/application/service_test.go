package application

import (
	"errors"
	"testing"

	"github.com/narinsak-u/backend/internal/auth/domain"
)

type testUsers struct { user *domain.User; createErr error }
func (r *testUsers) Create(user *domain.User) error { r.user=user; if r.createErr==nil { user.ID=9 }; return r.createErr }
func (r *testUsers) FindByEmail(string) (*domain.User,error) { return r.user,nil }
func (r *testUsers) FindByID(uint) (*domain.User,error) { return r.user,nil }
type testPasswords struct{}
func (testPasswords) Hash(value string)(string,error) { return "hash:"+value,nil }
func (testPasswords) Matches(value, hash string)bool { return hash=="hash:"+value }
type testTokens struct{}
func (testTokens) Generate(id uint, role string)(string,error) { return "token",nil }

func TestRegisterAlwaysCreatesCustomer(t *testing.T) {
	repo:=&testUsers{}
	service:=NewService(repo,testPasswords{},testTokens{})
	user,token,err:=service.Register("N","n@example.com","pw")
	if err!=nil { t.Fatal(err) }
	if user.Role!=domain.RoleCustomer || repo.user.Role!=domain.RoleCustomer { t.Fatalf("signup role = %q",user.Role) }
	if repo.user.Password!="hash:pw" || token!="token" { t.Fatalf("unexpected registration result: %#v %q",repo.user,token) }
}
func TestRegisterMapsDuplicateEmail(t *testing.T) {
	service:=NewService(&testUsers{createErr:ErrEmailAlreadyRegistered},testPasswords{},testTokens{})
	_,_,err:=service.Register("N","n@example.com","pw")
	if !errors.Is(err,ErrEmailAlreadyRegistered) { t.Fatalf("error = %v",err) }
}
