package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTime2Str(t *testing.T) {
	timeInt := int64(time.Now().Unix())
	s := Time2Str(timeInt)
	fmt.Println(s)
	t.Log("pass")
}
