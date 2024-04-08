package ports

type ExpectationSource interface {
	GetAll() []string
}
