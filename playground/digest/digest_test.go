package dg

import (
	"fmt"
	"testing"
)

func TestDigest(t *testing.T) {
	a := []string{"a", "b", "c"}
	fmt.Printf("len of a is %d\n", len(a))
	fmt.Println(a[0:0])
	/*
		dd := []byte{1, 2, 4}

		dgst := utils.GetDigest(dd)

		fmt.Printf("got digest of %s\n", dgst)
	*/
}
