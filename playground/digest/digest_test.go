package dg

import (
	"fmt"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"testing"
)

func TestDigest(t *testing.T) {

	a := make([]byte, configuration.BlobSize/2)

	digest := utils.GetDigest(a)

	fmt.Printf("got digest %s\n", digest)
}
