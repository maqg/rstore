package dg

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"octlink/rstore/digest"
	"testing"
)

func TestDigest(t *testing.T) {
	dd := []byte{1, 2, 3}
	dgst, _ := digest.FromReader(bytes.NewReader(dd))
	fmt.Printf("ddd[%s]\n", dgst)

	hash := sha256.New()
	hash.Write(dd)
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	fmt.Printf("ddd[%s]\n", mdStr)
}
