package customerror

import (
	"errors"
	"fmt"
	"reflect"
)

func NotFoundError(resource string, key interface{}, inner error) error {
	t := reflect.TypeOf(key)
	v := reflect.ValueOf(key)

	if errors.Is(inner, nil) {
		//nolint: wrapcheck,err113
		return fmt.Errorf("%s not found %v(%q)", resource, t.Name(), v)
	}

	//nolint: wrapcheck
	return fmt.Errorf("%s not found %v(%q): %w", resource, t.Name(), v, inner)
}
