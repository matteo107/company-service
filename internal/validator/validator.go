package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) AddError(field, message string) {
	v.Errors[field] = message
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) != 0
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Check(b bool, field, message string) {
	if !b {
		v.AddError(field, message)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}
