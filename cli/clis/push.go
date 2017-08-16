package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	PushCmd.Flags().StringVarP(&id, "id", "i", "", "id")
}

var PushCmd = &cobra.Command{

	Use:   "push",
	Short: "Push image to remote storage.",

	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			fmt.Printf("got uuid %s\n", id)
			return
		}

		cmd.Usage()
	},
}
