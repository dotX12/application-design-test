package vo

import "fmt"

type Quota struct {
	value int
}

var ErrQuotaNegative = fmt.Errorf("quota cannot be negative")

func (q *Quota) validate() error {
	if q.value < 0 {
		return ErrQuotaNegative
	}
	return nil
}

func (q *Quota) String() string {
	return fmt.Sprintf("%d", q.value)
}

func (q *Quota) Value() int {
	return q.value
}

func NewQuotaFromInt(value int) (*Quota, error) {
	quota := Quota{value: value}
	if err := quota.validate(); err != nil {
		return nil, err
	}
	return &quota, nil
}
