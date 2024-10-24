package scope

import (
	"fmt"
	"strings"
)

type Engine struct {
	source      string
	actionMap   map[string]Actioner
	resourceMap map[string]Resourcer
}

func NewEngine(source string, actionMap map[string]Actioner, resourceMap map[string]Resourcer) Engine {
	return Engine{
		source:      source,
		actionMap:   actionMap,
		resourceMap: resourceMap,
	}
}

func (engine Engine) New(action Actioner, resource Resourcer) Scope {
	return New(engine.source, action, resource)
}

func (engine Engine) ParseScope(s string) Scoper {
	s = strings.Trim(s, " ")
	if s == "" {
		return nil
	}

	optional := false
	if s[0] == '@' {
		s = s[1:]
		optional = true
	}

	sourceStr := fmt.Sprintf("[%s]", engine.source)
	if !strings.HasPrefix(s, sourceStr) {
		return NewUndefinedScope(s).WithOptional(optional)
	}

	actionStr, resourceStr, found := strings.Cut(s[len(sourceStr):], ":")
	if !found {
		actionStr = s[len(sourceStr):]
		resourceStr = ""
	}

	action, ok := engine.actionMap[actionStr]
	if !ok {
		return NewUndefinedScope(s).WithOptional(optional)
	}

	resource, ok := engine.resourceMap[resourceStr]
	if !ok {
		return NewUndefinedScope(s).WithOptional(optional)
	}

	scope := New(engine.source, action, resource).WithOptional(optional)
	return scope
}

func (engine Engine) ParseScopes(s string) Scopes {
	s = strings.Trim(s, " ")
	if s == "" {
		return Scopes{}
	}

	scopesStr := strings.Split(s, " ")
	scopes := Scopes{}
	for _, str := range scopesStr {
		if s := engine.ParseScope(str); s != nil {
			scopes = append(scopes, s)
		}
	}

	return scopes
}
