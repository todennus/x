package scope

import (
	"strings"
)

type Scopes []Scoper

func NewScopes(scopes ...Scoper) Scopes {
	return scopes
}

func (scopes Scopes) String() string {
	scopeValues := []string{}
	for i := range scopes {
		scopeValues = append(scopeValues, scopes[i].String())
	}

	return strings.Join(scopeValues, " ")
}

func (scopes Scopes) Contains(target Scoper) bool {
	for _, scope := range scopes {
		if scope.Contains(target) {
			return true
		}
	}

	return false
}

func (scopes Scopes) LessThanOrEqual(another Scopes) bool {
	for _, scope := range scopes {
		isSubset := false
		for _, targetScope := range another {
			relationship := relationship(scope, targetScope)
			if relationship == scopeRelationshipSubset || relationship == scopeRelationshipEqual {
				isSubset = true
				break
			}
		}

		if !isSubset {
			return false
		}
	}

	return true
}

func (scopes Scopes) Intersect(another Scopes) Scopes {
	result := Scopes{}
	for _, scope := range scopes {
		if scope.IsUndefined() {
			result = append(result, scope)
			continue
		}

		for _, targetScope := range another {
			if relationship := relationship(scope, targetScope); relationship != scopeRelationshipNone {
				if relationship == scopeRelationshipSubset {
					result = append(result, scope)
				} else {
					result = append(result, targetScope)
				}

				break
			}
		}
	}

	return result
}
