package dg

import (
	"fmt"
	"octlink/rstore/utils"
	"testing"
)

func TestDigest(t *testing.T) {
	dd := []byte{1, 2, 4}

	dgst := utils.GetDigest(dd)

	fmt.Printf("got digest of %s\n", dgst)
}
