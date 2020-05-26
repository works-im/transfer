package transfer

// Source for transfer
type Source interface {
	Reader(query Query) error
}

// Target for transfer
type Target interface {
	Writer(data []M) error
}

// Task for transfer
type Task struct {
	Source    Source
	Target    Target
	Goroutine int
	Mapping   Mapping

	// Channel
}

// Run task
func (task *Task) Run() error {

	// itr, err := task.Source.Reader()
	// defer itr.Close()
	// for itr.Next() {
	// 	key, value := itr.Value()
	// }
	// task.Source.Reader("", task.Target)

	return nil
}
