package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/kirillgrachoff/load_tester/pkg/net/multi_get"
)

var count int
var sources []string

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Start load test",
	Long:  `Start load test with given parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(sources) == 0 {
			log.Fatalln("sources not specified")
		}
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
		defer cancel()

		client := multi_get.NewClient(count, sources)
		if err := client.Run(ctx); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().IntVarP(&count, "count", "c", 1, "parallel GETs count")
	loadCmd.Flags().StringArrayVarP(&sources, "sources", "s", []string{}, "specify sources to GET")
}
