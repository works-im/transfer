package transfer

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
	SourceType interface{} `mapstructure:"source_type"`
	Target     string      `mapstructure:"target"`
	TargetType interface{} `mapstructure:"target_type"`
	Converter  string      `mapstructure:"converter"`
}
