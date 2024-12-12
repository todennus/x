package conversion

import (
	"github.com/xybor-x/snowflake"
)

func ConvertFromPointer[T any](p *T) T {
	if p == nil {
		var defaultT T
		return defaultT
	}

	return *p
}

func ConvertToPointer[T any](p T) *T {
	return &p
}

func MakeSnowflakePointerString(p *snowflake.ID) *string {
	if p == nil {
		return nil
	}

	s := p.String()
	return &s
}
