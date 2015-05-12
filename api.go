package vpanel

type Validatable interface {
	Validate() error
}

type ValidationError struct {
	Errors []error
}

func (v *ValidationError) AddError(err error) {
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

func (v *ValidationError) Error() string {
	message := ""
	for _, err := range v.Errors {
		if len(message) > 0 {
			message += " "
		}
		message += err.Error()
	}
	return message
}
