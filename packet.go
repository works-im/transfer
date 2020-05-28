package transfer

// Field with database mapping
// source database field name
// target database filed name
// target_type target database field data type
// TODO: converter transfer script template
type Field struct {
	Source     string `mapstructure:"source"`
	Target     string `mapstructure:"target"`
	TargetType string `mapstructure:"target_type"`
	Converter  string `mapstructure:"converter"`
}

// Mapping for database fields map
type Mapping []Field

// M transfer data type
type M map[string]interface{}

// Packet for transfer data
type Packet []M

// Packager for transfer
type Packager struct{}
