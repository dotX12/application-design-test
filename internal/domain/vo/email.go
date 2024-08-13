package vo

import (
	"applicationDesignTest/internal/domain/exception"
	"regexp"
)

type Email struct {
	value string
}

func (e Email) String() string {
	return e.value
}

func NewEmailFromString(value string) (Email, error) {
	email := Email{value: value}
	if err := email.validate(); err != nil {
		return Email{}, err
	}
	return email, nil
}

func (e Email) validate() error {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(e.value) {
		return exception.ErrEmailFormat
	}
	return nil
}
