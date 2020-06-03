package transfer

import (
	"github.com/works-im/transfer/database"
)

// Configuration task configuration
type Configuration struct {
	Source database.Driver `mapstructure:"source"`
	Target database.Driver `mapstructure:"target"`

	Mapping database.Mapping `mapstructure:"mapping"`
	Query   database.Query   `mapstructure:"query"`
}
