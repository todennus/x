package scope_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todennus/x/scope"
)

const namespace = "test"

type UserScope struct {
	s *scope.Scope
}

func NewUserScope(s string) *UserScope {
	return &UserScope{s: scope.New(s)}
}

func (s *UserScope) Scope() string {
	return fmt.Sprintf("%s/%s", namespace, s.s.Scope())
}

type AdminScope struct {
	s *scope.Scope
}

func NewAdminScope(s string) *AdminScope {
	return &AdminScope{s: scope.New(s)}
}

func (s *AdminScope) Scope() string {
	return fmt.Sprintf("%s/admin:%s", namespace, s.s.Scope())
}

var Engine = scope.NewEngine()

var (
	UserReadUser          = scope.Define(Engine, NewUserScope("read:user"))
	UserUpdateUserAvatar  = scope.Define(Engine, NewUserScope("update:user.avatar"))
	UserReadClientProfile = scope.Define(Engine, NewUserScope("read:client.profile"))
	UserReadClientOwner   = scope.Define(Engine, NewUserScope("read:client.owner"))
	UserCreateClient      = scope.Define(Engine, NewUserScope("create:client"))

	AdminReadUser         = scope.Define(Engine, NewAdminScope("read:user"))
	AdminUpdateUserAvatar = scope.Define(Engine, NewAdminScope("update:user.avatar"))
)

func Test_SerializeScopes(t *testing.T) {
	scopes := scope.NewScopes(UserReadUser, UserUpdateUserAvatar)

	assert.Equal(t, "test/read:user test/update:user.avatar", scopes.String())
	assert.Equal(t, "", scope.NewScopes().String())
}

func Test_Engine_ParseScopes(t *testing.T) {
	scopes := Engine.ParseAnyScopes("test/read:user test/update:user.avatar")

	assert.True(t, scopes.Contains(UserReadUser))
	assert.True(t, scopes.Contains(UserUpdateUserAvatar))

	assert.False(t, scopes.Contains(UserReadClientOwner))
	assert.False(t, scopes.Contains(UserCreateClient))
}

func Test_Engine_ParseScopes_Empty(t *testing.T) {
	scopes := Engine.ParseAnyScopes("")

	assert.False(t, scopes.Contains(UserReadUser))
	assert.False(t, scopes.Contains(UserReadClientProfile))
}

func Test_Engine_ParseScopes_Outside_WithoutNamespace(t *testing.T) {
	scopes := Engine.ParseAnyScopes("read:user update:user.avatar")

	assert.False(t, scopes.Contains(UserReadUser))
	assert.False(t, scopes.Contains(UserUpdateUserAvatar))

	assert.False(t, scopes.Contains(UserReadClientOwner))
	assert.False(t, scopes.Contains(UserCreateClient))
}

func Test_ScopesNil(t *testing.T) {
	scopes := scope.Scopes(nil)
	assert.Equal(t, "", scopes.String())
}

func Test_Scopes_Title(t *testing.T) {
	sc, _ := Engine.ParseScope("test/admin:read:user")
	assert.True(t, scope.Equal(sc, AdminReadUser))

	sc, _ = Engine.ParseScope("test/admin:update:user.avatar")
	assert.True(t, scope.Equal(sc, AdminUpdateUserAvatar))

	sc, _ = Engine.ParseScope("test/read:user")
	assert.True(t, scope.Equal(sc, UserReadUser))

	sc, _ = Engine.ParseScope("test/admin:client")
	_, ok := sc.(*scope.Scope)
	assert.True(t, ok)

	sc, _ = Engine.ParseScope("test/admin:")
	_, ok = sc.(*scope.Scope)
	assert.True(t, ok)

	sc, _ = Engine.ParseScope("test/admin")
	_, ok = sc.(*scope.Scope)
	assert.True(t, ok)
}

func Test_Scopes_FullContext(t *testing.T) {
	sc, _ := Engine.ParseScope("test/admin:create:client")
	_, ok := sc.(*scope.Scope)
	assert.True(t, ok)
	assert.Equal(t, "test/admin:create:client", sc.Scope())

	sc, _ = Engine.ParseScope("test/create:client")
	_, ok = sc.(*UserScope)
	assert.True(t, ok)
	assert.Equal(t, "test/create:client", sc.Scope())
}
