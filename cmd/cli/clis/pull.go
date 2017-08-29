package clis

import (
	"fmt"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	pullCmd.Flags().StringVarP(&installpath, "installpath", "i", "", "Install path like 'rstore://xxxxxx/xxxxxxxx'")
	pullCmd.Flags().StringVarP(&outpath, "outpath", "o", "./out.qcow2", "Output file path")
	pullCmd.Flags().StringVarP(&address, "address", "a", "localhost:5000", "Rstore Server Address")
}

func getBlob(name string, blobhash string) ([]byte, int, error) {
	url := fmt.Sprintf("http://%s/v1/%s/blobs/%s", address, name, blobhash)
	data, len, err := blobs.HTTPGetBlob(url)
	if err != nil {
		fmt.Printf("got blob from url %s error\n", url)
		return nil, 0, err
	}

	return data, len, nil
}

func pullImage() {

	name, digest := manifest.ParseInstallPath(installpath)
	if name == "" || digest == "" {
		fmt.Printf("name and digest must be specified\n")
		return
	}

	url := fmt.Sprintf("http://%s/v1/%s/manifests/%s", address, name, digest)
	manifest, err := manifest.HTTPGetManifest(url)
	if err != nil {
		fmt.Printf("get manifest from url %s error, manifest may not exist\n", url)
		return
	}

	fmt.Printf("got manifest %s\n", utils.JSON2String(manifest))

	url = fmt.Sprintf("http://%s/v1/%s/blobsmanifest/%s", address, name, manifest.BlobSum)
	blobs, err := blobsmanifest.HTTPGetBlobsManifest(url)
	if err != nil {
		fmt.Printf("got blobsmanifest error from url %s\n", url)
		return
	}
	fmt.Printf("got blobsmanifest %s\n", utils.JSON2String(blobs))

	fd, err := os.Create(outpath)
	if err != nil {
		fmt.Printf("create output file %s error\n", outpath)
		return
	}

	defer fd.Close()

	// Start to downloan blobs
	var total int
	for _, blob := range blobs.Chunks {
		data, len, err := getBlob(name, blob)
		if err != nil {
			fmt.Printf("got blob of %s error\n", blob)
			fmt.Printf("Pull Image of %s error\n", installpath)
			return
		}
		fd.Write(data)
		total += len
	}

	fmt.Printf("pull image of %s to %s length[%d] OK\n", installpath, outpath, total)
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
