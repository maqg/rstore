package clis

import (
	"encoding/json"
	"fmt"
	"io"
	"octlink/mirage/src/utils/merrors"
	"octlink/rstore/configuration"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/modules/revision"
	"octlink/rstore/utils"
	"octlink/rstore/utils/uuid"

	"os"

	"github.com/spf13/cobra"
)

func init() {
	ImportCmd.Flags().StringVarP(&id, "id", "i", "", "id")
	ImportCmd.Flags().StringVarP(&rootdirectory, "rootdirectory", "r", "", "RootDirectory of RSTORE")
	ImportCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file path of local image")
	ImportCmd.Flags().StringVarP(&callbackurl, "callbackurl", "b", "", "callbackurl to async")
}

func callbacking() {
	fmt.Printf("callbackurl of %s called\n", callbackurl)
}

func splitFile(filepath string) ([]string, error) {

	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("file of %s not exist\n", filepath)
		return nil, err
	}

	hashList := make([]string, 0)

	defer f.Close()

	for {
		buffer := make([]byte, configuration.BLOB_SIZE)
		n, err := f.Read(buffer)
		if err == io.EOF {
			fmt.Printf("reached end of file[%d]\n", n)
			break
		}

		if err != nil {
			fmt.Printf("read file error %s", err)
		}

		fmt.Printf("got size of %d\n", n)

		dgst := utils.GetDigest(buffer)

		hashList = append(hashList, dgst)
	}

	return hashList, nil
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

	fmt.Printf("got image id[%s],filepath[%s],callbackurl[%s]\n",
		id, filepath, callbackurl)

	if !checkParas() {
		fmt.Printf("check input paras failed\n")
		return merrors.ERR_UNACCP_PARAS
	}

	reposDir := rootdirectory + "/" + manifest.REPOS_DIR
	if !utils.IsFileExist(reposDir) {
		fmt.Printf("Directory of %s not exist\n", reposDir)
		return merrors.ERR_UNACCP_PARAS
	}

	mid := uuid.Generate().Simple()
	revisionDir := rootdirectory + fmt.Sprintf(manifest.REVISIONS_DIR_PROTO, id, mid)
	fmt.Printf("got new manifest dir %s\n", revisionDir)
	utils.CreateDir(revisionDir)

	revision := new(revision.Revision)
	revision.Name = id
	revision.Id = mid

	hashes, err := splitFile(filepath)
	if err != nil {
		fmt.Printf("got file hashlist error\n")
		return merrors.ERR_COMMON_ERR
	}

	if callbackurl != "" {
		callbacking()
	}

	d, _ := json.MarshalIndent(hashes, "", "  ")
	fmt.Printf("%s\n", string(d))

	fmt.Printf("Import image OK")

	return 0
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
