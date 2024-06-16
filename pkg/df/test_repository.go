package df

type TestRepository interface {
	All() ([]Testcase, error)
	Get(filename string) (Testcase, error)
}
