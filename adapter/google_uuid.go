package adapter

import "github.com/google/uuid"

type GoogleUUIDProvider struct {
}

func (g GoogleUUIDProvider) NewString() string {
	return uuid.NewString()
}
