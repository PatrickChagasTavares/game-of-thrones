package characters

import "errors"

var (
	ErrFind          = errors.New("failed to find lord")
	ErrHouseNotFound = errors.New("this lord is not found or deleted")
)
