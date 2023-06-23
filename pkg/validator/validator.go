package validator

type (
	Validator interface {
		Validate(v interface{}) *ValidationError
	}

	ValidationError struct {
		OriginalMessage string
		Message         string
		Violations      []Violation
	}
	Violation struct {
		Namespace string      `json:"-"`
		Field     string      `json:"-"`
		FieldJSON string      `json:"field"`
		Tag       string      `json:"error"`
		Value     interface{} `json:"value,omitempty"`
	}
)
