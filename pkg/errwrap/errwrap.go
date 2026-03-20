package errwrap

import "fmt"

func Wrap(s string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", s, err)
}
