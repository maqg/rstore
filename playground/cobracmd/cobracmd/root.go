package cobracmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showVersion bool

func init() {
	RootCmd.AddCommand(ServeCmd)
	RootCmd.AddCommand(HenryCmd)
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show the version and exit")
}

var RootCmd = &cobra.Command{
	Use:   "registry",
	Short: "`registry`",
	Long:  "`registry`",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("no version specified\n")
			return
		}
		cmd.Usage()
	},
}

func main() {

	fmt.Printf("uuid is %s\n", uuid)

}
