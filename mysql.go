package transfer

import (
	"fmt"
)

// MySQL database transfer
type MySQL struct {
	Driver    Driver
	TableName string
	Mapping   Mapping
}

// NewMySQL return MySQL transfer
func NewMySQL(args *DatabaseOptions) (*MySQL, error) {
	return &MySQL{}, nil
}

// Reader database
func (mongo *MySQL) Reader(query Query) (packet Packet, err error) {

	return nil, nil
}

// Writer data
func (mongo *MySQL) Writer(packet Packet) error {

	fmt.Printf("%#v\n", packet)

	return nil
}
