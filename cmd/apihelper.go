package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"octlink/rstore/api/v1"
)

func main() {
	descriptors := v1.RouteDescriptorsMap

	data, _ := json.Marshal(descriptors)

	var formated bytes.Buffer
	json.Indent(&formated, data, "", "\t")

	fmt.Printf("\n%s\n", formated)
}
