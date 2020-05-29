package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"transfer"
)

var databaseCmd = &cobra.Command{
	Use:   "database [no options!]",
	Short: "Database transfer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside subCmd Run with args: %v\n", args)

		var config transfer.Configuration
		err := viper.UnmarshalKey("task", &config)
		if err != nil {
			panic(fmt.Errorf("Fatal error unmarsha config: %s", err))
		}

		fmt.Printf("config: %#v\n", config)

		task, err := transfer.NewTask(config)
		if err != nil {
			panic(err.Error())
		}

		err = task.Run()
		if err != nil {
			panic(err.Error())
		}
	},
}
