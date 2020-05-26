package transfer

import (
	"fmt"
)

// MySQL database transfer
type MySQL struct {
	TableName string
	Fields    []Field
}

// NewMySQL return MySQL transfer
func NewMySQL() *MySQL {
	return &MySQL{}
}

// Reader database
func (mongo *MySQL) Reader(query interface{}, paginator *Pagination) error {

	return nil
}

// Writer data
func (mongo *MySQL) Writer(data []M) error {

	fmt.Printf("%#v\n", data)

	return nil
}
