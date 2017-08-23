package main

import (
	"fmt"
	"io"
	"octlink/rstore/configuration"
	"octlink/rstore/utils"
	"os"
)

func main() {
	filepath := "G:\\test.qcow2"

	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("file of %s not exist\n", filepath)
		return
	}

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

		fmt.Printf("ddd[%s]\n", dgst)
	}
}
