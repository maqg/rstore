package clis

import (
	"fmt"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/image"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/merrors"

	"github.com/spf13/cobra"
)

var conf *configuration.Configuration

func init() {
	importCmd.Flags().StringVarP(&id, "id", "i", "", "Image UUID")
	importCmd.Flags().StringVarP(&configfile, "configfile", "c", "./config.yml", "Config of RSTORE")
	importCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	importCmd.Flags().StringVarP(&callbackurl, "callbackurl", "b", "", "callbackurl to async")
}

func callbacking() {
	fmt.Printf("callbackurl of %s called\n", callbackurl)
}

func checkParas() bool {

	if !utils.IsFileExist(filepath) {
		fmt.Printf("filepath of %s not exist", filepath)
		return false
	}

	return true
}

func importImage() int {

	fmt.Printf("got image id[%s],filepath[%s],config[%s],callbackurl[%s]\n",
		id, filepath, configfile, callbackurl)

	if !checkParas() {
		fmt.Printf("check input paras failed\n")
		return merrors.ErrBadParas
	}

	c, err := configuration.ResolveConfig(configfile)
	if err != nil {
		fmt.Printf("parse config %s error\n", configfile)
		return merrors.ErrCmdErr
	}
	conf = c

	fmt.Printf("with config of HugeBlob %v\n", conf.HugeBlob)

	image.ReloadImages()

	reposDir := utils.TrimDir(conf.RootDirectory + "/" + manifest.ReposDir)
	if !utils.IsFileExist(reposDir) {
		fmt.Printf("Directory of %s not exist\n", reposDir)
		return merrors.ErrBadParas
	}

	var hashes []string
	var size int64

	if conf.HugeBlob {
		hashes, size, err = blobs.ImportHugeBlob(filepath)
	} else {
		hashes, size, err = blobs.ImportBlobs(filepath)
	}
	if err != nil {
		fmt.Printf("got file hashlist error\n")
		return merrors.ErrCommonErr
	}

	if size != utils.GetFileSize(filepath) {
		fmt.Printf("filelen not match %d\n", size)
		return merrors.ErrCmdErr
	}

	// write blobs-manifest config
	bm := new(blobsmanifest.BlobsManifest)
	bm.Size = size
	bm.Chunks = hashes
	bm.BlobSum = bm.GetBlobSum()
	err = bm.Write()
	if err != nil {
		fmt.Printf("write blobs-manifest error\n")
		return merrors.ErrSystemErr
	}

	// for hugeblob supporting, copy file to blobs
	if configuration.HugeBlob() {
		dstFileDir := configuration.RootDirectory() + manifest.ManifestDir + "/" + bm.BlobSum
		utils.CreateDir(dstFileDir)

		dstFilePath := dstFileDir + "/" + "image"
		fmt.Printf("for huge blob mode, copy %s to %s\n", filepath, dstFilePath)
		utils.CopyFile(filepath, dstFilePath)
	}

	// write manifest config
	manifest := new(manifest.Manifest)
	manifest.Name = id
	manifest.DiskSize = size
	manifest.VirtualSize = utils.GetVirtualSize(filepath)
	manifest.CreateTime = utils.CurrentTimeStr()
	manifest.BlobSum = bm.BlobSum

	err = manifest.Write()
	if err != nil {
		fmt.Printf("Create manifest error[%s]\n", err)
		// TDB,rollback
		return merrors.ErrSystemErr
	}
	err = image.UpdateImageCallback(manifest.Name, manifest.DiskSize, manifest.VirtualSize,
		manifest.BlobSum, config.ImageStatusReady)
	if err != nil {
		fmt.Printf("update image info %s for image %s error, and manifest created OK\n",
			manifest.Name, id)
	}

	if callbackurl != "" {
		callbacking()
	}

	utils.SendUserSignal("rstore")

	fmt.Printf("Import image OK\n")

	return 0
}

var importCmd = &cobra.Command{

	Use:   "import -filepath xxx -id xxx -callbackurl xxx",
	Short: "Import image from local to bs.",

	Run: func(cmd *cobra.Command, args []string) {

		if id == "" || filepath == "" || configfile == "" {
			fmt.Printf("id or filepath must specified,id:%s,filepath:%s,config:%s\n",
				id, filepath, configfile)
			cmd.Usage()
			return
		}

		importImage()
	},
}
