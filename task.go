package transfer

// Source for transfer
type Source interface {
	Reader(query Query) (Packet, error)
}

// Target for transfer
type Target interface {
	Writer(packet Packet) error
}

// Configuration task configuration
type Configuration struct {
	Source struct {
		Driver Driver `mapstructure:"driver"`
		Table  string `mapstructure:"table"`
	} `mapstructure:"source"`

	Target struct {
		Driver Driver `mapstructure:"driver"`
		Table  string `mapstructure:"table"`
	} `mapstructure:"target"`

	Mapping Mapping `mapstructure:"mapping"`
	Query   Query   `mapstructure:"query"`
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
		Driver:    config.Source.Driver,
		TableName: config.Source.Table,
		Mapping:   config.Mapping,
	}

	if task.Source, err = GenerateSourceTransfer(sourceOptions); err != nil {
		return nil, err
	}

	targetOptions := &DatabaseOptions{
		Driver:    config.Source.Driver,
		TableName: config.Source.Table,
		Mapping:   config.Mapping,
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
			break
		}

		if err = task.Target.Writer(result); err != nil {
			return err
		}

		query.Page++

		if uint(len(result)) < query.Size {
			break
		}
	}

	return
}
