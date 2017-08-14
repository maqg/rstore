package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showVersion bool

func init() {
	RootCmd.AddCommand(HenryCmd)
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show the version and exit")
}

var RootCmd = &cobra.Command{
	Use:   "rstcli",
	Short: "`rstcli`",
	Long:  "`rstcli`",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("V 0.0.1\n")
			return
		}
		cmd.Usage()
	},
}
