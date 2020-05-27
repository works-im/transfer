package transfer

// Mapping database field mapping
// source database field name
// target database filed name
// target_type target database field data type
// TODO: converter transfer script template
type Mapping map[string]struct {
	Source     string `mapstructure:"source"`
	Target     string `mapstructure:"target"`
	TargetType string `mapstructure:"target_type"`
	Converter  string `mapstructure:"converter"`
}

// M transfer data type
type M map[string]interface{}

// Packet for transfer data
type Packet []M

// Packager for transfer
type Packager struct{}
