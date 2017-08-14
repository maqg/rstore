package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	PushCmd.Flags().StringVarP(&imageuuid, "imageuuid", "u", "", "imageuuid")
}

var PushCmd = &cobra.Command{

	Use:   "push",
	Short: "Push image to remote storage.",

	Run: func(cmd *cobra.Command, args []string) {

		if imageuuid != "" {
			fmt.Printf("got uuid %s\n", imageuuid)
			return
		}

		cmd.Usage()
	},
}
