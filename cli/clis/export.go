package clis

import (
	"fmt"
	"octlink/mirage/src/utils/merrors"
	"octlink/rstore/configuration"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"

	"github.com/spf13/cobra"
)

func init() {
	exportCmd.Flags().StringVarP(&id, "id", "i", "", "manifest uuid")
	exportCmd.Flags().StringVarP(&name, "name", "n", "", "image uuid")
	exportCmd.Flags().StringVarP(&config, "config", "c", "./config.yml", "Config file for rstore")
	exportCmd.Flags().StringVarP(&outpath, "outpath", "o", "", "output file path of local image")
}

func exportImage() int {

	fmt.Printf("got image id[%s],out[%s],root[%s]\n",
		id, outpath, config)

	if id == "" || outpath == "" || config == "" || name == "" {
		fmt.Printf("id or filepath must specified,id:%s,filepath:%s,rootdir:%s,name:%s\n",
			id, outpath, config, name)
		return merrors.ERR_UNACCP_PARAS
	}

	conf, err := configuration.ResolveConfig(config)
	if err != nil {
		fmt.Printf("parse config %s error\n", config)
		return merrors.ERR_CMD_ERR
	}

	reposDir := conf.RootDirectory + "/" + manifest.ReposDir
	if !utils.IsFileExist(reposDir) {
		fmt.Printf("Directory of %s not exist\n", reposDir)
		return merrors.ERR_UNACCP_PARAS
	}

	manifest := manifest.GetManifest(name, id)
	if manifest == nil {
		fmt.Printf("manifest of %s:%s not exist\n", name, id)
		return merrors.ERR_SEGMENT_NOT_EXIST
	}

	bm := blobsmanifest.GetBlobsManifest(manifest.BlobSum)
	if bm == nil {
		fmt.Printf("blobs manifest of %s not exist\n", manifest.BlobSum)
		return merrors.ERR_SEGMENT_NOT_EXIST
	}

	fmt.Printf(utils.JSON2String(bm))

	return 0
}

var exportCmd = &cobra.Command{

	Use:   "export -outpath xxx -id xxx -r ./",
	Short: "Export image to local from bs.",

	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			exportImage()
			return
		}

		cmd.Usage()
	},
}
