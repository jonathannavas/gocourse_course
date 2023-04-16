package course

import (
	"errors"
	"fmt"
)

var errNameRequired = errors.New("name is required")
var errStartDateRequired = errors.New("start date is required")
var errEndDateRequired = errors.New("end date is required")
var errDateValidation = errors.New("start date must be before end date")

type ErrNotFound struct {
	courseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Course '%s' doesn't exist", e.courseID)
}
