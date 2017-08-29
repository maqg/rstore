package clis

import (
	"fmt"
	"octlink/mirage/src/utils/merrors"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"

	"github.com/spf13/cobra"
)

func init() {
	exportCmd.Flags().StringVarP(&blobsum, "blobsum", "s", "", "manifest blob sum")
	exportCmd.Flags().StringVarP(&config, "config", "c", "./config.yml", "Config file for rstore")
	exportCmd.Flags().StringVarP(&outpath, "outpath", "o", "./out.qcow2", "output file path of local image")
}

func exportImage() int {

	if blobsum == "" || outpath == "" || config == "" {
		fmt.Printf("id or filepath must specified,blobsum:%s,out:%s,config:%s\n",
			blobsum, outpath, config)
		return merrors.ERR_UNACCP_PARAS
	}

	conf, err := configuration.ResolveConfig(config)
	if err != nil {
		fmt.Printf("parse config %s error\n", config)
		return merrors.ERR_CMD_ERR
	}

	reposDir := utils.TrimDir(conf.RootDirectory + "/" + manifest.ReposDir)
	if !utils.IsFileExist(reposDir) {
		fmt.Printf("Directory of %s not exist\n", reposDir)
		return merrors.ERR_UNACCP_PARAS
	}

	bm := blobsmanifest.GetBlobsManifest(blobsum)
	if bm == nil {
		fmt.Printf("blobs manifest of %s not exist\n", blobsum)
		return merrors.ERR_SEGMENT_NOT_EXIST
	}

	fmt.Println(utils.JSON2String(bm))

	err = bm.Export(outpath)
	if err != nil {
		fmt.Printf("export image to %s error\n", outpath)
	}

	fmt.Printf("Export image of %s to %s OK\n", blobsum, outpath)

	return 0
}

var exportCmd = &cobra.Command{

	Use:   "export -outpath xxx -id xxx -r ./",
	Short: "Export image to local from bs.",

	Run: func(cmd *cobra.Command, args []string) {

		if blobsum != "" {
			exportImage()
			return
		}

		cmd.Usage()
	},
}
