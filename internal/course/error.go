package course

import (
	"errors"
	"fmt"
)

var errNameRequired = errors.New("Name is required")
var errStartDateRequired = errors.New("Start date is required")
var errEndDateRequired = errors.New("End date is required")

type ErrNotFound struct {
	courseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Course '%s' doesn't exist", e.courseID)
}
