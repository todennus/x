package scope

import (
	"fmt"
	"strings"
)

type Engine struct {
	namespace       string
	namespacePrefix string
	actionMap       map[string]Actioner
	resourceMap     map[string]Resourcer
	titles          []string
	validScope      map[string]bool
}

func NewEngine(namespace string, actionMap map[string]Actioner, resourceMap map[string]Resourcer) Engine {
	return Engine{
		namespace:       namespace,
		namespacePrefix: fmt.Sprintf("%s/", namespace),
		actionMap:       actionMap,
		resourceMap:     resourceMap,
		titles:          make([]string, 0),
		validScope:      make(map[string]bool),
	}
}

func (engine *Engine) DefineTitle(title string) {
	engine.titles = append(engine.titles, title)
}

func (engine *Engine) New(action Actioner, resource Resourcer) Scope {
	return newScope(engine, action, resource)
}

func (engine *Engine) ParseScope(s string) Scoper {
	s = strings.Trim(s, " ")
	if s == "" {
		return nil
	}

	isOptional := false
	if strings.HasPrefix(s, optionalPrefix) {
		isOptional = true
		s = s[len(optionalPrefix):]
	}

	if !strings.HasPrefix(s, engine.namespacePrefix) {
		return NewUndefinedScope(s).WithOptional(isOptional)
	}

	s = s[len(engine.namespacePrefix):]

	detectedTitle := ""
	if len(engine.titles) > 0 {
		detectedTitle = engine.titles[0]
	}

	for i := range engine.titles {
		titlePrefix := fmt.Sprintf("%s:", engine.titles[i])
		if strings.HasPrefix(s, titlePrefix) {
			detectedTitle = engine.titles[i]
			s = s[len(titlePrefix):]
			break
		}
	}

	actionStr, resourceStr, found := strings.Cut(s, ":")
	if !found {
		actionStr = s
		resourceStr = ""
	}

	action, ok := engine.actionMap[actionStr]
	if !ok {
		return NewUndefinedScope(s).WithOptional(isOptional)
	}

	resource, ok := engine.resourceMap[resourceStr]
	if !ok {
		return NewUndefinedScope(s).WithOptional(isOptional)
	}

	scope := newScope(engine, action, resource).WithOptional(isOptional).WithTitle(detectedTitle)
	return scope
}

func (engine *Engine) ParseScopes(s string) Scopes {
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

type validScopeDefiner struct {
	engine *Engine
	scope  Scope
}

func (engine *Engine) DefineScope(action Actioner, resource Resourcer) *validScopeDefiner {
	return &validScopeDefiner{
		engine: engine,
		scope:  engine.New(action, resource),
	}
}

func (d *validScopeDefiner) WithTitle(title string) *validScopeDefiner {
	d.engine.validScope[d.scope.WithTitle(title).String()] = true
	return d
}
