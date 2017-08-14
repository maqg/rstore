package clis

import "github.com/spf13/cobra"

func init() {
	ImagesCmd.Flags().StringVarP(&imageuuid, "imageuuid", "u", "", "imageuuid")
	ImagesCmd.Flags().StringVarP(&name, "name", "n", "", "Image Name")
	ImagesCmd.Flags().StringVarP(&account, "account", "a", "", "Account Id")
}

var ImagesCmd = &cobra.Command{

	Use:   "images",
	Short: "List images of rstore.",

	Run: func(cmd *cobra.Command, args []string) {

		cmd.Usage()
	},
}
