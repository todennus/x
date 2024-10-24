package conversion

import (
	"errors"
	"fmt"
	"strconv"
)

func ToInt(a any, strict bool) (int64, error) {
	if a == nil {
		a = 0
	}

	switch t := a.(type) {
	case int:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return int64(t), nil
	case float32:
		if t == float32(int64(t)) {
			return int64(t), nil
		}

		return 0, errors.New("expected int, got float")
	case float64:
		if t == float64(int64(t)) {
			return int64(t), nil
		}

		return 0, errors.New("expected int, got float")
	case string:
		if !strict {
			if t == "" {
				return 0, nil
			}

			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				return n, nil
			}
		}

		return 0, errors.New("expected int, got string")
	case bool:
		return 0, errors.New("expected int, got bool")
	default:
		return 0, fmt.Errorf("not handle for %T", t)
	}
}

func ToFloat(a any, strict bool) (float64, error) {
	if a == nil {
		a = 0.0
	}

	switch t := a.(type) {
	case int:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float32:
		return float64(t), nil
	case float64:
		return float64(t), nil
	case string:
		if !strict {
			if t == "" {
				return 0, nil
			}

			if n, err := strconv.ParseFloat(t, 64); err == nil {
				return n, nil
			}
		}

		return 0, errors.New("expected float, got string")
	case bool:
		return 0, errors.New("expected float, got bool")
	default:
		return 0, fmt.Errorf("not handle for %T", t)
	}
}

func ToBool(a any, strict bool) (bool, error) {
	if a == nil {
		a = false
	}

	switch t := a.(type) {
	case int:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case int8:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case int16:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case int32:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case int64:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case float32:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case float64:
		if !strict {
			switch t {
			case 0:
				return false, nil
			case 1:
				return true, nil
			}
		}
		return false, errors.New("expected bool, got number")
	case string:
		if !strict {
			if t == "" {
				return false, nil
			}

			if b, err := strconv.ParseBool(t); err == nil {
				return b, nil
			}
		}

		return false, errors.New("expected bool, got string")
	case bool:
		return t, nil
	default:
		return false, fmt.Errorf("not handle for %T", t)
	}
}

func ToString(a any, strict bool) (string, error) {
	if a == nil {
		a = ""
	}

	switch t := a.(type) {
	case int:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", errors.New("expected string, got number")
	case int8:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", errors.New("expected string, got number")
	case int16:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", errors.New("expected string, got number")
	case int32:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", errors.New("expected string, got number")
	case int64:
		if !strict {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return "", errors.New("expected string, got number")
	case float32:
		if !strict {
			return strconv.FormatFloat(float64(t), 'f', -1, 32), nil
		}
		return "", errors.New("expected string, got number")
	case float64:
		if !strict {
			return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
		}
		return "", errors.New("expected string, got number")
	case string:
		return t, nil
	case bool:
		if !strict {
			return strconv.FormatBool(t), nil
		}
		return "", errors.New("expected string, got bool")
	default:
		return "", fmt.Errorf("not handle for %T", t)
	}
}
