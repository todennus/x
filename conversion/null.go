package conversion

import (
	"github.com/todennus/x/enum"
	"github.com/xybor-x/snowflake"
)

func ConvertPointer[T any](p *T) T {
	if p == nil {
		var defaultT T
		return defaultT
	}

	return *p
}

func MakeSnowflakePointerString(p *snowflake.ID) *string {
	if p == nil {
		return nil
	}

	s := p.String()
	return &s
}

func MakeEnumPointerString[T any](p *enum.Enum[T]) *string {
	if p == nil {
		return nil
	}

	s := p.String()
	return &s
}
