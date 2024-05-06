package ports

type UUIDProvider interface {
	NewString() string
}
