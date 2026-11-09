package presentation

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/auth/domain"
)

type workflowsStub struct { user *domain.User }
func (s workflowsStub) Register(name,email,password string)(*domain.User,string,error) { return s.user,"signed-token",nil }
func (s workflowsStub) Login(string,string)(*domain.User,string,error) { return s.user,"signed-token",nil }
func (s workflowsStub) Profile(uint)(*domain.User,error) { return s.user,nil }

func TestRegisterIgnoresRequestedRoleAndHidesPassword(t *testing.T) {
	app:=fiber.New()
	app.Post("/api/auth/register",NewHandler(workflowsStub{&domain.User{ID:3,Name:"N",Email:"n@example.com",Password:"secret-hash",Role:domain.RoleCustomer}}).Register)
	request:=httptest.NewRequest("POST","/api/auth/register",strings.NewReader(`{"name":"N","email":"n@example.com","password":"pw","role":"admin"}`))
	request.Header.Set("Content-Type","application/json")
	response,err:=app.Test(request)
	if err!=nil { t.Fatal(err) }
	defer response.Body.Close()
	var body map[string]any
	if err:=json.NewDecoder(response.Body).Decode(&body); err!=nil { t.Fatal(err) }
	data:=body["data"].(map[string]any)
	user:=data["user"].(map[string]any)
	if user["role"]!="customer" { t.Fatalf("role = %v",user["role"]) }
	if _,ok:=user["password"]; ok { t.Fatal("password exposed") }
	if response.Header.Get("Set-Cookie")=="" { t.Fatal("auth cookie missing") }
}
