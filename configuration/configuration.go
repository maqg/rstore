package configuration

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {

	// Version is the version which defines the format of the rest of the configuration
	Version string `yaml:"version"`

	// Loglevel is the level at which registry operations are logged.
	Loglevel string `yaml:"loglevel,omitempty"`

	// Debuglevel is the level for debugging.
	Debuglevel string `yaml:"debuglevel,omitempty"`

	HTTP struct {
		Addr string `yaml:"addr,omitempty"`
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
