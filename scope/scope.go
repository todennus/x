package scope

import (
	"fmt"
)

type scopeRelationship int

const (
	scopeRelationshipNone scopeRelationship = iota
	scopeRelationshipEqual
	scopeRelationshipSubset
	scopeRelationshipSuperset
)

const (
	optionalPrefix = "?"
)

type Scoper interface {
	IsUndefined() bool
	IsOptional() bool
	Contains(another Scoper) bool
	String() string
}

var _ Scoper = Scope{}

type Scope struct {
	engine *Engine

	title      string
	isOptional bool
	action     Actioner
	resource   Resourcer
}

func newScope(engine *Engine, action Actioner, resource Resourcer) Scope {
	return Scope{
		engine:   engine,
		action:   action,
		resource: resource,
	}
}

func (scope Scope) WithOptional(optional bool) Scope {
	return Scope{
		engine:     scope.engine,
		action:     scope.action,
		resource:   scope.resource,
		isOptional: optional,
		title:      scope.title,
	}
}

func (scope Scope) WithTitle(title string) Scope {
	return Scope{
		engine:     scope.engine,
		action:     scope.action,
		resource:   scope.resource,
		isOptional: scope.isOptional,
		title:      title,
	}
}

func (scope Scope) AsScopes() Scopes {
	return NewScopes(scope)
}

func (scope Scope) Title() string {
	return scope.title
}

func (scope Scope) IsOptional() bool {
	return scope.isOptional
}

// String returns the scope as the format ?todennus/title:read:write
func (scope Scope) String() string {
	s := ""
	if scope.isOptional {
		s += optionalPrefix
	}

	if scope.engine.namespace != "" {
		s += scope.engine.namespace + "/"
	}

	if scope.title != "" && len(scope.engine.titles) > 0 && scope.title != scope.engine.titles[0] {
		s += scope.title + ":"
	}

	s += scope.action.String()

	if resource := scope.resource.String(); resource != "" {
		s += ":" + resource
	}

	return s
}

func (scope Scope) Contains(another Scoper) bool {
	anotherScope, ok := another.(Scope)
	if !ok {
		return false
	}

	if scope.title != anotherScope.title {
		return false
	}

	return anotherScope.action.IsSubset(scope.action) && anotherScope.resource.IsSubset(scope.resource)
}

func (scope Scope) IsUndefined() bool {
	return false
}

var _ Scoper = UndefinedScope{}

type UndefinedScope struct {
	value    string
	optional bool
}

func NewUndefinedScope(value string) UndefinedScope {
	return UndefinedScope{value: value}
}

func (scope UndefinedScope) String() string {
	prefix := ""
	if scope.IsOptional() {
		prefix = optionalPrefix
	}

	return fmt.Sprintf("%s%s", prefix, scope.value)
}

func (scope UndefinedScope) Contains(another Scoper) bool {
	if !another.IsUndefined() {
		return false
	}

	return scope.value == another.String()
}

func (scope UndefinedScope) IsUndefined() bool {
	return true
}

func (scope UndefinedScope) WithOptional(optional bool) UndefinedScope {
	return UndefinedScope{
		value:    scope.value,
		optional: optional,
	}
}

func (scope UndefinedScope) IsOptional() bool {
	return scope.optional
}

func relationship(scope, another Scoper) scopeRelationship {
	containAnother := scope.Contains(another)
	anotherContain := another.Contains(scope)

	switch {
	case containAnother && anotherContain:
		return scopeRelationshipEqual
	case containAnother:
		return scopeRelationshipSuperset
	case anotherContain:
		return scopeRelationshipSubset
	default:
		return scopeRelationshipNone
	}
}
