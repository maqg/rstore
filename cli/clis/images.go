package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	ImagesCmd.Flags().StringVarP(&imageuuid, "imageuuid", "u", "", "imageuuid")
	ImagesCmd.Flags().StringVarP(&name, "name", "n", "", "Image Name")
	ImagesCmd.Flags().StringVarP(&account, "account", "a", "", "Account Id")
}

var IMAGES_DISP_HEADER = "  %-16s%-37s%-16s%-16s\n"

func listImages() {
	fmt.Printf(IMAGES_DISP_HEADER, "Name", "Uuid", "Account", "LastSync")
	fmt.Printf(IMAGES_DISP_HEADER, "Image1", "2264713e-80d3-11e7-9b5f-525400659eb7", "admin", "2017-08-10 22:33:33")
}

func imageDetail() {
	fmt.Printf("\nName: %s\n", "Image1")
	fmt.Printf("Uuid: %s\n", "2264713e-80d3-11e7-9b5f-525400659eb7")
	fmt.Printf("Account: %s\n", "admin")
	fmt.Printf("Size: %s\n", "800M")
	fmt.Printf("VirtualSize: %s\n", "100G")
	fmt.Printf("LastSync: %s\n", "2017-08-10 22:33:33")
	fmt.Printf("createTime: %s\n", "2017-08-10 22:33:33")
	fmt.Printf("storePath: %s\n", "rstore://2264713e-80d3-11e7-9b5f-525400659eb7/2264713e-80d3-11e7-9b5f-525400659eb7kkkk")
	fmt.Printf("Desc: %s\n\n", "fadfdasfadsf")
}

var ImagesCmd = &cobra.Command{

	Use:   "images",
	Short: "List images of rstore.",

	Run: func(cmd *cobra.Command, args []string) {
		if imageuuid != "" {
			imageDetail()
		} else {
			listImages()
		}
	},
}
