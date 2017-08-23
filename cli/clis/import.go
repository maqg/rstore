package clis

import (
	"fmt"
	"octlink/mirage/src/utils"

	"github.com/spf13/cobra"
)

func init() {
	ImportCmd.Flags().StringVarP(&id, "id", "i", "", "id")
	ImportCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	ImportCmd.Flags().StringVarP(&callbackurl, "callbackurl", "b", "", "callbackurl to async")
}

func callbacking() {
	fmt.Printf("callbackurl of %s called\n", callbackurl)
}

func splitFile(filepath string) {
	
}

func checkParas() bool {
	if id == "" || filepath == "" {
		fmt.Printf("id or filepath must specified,id:%s,filepath:%s\n", id, filepath)
		return false
	}

	if !utils.IsFileExist(filepath) {
		fmt.Printf("filepath of %s not exist", filepath)
		return false
	}

	return true
}

func importImage() {

	fmt.Printf("got image id[%s],filepath[%s],callbackurl[%s]\n",
		id, filepath, callbackurl)

	if !checkParas() {
		fmt.Printf("check input paras failed\n")
		return
	}

	if callbackurl != "" {
		callbacking()
	}

	fmt.Printf("Import image OK")
}

var ImportCmd = &cobra.Command{

	Use:   "import -filepath xxx -id xxx -callbackurl xxx",
	Short: "Import image from local to bs.",

	Run: func(cmd *cobra.Command, args []string) {

		if filepath != "" {
			importImage()
			return
		}

		cmd.Usage()
	},
}
