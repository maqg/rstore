package helper

import (
	"encoding/json"
	"fmt"
	"octlink/rstore/api/v1"

	"github.com/spf13/cobra"
)

func init() {
	MethodHelperCmd.Flags().StringVarP(&name, "name", "n", "", "name")
	MethodHelperCmd.Flags().StringVarP(&method, "method", "m", "", "METHOD:GET,PUT,POST,PATCH,HEAD,DELETE")
}

func printMessage(obj interface{}) {
	data, _ := json.MarshalIndent(obj, "", "    ")
	fmt.Printf("%s\n", data)
}

func helper() {

	descriptors := v1.RouteDescriptorsMap
	for _, desc := range descriptors {
		if desc.Name == name {
			if method != "" {
				for _, m := range desc.Methods {
					if m.Method == method {
						printMessage(m)
						return
					}
				}
				fmt.Printf("No Such Method %s in API %s\n", method, name)
				break
			} else {
				printMessage(desc)
				return
			}
		}
	}

	fmt.Printf("No Such CMD of %s\n", name)
}

var MethodHelperCmd = &cobra.Command{

	Use:   "method -m xxx -n xxx",
	Short: "Get Helper Message of API",

	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			cmd.Usage()
			return
		}
		helper()
	},
}
