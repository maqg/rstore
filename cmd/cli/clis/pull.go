package clis

import (
	"fmt"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/merrors"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	pullCmd.Flags().StringVarP(&installpath, "installpath", "i", "", "Install path like 'rstore://xxxxxx/xxxxxxxx'")
	pullCmd.Flags().StringVarP(&outpath, "outpath", "o", "./out.qcow2", "Output file path")
	pullCmd.Flags().StringVarP(&address, "address", "a", "localhost:5000", "Rstore Server Address")
}

func getBlob(name string, blobhash string) ([]byte, int, error) {

	if blobhash == config.ZeroDataDigest8M {
		// zero data no need fetch from remote
		fmt.Printf("zero data no need fetch from remote\n")
		return config.ZeroData8M, configuration.BlobSize, nil
	}

	url := fmt.Sprintf("http://%s/v1/%s/blobs/%s", address, name, blobhash)
	data, len, err := blobs.HTTPGetBlob(url)
	if err != nil {
		fmt.Printf("got blob from url %s error\n", url)
		return nil, 0, err
	}

	return data, len, nil
}

func pullImage() int {

	name, digest := manifest.ParseInstallPath(installpath)
	if name == "" || digest == "" {
		fmt.Printf("name and digest must be specified\n")
		return merrors.ErrBadParas
	}

	url := fmt.Sprintf("http://%s/v1/%s/manifests/%s", address, name, digest)
	manifest, err := manifest.HTTPGetManifest(url)
	if err != nil {
		fmt.Printf("get manifest from url %s error, manifest may not exist\n", url)
		return merrors.ErrSegmentNotExist
	}

	fmt.Printf("got manifest %s\n", utils.JSON2String(manifest))

	url = fmt.Sprintf("http://%s/v1/%s/blobsmanifest/%s", address, name, manifest.BlobSum)
	blobs, err := blobsmanifest.HTTPGetBlobsManifest(url)
	if err != nil {
		fmt.Printf("got blobsmanifest error from url %s\n", url)
		return merrors.ErrUserNotExist
	}
	fmt.Printf("got blobsmanifest %s\n", utils.JSON2String(blobs))

	fd, err := os.Create(outpath)
	if err != nil {
		fmt.Printf("create output file %s error\n", outpath)
		return merrors.ErrCmdErr
	}

	defer fd.Close()

	// Start to downloan blobs
	var total int64
	for _, blob := range blobs.Chunks {
		data, len, err := getBlob(name, blob)
		if err != nil {
			fmt.Printf("got blob of %s error\n", blob)
			fmt.Printf("Pull Image of %s error\n", installpath)
			return merrors.ErrSyscallErr
		}
		fd.Write(data)
		total += len
	}

	if total != blobs.Size {
		fmt.Printf("size fetched %lld not match size %lld\n", total, blobs.Size)
		return merrors.ErrCommonErr
	}

	fmt.Printf("pull image of %s to %s length[%d] OK\n", installpath, outpath, total)

	return merrors.ErrSuccess
}

var pullCmd = &cobra.Command{

	Use:   "pull",
	Short: "Pull image from rstore.",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("running in pull service\n")

		if installpath != "" {
			os.Exit(pullImage())
		}

		cmd.Usage()

		os.Exit(merrors.ErrBadParas)
	},
}
