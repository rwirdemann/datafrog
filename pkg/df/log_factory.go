package df

type LogFactory interface {
	Create(filename string) Log
}
