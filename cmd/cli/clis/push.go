package clis

import (
	"fmt"
	"octlink/rstore/api/v1"
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

	if !utils.IsFileExist(filepath) {
		fmt.Printf("filepath of %s not exist", filepath)
		return merrors.ErrBadParas
	}

	urlPattern := fmt.Sprintf(v1.APIURLFormatBlobUpload, address, id)
	hashes, size, err := blobs.HTTPWriteBlobs(filepath, urlPattern)
	if err != nil {
		fmt.Printf("got file hashlist error\n")
		return merrors.ErrCommonErr
	}

	// write blobs-manifest config
	bm := new(blobsmanifest.BlobsManifest)
	bm.Size = size
	bm.Chunks = hashes
	bm.BlobSum = bm.GetBlobSum()
	err = bm.HTTPWrite(fmt.Sprintf(v1.APIURLFormatBlobsManifest, address, id, bm.BlobSum))
	if err != nil {
		fmt.Printf("write blobs-manifest error\n")
		// TBD remove all posted blobs
		return merrors.ErrSystemErr
	}

	// write manifest config
	manifest := new(manifest.Manifest)
	manifest.Name = id
	manifest.DiskSize = size
	manifest.VirtualSize = utils.GetVirtualSize(filepath)
	manifest.CreateTime = utils.CurrentTimeStr()
	manifest.BlobSum = bm.BlobSum

	err = manifest.HTTPWrite(fmt.Sprintf(v1.APIURLFormatManifests, address, id, bm.BlobSum))
	if err != nil {
		fmt.Printf("Create manifest error[%s]\n", err)
		// TDB,rollback
		return merrors.ErrSystemErr
	}

	if callbackurl != "" {
		callbacking()
	}

	fmt.Printf("Import image %s to %s OK\n", filepath, id)

	return merrors.ErrSuccess
}

var pushCmd = &cobra.Command{

	Use:   "push",
	Short: "push image to remote storage.",

	Run: func(cmd *cobra.Command, args []string) {

		if id != "" && filepath != "" {
			pushImage()
			return
		}
		fmt.Printf("id and filepath must specified,id:%s,filepath:%s\n",
			id, filepath)

		cmd.Usage()
	},
}
