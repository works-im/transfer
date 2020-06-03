package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/works-im/transfer"
)

var databaseCmd = &cobra.Command{
	Use:   "database [no options!]",
	Short: "Database transfer",
	Run: func(cmd *cobra.Command, args []string) {

		if len(taskName) > 0 {
			var config transfer.Configuration

			err := viper.UnmarshalKey(taskName, &config)
			if err != nil {
				panic(fmt.Errorf("Fatal error unmarsha config: %s", err))
			}

			task, err := transfer.NewTask(config)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if err = task.Run(); err != nil {
				fmt.Println(err.Error())
				return
			}

		} else {

			var configs []transfer.Configuration
			err := viper.UnmarshalKey("tasks", &configs)
			if err != nil {
				panic(fmt.Errorf("Fatal error unmarsha config: %s", err))
			}

			var wg sync.WaitGroup

			for _, config := range configs {
				fmt.Printf("task config: %#v\n", config)

				wg.Add(1)

				go func(config transfer.Configuration) {
					defer wg.Done()

					task, err := transfer.NewTask(config)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					if err = task.Run(); err != nil {
						fmt.Println(err.Error())
						return
					}

				}(config)
			}

			wg.Wait()
		}
	},
}
