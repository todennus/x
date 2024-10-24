package xerror

import (
	"errors"
)

func Is(err error, target error, others ...error) bool {
	if errors.Is(err, target) {
		return true
	}

	for _, target := range others {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func Iss(err error, targets ...error) bool {
	if len(targets) == 0 {
		panic("invalid targets")
	}

	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}
