package device

import "fmt"

type Status struct {
	value string
}

var (
	StatusOffline = Status{"offline"}
	StatusOnline  = Status{"online"}
)

func NewStatus(value string) (Status, error) {
	switch value {
	case StatusOffline.value:
		return StatusOffline, nil
	case StatusOnline.value:
		return StatusOnline, nil
	default:
		return Status{}, fmt.Errorf("invalid device status: %s", value)
	}
}

func (s Status) String() string {
	return s.value
}
