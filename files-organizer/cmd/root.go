package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "organizer",
	Short: "Files organizer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		fmt.Println(cmd.UsageString())
	},
}
