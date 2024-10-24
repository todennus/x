package scope_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todennus/x/scope"
)

type AllActions struct {
	*scope.BaseAction

	Read  *scope.BaseAction
	Write *WriteAction
}

type WriteAction struct {
	*scope.BaseAction

	Create *scope.BaseAction
	Update *scope.BaseAction
	Delete *scope.BaseAction
}

type AllResources struct {
	*scope.BaseResource

	Client *scope.BaseResource
	User   *UserResource
}

type UserResource struct {
	*scope.BaseResource

	Email       *scope.BaseResource
	DisplayName *scope.BaseResource `resource:"display_name"`
	Avatar      *scope.BaseResource
}

var ResourceAll, resourceMap = scope.DefineResource[AllResources]()
var ActionSet, actionMap = scope.DefineAction[AllActions]()
var Engine = scope.NewEngine("x", actionMap, resourceMap)

func Test_DefineAction(t *testing.T) {
	assert.Equal(t, ActionSet.String(), actionMap["*"].String())
	assert.Equal(t, ActionSet.Read.String(), actionMap["read"].String())
	assert.Equal(t, ActionSet.Read.String(), actionMap["read"].String())
	assert.Equal(t, ActionSet.Write.String(), actionMap["write"].String())
	assert.Equal(t, ActionSet.Write.Create.String(), actionMap["create"].String())
	assert.Equal(t, ActionSet.Write.Update.String(), actionMap["update"].String())
	assert.Equal(t, ActionSet.Write.Delete.String(), actionMap["delete"].String())

	assert.Equal(t, "*", ActionSet.String())
	assert.Equal(t, "read", ActionSet.Read.String())
	assert.Equal(t, "write", ActionSet.Write.String())
	assert.Equal(t, "create", ActionSet.Write.Create.String())
	assert.Equal(t, "update", ActionSet.Write.Update.String())
	assert.Equal(t, "delete", ActionSet.Write.Delete.String())

	assert.True(t, ActionSet.Write.Delete.IsSubset(ActionSet))
	assert.True(t, ActionSet.Write.Delete.IsSubset(ActionSet.Write))
	assert.True(t, ActionSet.Write.Delete.IsSubset(ActionSet.Write.Delete))
	assert.False(t, ActionSet.Write.Delete.IsSubset(ActionSet.Read))
	assert.False(t, ActionSet.Write.Delete.IsSubset(ActionSet.Write.Create))
}

func Test_DefineResource(t *testing.T) {
	assert.Equal(t, ResourceAll.String(), resourceMap[""].String())
	assert.Equal(t, ResourceAll.User.String(), resourceMap["user"].String())
	assert.Equal(t, ResourceAll.Client.String(), resourceMap["client"].String())
	assert.Equal(t, ResourceAll.User.Email.String(), resourceMap["user.email"].String())
	assert.Equal(t, ResourceAll.User.Avatar.String(), resourceMap["user.avatar"].String())
	assert.Equal(t, ResourceAll.User.DisplayName.String(), resourceMap["user.display_name"].String())

	assert.Equal(t, "", ResourceAll.String())

	assert.Equal(t, "user", ResourceAll.User.String())
	assert.Equal(t, "client", ResourceAll.Client.String())

	assert.Equal(t, "user.email", ResourceAll.User.Email.String())
	assert.Equal(t, "user.display_name", ResourceAll.User.DisplayName.String())
	assert.Equal(t, "user.avatar", ResourceAll.User.Avatar.String())

	assert.True(t, ResourceAll.User.IsSubset(ResourceAll))
	assert.True(t, ResourceAll.User.Email.IsSubset(ResourceAll))
	assert.True(t, ResourceAll.Client.IsSubset(ResourceAll))

	assert.False(t, ResourceAll.Client.IsSubset(ResourceAll.User))
	assert.True(t, ResourceAll.User.Email.IsSubset(ResourceAll.User))
	assert.False(t, ResourceAll.User.Email.IsSubset(ResourceAll.Client))
}

func Test_SerializeScope(t *testing.T) {
	scopes := scope.NewScopes(
		Engine.New(ActionSet.Read, ResourceAll.User),
		Engine.New(ActionSet.Write, ResourceAll.User.Avatar),
	)

	assert.Equal(t, "[x]read:user [x]write:user.avatar", scopes.String())
	assert.Equal(t, "", scope.NewScopes().String())
}

func Test_Engine_ParseScope(t *testing.T) {
	engine := scope.NewEngine("custom", actionMap, resourceMap)
	scopes := engine.ParseScopes("[custom]read:user write:user.avatar")

	assert.True(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User)))
	assert.True(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User.Email)))

	assert.False(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User.Avatar)))
	assert.False(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User)))
}

func Test_Engine_ParseScopeEmpty(t *testing.T) {
	engine := scope.NewEngine("custom", actionMap, resourceMap)
	scopes := engine.ParseScopes("")

	assert.False(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User)))
	assert.False(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User.Email)))

	assert.False(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User.Avatar)))
	assert.False(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User)))
}

func Test_Engine_ParseScopeAll(t *testing.T) {
	engine := scope.NewEngine("custom", actionMap, resourceMap)
	scopes := engine.ParseScopes("[custom]*")

	assert.True(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User)))
	assert.True(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User.Email)))

	assert.True(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User.Avatar)))
	assert.True(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User)))
}

func Test_Engine_ParseScopeAll_Outside(t *testing.T) {
	engine := scope.NewEngine("custom", actionMap, resourceMap)
	scopes := engine.ParseScopes("*")

	assert.False(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User)))
	assert.False(t, scopes.Contains(engine.New(ActionSet.Read, ResourceAll.User.Email)))

	assert.False(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User.Avatar)))
	assert.False(t, scopes.Contains(engine.New(ActionSet.Write, ResourceAll.User)))
}

func Test_ScopesNil(t *testing.T) {
	scopes := scope.Scopes(nil)
	assert.Equal(t, "", scopes.String())
}

func Test_Scopes_Lessthan(t *testing.T) {
	parser := scope.NewEngine("z", actionMap, resourceMap)

	scopes := parser.ParseScopes("create:client something")
	assert.Equal(t, "create:client something", scopes.String())
	assert.False(t, scopes.LessThanOrEqual(parser.ParseScopes("[z]*")))

	scopes = parser.ParseScopes("[z]create:client")
	assert.Equal(t, "[z]create:client", scopes.String())
	assert.True(t, scopes.LessThanOrEqual(parser.ParseScopes("[z]*")))

	scopeA := parser.ParseScopes("[z]read:client")
	scopeB := parser.ParseScopes("[z]read:client [z]read:user")
	assert.True(t, scopeA.LessThanOrEqual(scopeB))

	scopeA = parser.ParseScopes("[z]read:client")
	scopeB = parser.ParseScopes("[z]read:client")
	assert.False(t, !scopeA.LessThanOrEqual(scopeB))
}
