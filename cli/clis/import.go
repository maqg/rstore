package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	ImportCmd.Flags().StringVarP(&id, "id", "u", "", "id")
	ImportCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	ImportCmd.Flags().StringVarP(&callbackurl, "callbackurl", "c", "", "callbackurl to async")
}

var ImportCmd = &cobra.Command{

	Use:   "import -filepath xxx -id xxx -callbackurl xxx",
	Short: "Import image from local to bs.",

	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			fmt.Printf("got id %s\n", id)
			return
		}

		cmd.Usage()
	},
}
