package transfer

// Source for transfer
type Source interface {
	Reader() error
}

// Target for transfer
type Target interface {
	Writer() error
}

// Task for transfer
type Task struct {
	Source    Source
	Target    Target
	Goroutine int
	Mapping   Mapping
}
