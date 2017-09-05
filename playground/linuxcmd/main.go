package main

import (
	"fmt"
	"octlink/rstore/utils"
)

func main() {
	x := utils.GetVirtualSize("../../../test.qcow2")
	fmt.Printf("got virtual len of %d\n", x)

	y, err := utils.OCTSystem("ls -al")
	if err != nil {
		fmt.Printf("cmd error %s\n", err)
		return
	}

	fmt.Printf("got result of %s\n", y)
}
