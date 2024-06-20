package df

type TestRepository interface {
	All() ([]Testcase, error)
	Get(testname string) (Testcase, error)
	Exists(filename string) bool
}
