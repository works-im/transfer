package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	assert := assert.New(t)

	args := &DatabaseOptions{
		Driver: Driver{
			Host:     "localhost",
			Port:     "27017",
			Username: "root",
			Password: "",
			Database: "qvm-order",
		},
		TableName: "atomic",
	}

	db, err := NewMongoDB(args)

	assert.NotNil(db)
	assert.Nil(err)

	query := Query{
		Q: []M{
			{
				"$project": M{
					"_id":           1,
					"uid":           1,
					"resource_type": "$atomic.atomic_class_info.resource_type",
				},
			},
		},
	}

	result, err := db.Reader(query)
	assert.Nil(err)
	assert.NotEmpty(result)
}
