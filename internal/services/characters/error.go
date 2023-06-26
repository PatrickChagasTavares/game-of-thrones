package characters

import "errors"

var (
	ErrFind              = errors.New("failed to find character")
	ErrCharacterNotFound = errors.New("this character is not found or deleted")
)
