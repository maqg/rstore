package octmysql

import (
	"octlink/mirage/src/utils/octlog"
	"testing"
)

func TestTest(t *testing.T) {
	t.Log("test OK")
}

func init() {
	InitLog(octlog.DEBUG_LEVEL)
	octlog.InitDebugConfig(octlog.DEBUG_LEVEL)
}
