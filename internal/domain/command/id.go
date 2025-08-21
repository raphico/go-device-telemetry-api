package command

import "github.com/google/uuid"

type CommandID uuid.UUID

func NewCommandID(id string) (CommandID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return CommandID(uuid.Nil), err
	}

	return CommandID(parsed), nil
}

func (c CommandID) String() string {
	return uuid.UUID(c).String()
}
