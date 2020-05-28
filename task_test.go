package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	assert := assert.New(t)

	config := Configuration{
		Source: Driver{
			Driver:   "mongodb",
			Host:     "localhost",
			Port:     27017,
			Username: "root",
			Password: "",
			Database: "qvm-order",
			Table:    "atomic",
		},

		Target: Driver{
			Driver:   "mysql",
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "root",
			Database: "qvm-order",
			Table:    "atomics",
		},

		Mapping: []Field{
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

		Query: Query{
			Q: []M{
				{
					"$match": M{
						"uid": M{"$gt": 0},
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
