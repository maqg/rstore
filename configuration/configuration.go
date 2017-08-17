package configuration

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const (
	TB_USER         = "tb_user"
	TB_RELUSERGROUP = "tb_relusergroup"
	TB_USERGROUP    = "tb_usergroup"
	TB_ACCOUNT      = "tb_account"
	TB_MISC         = "tb_misc"
	TB_SESSION      = "tb_session"
)

const (
	ACCOUNT_STATE_DISABLE = 0
	ACCOUNT_STATE_ENABLE  = 1
)

const (
	ACCOUNT_TYPE_SUPERADMIN = 7
	ACCOUNT_TYPE_ADMIN      = 3
	ACCOUNT_TYPE_AUDIT      = 1

	USER_TYPE_USER = 0
)

type Configuration struct {

	// Version is the version which defines the format of the rest of the configuration
	Version string `yaml:"version"`

	// Loglevel is the level at which registry operations are logged.
	Loglevel string `yaml:"loglevel,omitempty"`

	// Debuglevel is the level for debugging.
	Debuglevel string `yaml:"debuglevel,omitempty"`

	HTTP struct {
		Addr    string `yaml:"addr,omitempty"`
		ApiAddr string `yaml:"apiAddr,omitempty"`
	}
}

func ParseConfig(in []byte) (*Configuration, error) {

	config := new(Configuration)

	if err := yaml.Unmarshal(in, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func Parse(rd io.Reader) (*Configuration, error) {

	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	return ParseConfig(in)
}
