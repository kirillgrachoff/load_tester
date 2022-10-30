package cmd

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var versionFlag bool

const version = "1.1.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "load_tester",
	Short: "Simple app for load test",
	Long:  "Sends many Get requests to an address",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Println(version)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer cancel()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Printf("error: %s", err)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "show version")
}
