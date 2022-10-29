package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/kirillgrachoff/load_tester/pkg/net/multi_get"
)

var count int

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load host... [ -c count ]",
	Short: "Start load test",
	Long:  `Start load test with given parameters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("sources not specified")
		}

		client := multi_get.NewClient(count, args)
		return client.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().IntVarP(&count, "count", "c", 10, "parallel GETs count")
}
