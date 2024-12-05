package xerror

import "fmt"

func Join(err error, other error) error {
	return JoinByFormat(err, other, "%w | %w")
}

func JoinByFormat(err error, other error, format string) error {
	if err == nil {
		return other
	}

	if other == nil {
		return err
	}

	return fmt.Errorf(format, err, other)
}
