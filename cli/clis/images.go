package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	imagesCmd.Flags().StringVarP(&id, "id", "i", "", "id")
	imagesCmd.Flags().StringVarP(&name, "name", "n", "", "Image Name")
	imagesCmd.Flags().StringVarP(&account, "user", "u", "", "Account Id")
	imagesCmd.Flags().StringVarP(&address, "address", "a", "localhost:5000", "Rstore Server Address")
}

var imagesDispHeader = "%-33s%-33s%-16s%-6s%-7s%-7s%-16s\n"

func listImages() {
	fmt.Printf(imagesDispHeader, "Name", "Uuid",
		"Account", "Arch",
		"R-Size", "V-SIZE",
		"CreateTime")
	fmt.Printf(imagesDispHeader, "c4fffb59fc8e40899abd824d654ce416",
		"2264713e80d311e79b5f525400659eb7",
		"admin", "amd64", "500M", "100G", "2017-08-10 22:33:33")
}

func imageDetail() {
	fmt.Printf("\nName: %s\n", "c4fffb59fc8e40899abd824d654ce416")
	fmt.Printf("Uuid: %s\n", "2264713e80d311e79b5f525400659eb7")
	fmt.Printf("Account: %s\n", "admin")
	fmt.Printf("Blobsum: %s\n", "f5755a250b60cb7f555a7536e956f8562ab600188850b9acead1a38c5de42360")
	fmt.Printf("Arch: %s\n", "amd64")
	fmt.Printf("Size: %s\n", "800M")
	fmt.Printf("VirtualSize: %s\n", "100G")
	fmt.Printf("createTime: %s\n", "2017-08-10 22:33:33")
}

var imagesCmd = &cobra.Command{

	Use:   "images",
	Short: "List images of rstore.",

	Run: func(cmd *cobra.Command, args []string) {
		if id != "" {
			imageDetail()
		} else {
			listImages()
		}
	},
}
