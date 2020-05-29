package transfer

// Configuration task configuration
type Configuration struct {
	Source Driver `mapstructure:"source"`
	Target Driver `mapstructure:"target"`

	Mapping Mapping `mapstructure:"mapping"`
	Query   Query   `mapstructure:"query"`
}
