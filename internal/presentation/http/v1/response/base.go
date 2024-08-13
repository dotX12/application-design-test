package response

import (
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1]

	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", d.Format("2006-01-02"))
	return []byte(formatted), nil
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}
