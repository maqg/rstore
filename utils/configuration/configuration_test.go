package configuration

import (
	"fmt"
	"testing"
)

var configYaml = `
version: 0.1
debuglevel: 5
loglevel: 5
http:
  addr: :5000
`

func TestParseConfig(t *testing.T) {

	config, err := ParseConfig([]byte(configYaml))
	if err != nil {
		fmt.Printf("parse config error\n")
		return
	}

	fmt.Printf("version %s\n", config.Version)
	fmt.Printf("debuglevel %s\n", config.Debuglevel)
	fmt.Printf("loglevel %s\n", config.Loglevel)
}
