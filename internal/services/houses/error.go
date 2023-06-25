package houses

import "errors"

var (
	ErrNameUsed      = errors.New("name informed already used in another house")
	ErrFind          = errors.New("failed to find houses")
	ErrHouseNotFound = errors.New("this house is not found or deleted")
)
