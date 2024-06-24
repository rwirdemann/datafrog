package df

type TestRepository interface {
	All() ([]Testcase, error)
	Get(testname string) (Testcase, error)
	Exists(filename string) bool
	Write(testname string, testcase Testcase) error
	Delete(testname string) error
}
