package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	assert := assert.New(t)

	args := &MongoOptions{
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

	query := []M{
		{
			"$project": M{
				"_id":           1,
				"uid":           1,
				"resource_type": "$atomic.atomic_class_info.resource_type",
			},
		},
	}

	paginator := &Pagination{}

	result, err := db.Reader(query, paginator)
	assert.Nil(err)
	assert.NotEmpty(result)
}
