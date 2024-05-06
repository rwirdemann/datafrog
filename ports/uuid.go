package ports

// UUIDProvider abstracts the generation of UUIDs.
type UUIDProvider interface {
	NewString() string
}
