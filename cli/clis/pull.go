package clis

import (
	"fmt"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	pullCmd.Flags().StringVarP(&installpath, "installpath", "i", "", "Install path like 'rstore://xxxxxx/xxxxxxxx'")
	pullCmd.Flags().StringVarP(&outpath, "outpath", "o", "./out.qcow2", "Output file path")
	pullCmd.Flags().StringVarP(&address, "address", "a", "localhost:5000", "Rstore Server Address")
}

func pullImage() {
	segs := strings.Split(installpath, "/")
	len := len(segs)
	name := segs[len-2]
	digest := segs[len-1]

	url := fmt.Sprintf("http://%s/v1/%s/manifests/%s", address, name, digest)
	manifest, err := manifest.HTTPGetManifest(url)
	if err != nil {
		fmt.Printf("get manifest from url %s error \n", url)
		return
	}

	fmt.Printf("got manifest %s\n", utils.JSON2String(manifest))
}

var pullCmd = &cobra.Command{

	Use:   "pull",
	Short: "Pull image from rstore.",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("running in pull service\n")

		if installpath != "" {
			pullImage()
			return
		}

		cmd.Usage()
	},
}
