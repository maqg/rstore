package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	ImportCmd.Flags().StringVarP(&imageuuid, "imageuuid", "u", "", "imageuuid")
	ImportCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	ImportCmd.Flags().StringVarP(&callbackurl, "callbackurl", "c", "", "callbackurl to async")
}

var ImportCmd = &cobra.Command{

	Use:   "import -filepath xxx -imageuuid xxx -callbackurl xxx",
	Short: "Import image from local to bs.",

	Run: func(cmd *cobra.Command, args []string) {

		if imageuuid != "" {
			fmt.Printf("got imageuuid %s\n", imageuuid)
			return
		}

		cmd.Usage()
	},
}
