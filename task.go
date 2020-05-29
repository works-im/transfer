package transfer

// Source for transfer
type Source interface {
	Reader(query Query) (Packet, error)
}

// Target for transfer
type Target interface {
	Writer(packet Packet) error
}

// Task for transfer
type Task struct {
	config Configuration
	Query  Query
	Source Source
	Target Target
}

// NewTask return task
func NewTask(config Configuration) (task *Task, err error) {

	task = &Task{
		config: config,
		Query:  config.Query,
	}

	sourceOptions := &DatabaseOptions{
		Driver:  config.Source,
		Mapping: config.Mapping,
	}

	if task.Source, err = GenerateSourceTransfer(sourceOptions); err != nil {
		return nil, err
	}

	targetOptions := &DatabaseOptions{
		Driver:  config.Target,
		Mapping: config.Mapping,
	}

	if task.Target, err = GenerateTargetTransfer(targetOptions); err != nil {
		return nil, err
	}

	return task, nil
}

// Run task
func (task *Task) Run() (err error) {

	var (
		result []M
		query  = task.Query
	)

	for {
		if result, err = task.Source.Reader(query); err != nil {
			return err
		}

		if err = task.Target.Writer(result); err != nil {
			return err
		}

		query.Page++

		if len(result) < query.Size {
			break
		}
	}

	return
}
