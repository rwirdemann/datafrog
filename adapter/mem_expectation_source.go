package adapter

type MemExpectationSource struct {
	expecations []string
}

func NewMemExpectationSource(expecations []string) MemExpectationSource {
	return MemExpectationSource{expecations: expecations}
}

func (es MemExpectationSource) GetAll() []string {
	return es.expecations
}
