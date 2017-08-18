package helper

import (
	"fmt"
	"octlink/rstore/api/v1"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	APIHelperCmd.Flags().StringVarP(&name, "name", "n", "", "name")
}

func parseMethod(route v1.RouteDescriptor) string {

	methods := make([]string, 0)

	for _, method := range route.Methods {
		methods = append(methods, method.Method)
	}

	return strings.Join(methods, ",")
}

func api_helper() {

	descriptors := v1.RouteDescriptorsMap
	for _, desc := range descriptors {
		if name != "" {
			if desc.Name == name {
				fmt.Printf("%s(%s): %s\n",
					desc.Name, desc.PathSimple, parseMethod(desc))
				break
			}
		} else {
			fmt.Printf("%s(%s): %s\n",
				desc.Name, desc.PathSimple, parseMethod(desc))
		}
	}
}

var APIHelperCmd = &cobra.Command{

	Use:   "api -n xxx",
	Short: "Get Helper Message of API",

	Run: func(cmd *cobra.Command, args []string) {
		api_helper()
	},
}
