package transfer

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Driver for database
type Driver struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Field database field
type Field struct {
	Name string
	Type interface{}
}

// Mapping database field mapping
// key: source
// value: target
type Mapping map[string]struct {
	Source     string      `mapstructure:"source"`
	Target     string      `mapstructure:"target"`
	TargetType interface{} `mapstructure:"target_type"`
	Converter  string      `mapstructure:"converter"`
}

// M transfer data type
type M map[string]interface{}

// Pagination for database
type Pagination struct {
	Page int
	Size int
}

// Query database query
type Query struct {
	Q    interface{}
	Page *int
	Size *int
}

// UnmarshalQuery implements the json.Marshaler interface.
func (q Query) UnmarshalQuery(v interface{}) error {

	byteData, err := json.Marshal(q)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(byteData), &v)
}
