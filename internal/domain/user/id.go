package user

import "github.com/google/uuid"

type UserID uuid.UUID

func NewUserID(id string) (UserID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), err
	}

	return UserID(parsed), nil
}

func (u UserID) String() string {
	return uuid.UUID(u).String()
}
