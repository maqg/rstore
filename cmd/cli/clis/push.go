package clis

import (
	"fmt"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/merrors"

	"github.com/spf13/cobra"
)

func init() {
	pushCmd.Flags().StringVarP(&id, "id", "i", "", "id")
	pushCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	pushCmd.Flags().StringVarP(&callbackurl, "callbackurl", "b", "", "callbackurl to async")
	pushCmd.Flags().StringVarP(&address, "address", "a", "localhost:5000", "Rstore Server Address")
}

func pushImage() int {

	if id == "" || filepath == "" {
		fmt.Printf("id or filepath must specified,id:%s,filepath:%s\n",
			id, filepath)
		return merrors.ErrBadParas
	}

	if !utils.IsFileExist(filepath) {
		fmt.Printf("filepath of %s not exist", filepath)
		return merrors.ErrBadParas
	}

	urlPattern := fmt.Sprintf("http://%s/v1/%s/blobs/", address, id)
	hashes, size, err := blobs.HTTPWriteBlobs(filepath, urlPattern)
	if err != nil {
		fmt.Printf("got file hashlist error\n")
		return merrors.ErrCommonErr
	}

	if 1 == 1 {
		return 0
	}

	// write blobs-manifest config
	bm := new(blobsmanifest.BlobsManifest)
	bm.Size = size
	bm.Chunks = hashes
	bm.BlobSum = bm.GetBlobSum()
	err = bm.HTTPWrite()
	if err != nil {
		fmt.Printf("write blobs-manifest error\n")
		return merrors.ErrSystemErr
	}

	// write manifest config
	mid := utils.GetDigestStr(id)
	manifest := new(manifest.Manifest)
	manifest.Name = id
	manifest.ID = mid
	manifest.DiskSize = size
	manifest.VirtualSize = utils.GetVirtualSize(filepath)
	manifest.CreateTime = utils.CurrentTimeStr()
	manifest.BlobSum = bm.BlobSum

	err = manifest.HTTPWrite()
	if err != nil {
		fmt.Printf("Create manifest error[%s]\n", err)
		// TDB,rollback
		return merrors.ErrSystemErr
	}

	if callbackurl != "" {
		callbacking()
	}

	fmt.Printf("Import image OK")

	return merrors.ErrSuccess
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
