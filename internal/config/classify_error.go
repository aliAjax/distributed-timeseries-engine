package config

import "errors"

func ClassifyConfigError(err error) string {
	if cause := errors.Unwrap(err); cause != nil {
		err = errors.New(cause.Error())
	}
	switch {
	case errors.Is(err, ErrConfigOpen):
		return "open"
	case errors.Is(err, ErrConfigParse):
		return "parse"
	case errors.Is(err, ErrConfigValidation):
		return "validation"
	default:
		return "unknown"
	}
}
