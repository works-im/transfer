package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/works-im/transfer/database"
)

func TestTask(t *testing.T) {
	assert := assert.New(t)

	config := Configuration{
		Source: database.Driver{
			Driver:   "mongodb",
			Host:     "localhost",
			Port:     27017,
			Username: "root",
			Password: "",
			Database: "qvm-order",
			Table:    "atomic",
		},

		Target: database.Driver{
			Driver:   "mysql",
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "root",
			Database: "qvm-order",
			Table:    "atomics",
		},

		Mapping: []*database.Field{
			{
				Source:     "_id",
				Target:     "id",
				TargetType: "string",
			},
			{
				Source:     "uid",
				Target:     "uid",
				TargetType: "uint",
			},
			{
				Source:     "atomic.atomic_class_info.resource_type",
				Target:     "resource_type",
				TargetType: "string",
			},
		},

		Query: database.Query{
			Q: []database.M{
				{
					"$match": database.M{
						"uid": database.M{"$gt": 0},
					},
				},
			},
			Page: 1,
			Size: 100,
		},
	}

	task, err := NewTask(config)
	assert.Nil(err)
	assert.NotNil(task)

	err = task.Run()
	assert.Nil(err)
}
