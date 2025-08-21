package command

import "fmt"

type Status struct {
	value string
}

var (
	StatusPending  = Status{"pending"}
	StatusExecuted = Status{"executed"}
	StatusFailed   = Status{"failed"}
)

func NewStatus(value string) (Status, error) {
	switch value {
	case StatusPending.value:
		return StatusPending, nil
	case StatusExecuted.value:
		return StatusExecuted, nil
	case StatusFailed.value:
		return StatusFailed, nil
	default:
		return Status{}, fmt.Errorf("invalid command status: %s", value)
	}
}

func (s Status) String() string {
	return s.value
}
