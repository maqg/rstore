package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	pushCmd.Flags().StringVarP(&id, "id", "i", "", "id")
	pushCmd.Flags().StringVarP(&address, "address", "a", "localhost:8000", "Rstore Server Address")
}

var pushCmd = &cobra.Command{

	Use:   "Push",
	Short: "Push image to remote storage.",

	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			fmt.Printf("got uuid %s\n", id)
			return
		}

		cmd.Usage()
	},
}
