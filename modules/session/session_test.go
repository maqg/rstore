package session

import (
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/octmysql"
	"testing"
)

var DB *octmysql.OctMysql

func TestFindSession(t *testing.T) {
	sess := FindSession(DB, "00000000000000000000000000000000")
	if sess == nil {
		t.Log("Fail")
		return
	}
	t.Log("pass")
}

func TestNewSession(t *testing.T) {
	sess := NewSession(DB, "TestUserId", "Test")
	if sess == nil {
		t.Log("Fail")
		return
	}
	t.Log("pass")
}

func TestClearSession(t *testing.T) {
	ClearSession(DB, "TestUserId")
	t.Log("pass")
}

func init() {

	DB = new(octmysql.OctMysql)

	InitLog(octlog.DEBUG_LEVEL)
	octlog.InitDebugConfig(octlog.DEBUG_LEVEL)
}
