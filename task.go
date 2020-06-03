package transfer

import (
	"transfer/database"
	"transfer/mongo"
	"transfer/mysql"
)

// Source for transfer
type Source interface {
	Reader(query database.Query) (database.Packet, error)
}

// Target for transfer
type Target interface {
	Writer(packet database.Packet) error
}

// Task for transfer
type Task struct {
	config Configuration
	Query  database.Query
	Source Source
	Target Target
}

// NewTask return task
func NewTask(config Configuration) (task *Task, err error) {

	task = &Task{
		config: config,
		Query:  config.Query,
	}

	for _, field := range config.Mapping {
		if len(field.Target) == 0 {
			field.Target = field.Source
		}
	}

	sourceOptions := &database.Options{
		Driver:  config.Source,
		Mapping: config.Mapping,
	}

	if task.Source, err = GenerateSourceTransfer(sourceOptions); err != nil {
		return nil, err
	}

	targetOptions := &database.Options{
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
		result []database.M
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

// GenerateSourceTransfer return source transfer
func GenerateSourceTransfer(args *database.Options) (source Source, err error) {

	switch args.Driver.Driver {
	case "mongodb":
		source, err = mongo.NewMongoDB(args)
	case "mysql":
		source, err = mysql.NewMySQL(args)
	}

	return
}

// GenerateTargetTransfer return target transfer
func GenerateTargetTransfer(args *database.Options) (target Target, err error) {

	switch args.Driver.Driver {
	case "mongodb":
		target, err = mongo.NewMongoDB(args)
	case "mysql":
		target, err = mysql.NewMySQL(args)
	}

	return
}
