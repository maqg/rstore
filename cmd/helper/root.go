package helper

import "github.com/spf13/cobra"

var (
	name   string // blob,blob-upload,manifest
	method string // GET,PUT,HEAD,POST,DELETE,PATCH
)

func init() {
	RootCmd.AddCommand(APIHelperCmd)
	RootCmd.AddCommand(MethodHelperCmd)
}

var RootCmd = &cobra.Command{

	Use:   "apihelper",
	Short: "RSTORE API HELPER TOOLS",

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
