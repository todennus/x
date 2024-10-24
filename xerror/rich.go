package xerror

import (
	"fmt"
)

type richErrorCore struct {
	code        error
	description string
	detail      error
	event       string
	attributes  []any
}

type RichError struct {
	core *richErrorCore
}

func (err RichError) Error() string {
	if err.core.detail == nil {
		return fmt.Sprintf("%s:%s", err.core.code, err.core.description)
	} else {
		return fmt.Sprintf("%s:%s (%s:%s) %v",
			err.core.code,
			err.core.description,
			err.core.event,
			err.core.detail,
			err.core.attributes,
		)
	}
}

func (err RichError) Code() error {
	return err.core.code
}

func (err RichError) Description() string {
	return err.core.description
}

func (err RichError) Event() string {
	return err.core.event
}

func (err RichError) Detail() error {
	return err.core.detail
}

func (err RichError) Attributes() []any {
	return err.core.attributes
}

func (err RichError) Unwrap() error {
	return err.core.code
}

func (err RichError) Reduce() error {
	return RichError{
		core: &richErrorCore{
			code:        err.core.code,
			description: err.core.description,
		},
	}
}

func Enrich(code error, description string, format ...any) RichError {
	return RichError{core: &richErrorCore{code: code, description: fmt.Sprintf(description, format...)}}
}

func (err RichError) Hide(detail error, event string, attributes ...any) RichError {
	if len(attributes)%2 != 0 {
		panic("invalid key-value attributes")
	}

	err.core.event = event
	err.core.detail = detail
	err.core.attributes = attributes
	return err
}
