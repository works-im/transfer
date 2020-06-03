package database

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Driver for database
type Driver struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Table    string `mapstructure:"table"`
}

// FieldMeta database table schema
type FieldMeta struct {
	Field   string
	Type    string
	Null    interface{}
	Key     interface{}
	Default interface{}
	Extra   interface{}
}

// Schema for database table
type Schema []FieldMeta

// FieldMap return FieldMeta map
func (schema Schema) FieldMap() map[string]FieldMeta {
	var m = map[string]FieldMeta{}
	for _, f := range schema {
		m[f.Field] = f
	}
	return m
}

// Query database query
type Query struct {
	Q    interface{} `mapstructure:"q"`
	Page int         `mapstructure:"page"`
	Size int         `mapstructure:"size"`
}

// UnmarshalQuery implements the json.Marshaler interface.
func (q Query) UnmarshalQuery(v interface{}) error {

	byteData, err := json.Marshal(q.Q)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(byteData), &v)
}

// Options for transfer
type Options struct {
	Driver  Driver
	Mapping Mapping
}
