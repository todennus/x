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

type Scoper interface {
	Contains(another Scoper) bool

	String() string

	IsUndefined() bool

	IsOptional() bool
}

var _ Scoper = Scope{}

type Scope struct {
	source   string
	action   Actioner
	resource Resourcer
	optional bool
}

func New(source string, action Actioner, resource Resourcer) Scope {
	return Scope{
		source:   source,
		action:   action,
		resource: resource,
	}
}

func (scope Scope) WithOptional(optional bool) Scope {
	return Scope{
		source:   scope.source,
		action:   scope.action,
		resource: scope.resource,
		optional: optional,
	}
}

func (scope Scope) IsOptional() bool {
	return scope.optional
}

func (scope Scope) String() string {
	prefix := ""
	if scope.IsOptional() {
		prefix = "@"
	}

	if scope.resource.String() == "" {
		return fmt.Sprintf("%s[%s]%s", prefix, scope.source, scope.action.String())
	}

	return fmt.Sprintf("%s[%s]%s:%s", prefix, scope.source, scope.action.String(), scope.resource.String())
}

func (scope Scope) Contains(another Scoper) bool {
	anotherScope, ok := another.(Scope)
	if !ok {
		return false
	}

	return anotherScope.action.IsSubset(scope.action) && anotherScope.resource.IsSubset(scope.resource)
}

func (scope Scope) AsScopes() Scopes {
	return NewScopes(scope)
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
		prefix = "@"
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
