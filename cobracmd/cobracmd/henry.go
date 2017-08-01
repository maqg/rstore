package cobracmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uuid string

func init() {
	HenryCmd.Flags().StringVarP(&uuid, "uuid", "u", "", "imageuuid")
}

var HenryCmd = &cobra.Command{

	Use:   "henry <config>",
	Short: "`henry` start henry service",
	Long:  "`henry` start henry services.",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("running in henry service\n")

		if uuid != "" {
			fmt.Printf("got uuid %s\n", uuid)
			return
		}

		cmd.Usage()
	},
}
