package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	SearchCmd.Flags().StringVarP(&imageuuid, "imageuuid", "u", "", "imageuuid")
}

var SearchCmd = &cobra.Command{

	Use:   "search",
	Short: "Search image by name,account,uuid",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("running in search service\n")

		if imageuuid != "" {
			fmt.Printf("got uuid %s\n", imageuuid)
			return
		}

		cmd.Usage()
	},
}
