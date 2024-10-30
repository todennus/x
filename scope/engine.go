package scope

import (
	"strings"
)

type Engine struct {
	valid map[string]Scoper
}

func NewEngine() *Engine {
	return &Engine{
		valid: make(map[string]Scoper),
	}
}

func Define[S Scoper](engine *Engine, scope S) S {
	engine.valid[scope.Scope()] = scope
	return scope
}

func (engine *Engine) ParseScope(s string) (Scoper, bool) {
	if scope, ok := engine.valid[s]; ok {
		return scope, true
	}

	return New(s), false
}

func (engine *Engine) ParseDefinedScopes(s string) Scopes {
	s = strings.Trim(s, " ")
	if s == "" {
		return nil
	}

	scopesStr := strings.Split(s, " ")
	scopes := Scopes{}
	for _, str := range scopesStr {
		if scope, ok := engine.ParseScope(str); ok {
			scopes = append(scopes, scope)
		}
	}

	return scopes
}

func (engine *Engine) ParseUndefinedScopes(s string) Scopes {
	s = strings.Trim(s, " ")
	if s == "" {
		return nil
	}

	scopesStr := strings.Split(s, " ")
	scopes := Scopes{}
	for _, str := range scopesStr {
		if scope, ok := engine.ParseScope(str); !ok {
			scopes = append(scopes, scope)
		}
	}

	return scopes
}

func (engine *Engine) ParseAnyScopes(s string) Scopes {
	s = strings.Trim(s, " ")
	if s == "" {
		return nil
	}

	scopesStr := strings.Split(s, " ")
	scopes := Scopes{}
	for _, str := range scopesStr {
		scope, _ := engine.ParseScope(str)
		scopes = append(scopes, scope)
	}

	return scopes
}
