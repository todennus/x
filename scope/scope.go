package scope

import (
	"strings"
)

type Scoper interface {
	Scope() string
}

type Scope struct {
	value string
}

func New(value string) *Scope {
	return &Scope{value: value}
}

func (scope *Scope) Scope() string {
	return scope.value
}

type Scopes []Scoper

func NewScopes(s ...Scoper) Scopes {
	return s
}

func (s Scopes) String() string {
	a := []string{}
	for i := range s {
		a = append(a, s[i].Scope())
	}

	return strings.Join(a, " ")
}

func (scopes Scopes) Contains(target ...Scoper) bool {
	if len(scopes) == 0 {
		return len(target) == 0
	}

	for i := range target {
		contains := false
		for j := range scopes {
			if Equal(target[i], scopes[j]) {
				contains = true
				break
			}
		}

		if !contains {
			return false
		}
	}

	return true
}

func Equal(a, b Scoper) bool {
	return a.Scope() == b.Scope()
}
