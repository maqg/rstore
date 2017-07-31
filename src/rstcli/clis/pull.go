package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uuid string

func init() {
	HenryCmd.Flags().StringVarP(&uuid, "uuid", "u", "", "imageuuid")
}

var HenryCmd = &cobra.Command{

	Use:   "pull <config>",
	Short: "`pull` pull image from rstore",
	Long:  "`pull` pull image from rstore.",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("running in pull service\n")

		if uuid != "" {
			fmt.Printf("got uuid %s\n", uuid)
			return
		}

		cmd.Usage()
	},
}
