package df

type TestRepository interface {
	All() ([]Testcase, error)
	Get(filename string) (Testcase, error)
	Exists(filename string) bool
}
