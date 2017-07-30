package cobracmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "start <config>",
	Short: "`start` start registry service",
	Long:  "`serve` start registry services.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("running registry service\n")
	},
}
