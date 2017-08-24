package clis

import (
	"fmt"
	"octlink/mirage/src/utils/merrors"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/uuid"

	"github.com/spf13/cobra"
)

func init() {
	importCmd.Flags().StringVarP(&id, "id", "i", "", "id")
	importCmd.Flags().StringVarP(&rootdirectory, "rootdirectory", "r", "", "RootDirectory of RSTORE")
	importCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	importCmd.Flags().StringVarP(&callbackurl, "callbackurl", "b", "", "callbackurl to async")
}

func callbacking() {
	fmt.Printf("callbackurl of %s called\n", callbackurl)
}

func checkParas() bool {
	if id == "" || filepath == "" || rootdirectory == "" {
		fmt.Printf("id or filepath must specified,id:%s,filepath:%s,rootdir:%s\n",
			id, filepath, rootdirectory)
		return false
	}

	if !utils.IsFileExist(filepath) {
		fmt.Printf("filepath of %s not exist", filepath)
		return false
	}

	return true
}

func importImage() int {

	fmt.Printf("got image id[%s],filepath[%s],root[%s],callbackurl[%s]\n",
		id, filepath, rootdirectory, callbackurl)

	if !checkParas() {
		fmt.Printf("check input paras failed\n")
		return merrors.ERR_UNACCP_PARAS
	}

	reposDir := rootdirectory + "/" + manifest.ReposDir
	if !utils.IsFileExist(reposDir) {
		fmt.Printf("Directory of %s not exist\n", reposDir)
		return merrors.ERR_UNACCP_PARAS
	}

	_, blobSum, err := blobs.WriteBlobs(filepath, rootdirectory)
	if err != nil {
		fmt.Printf("got file hashlist error\n")
		return merrors.ERR_COMMON_ERR
	}

	mid := uuid.Generate().Simple()
	manifestDir := rootdirectory + fmt.Sprintf(manifest.ManifestDirProto, id)

	manifest := new(manifest.Manifest)
	manifest.Name = id
	manifest.ID = mid
	manifest.Path = manifestDir
	manifest.CreateTime = utils.CurrentTimeStr()
	manifest.BlobSum = blobSum

	err = manifest.Write()
	if err != nil {
		fmt.Printf("Create manifest error[%s]\n", err)
		return merrors.ERR_SYSTEM_ERR
	}

	if callbackurl != "" {
		callbacking()
	}

	fmt.Printf("Import image OK")

	return 0
}

var importCmd = &cobra.Command{

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
