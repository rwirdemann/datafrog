package adapter

import "github.com/google/uuid"

// GoogleUUIDProvider generates uuid's via google's uuid package.
type GoogleUUIDProvider struct {
}

// NewString creates a new random UUID and returns it as a string or panics.
func (g GoogleUUIDProvider) NewString() string {
	return uuid.NewString()
}
