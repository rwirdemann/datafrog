package ports

// ExpectationSource defines methods to retrieve and remove string based
// expecations from the underlying source, e.g. files. While GetAll returns a
// list of all expectations in its raw format RemoveFirst removes the first
// expectations from the source if it matches the given pattern.
type ExpectationSource interface {
	GetAll() []string
	RemoveFirst(pattern string) error
}
