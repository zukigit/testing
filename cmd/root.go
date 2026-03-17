package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	tickets "github.com/zukigit/testing/tickets"
)

var tkts []tickets.Ticket

var rootCmd = &cobra.Command{
	Use:   "testing",
	Short: "run all tickets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
