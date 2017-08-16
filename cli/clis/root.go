package clis

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	showVersion bool
	id          string
	filepath    string // image file path
	callbackurl string
	taskid      string
	name        string // image name
	account     string // account uuid
	cachedir    string // temp dir for image files caching
	installpath string // where to store image of local
	storepath   string // store path like rstore://uuid:manifest
	address     string // address like 10.10.0.100:8000
)

func init() {
	RootCmd.AddCommand(PullCmd)
	RootCmd.AddCommand(PushCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(ImagesCmd)
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show the version and exit")
}

var RootCmd = &cobra.Command{

	Use:   "rstcli",
	Short: "RSTORE CLI TOOLS",

	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("V 0.0.1\n")
			return
		}
		cmd.Usage()
	},
}
