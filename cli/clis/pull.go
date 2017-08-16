package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	PullCmd.Flags().StringVarP(&id, "id", "u", "", "id")
}

var PullCmd = &cobra.Command{

	Use:   "pull",
	Short: "Pull image from rstore.",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("running in pull service\n")

		if id != "" {
			fmt.Printf("got uuid %s\n", id)
			return
		}

		cmd.Usage()
	},
}
