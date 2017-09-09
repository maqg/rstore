package configuration

import (
	"fmt"
	"io"
	"io/ioutil"
	"octlink/rstore/utils/octlog"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	// BlobSize for blob manage in bytes
	BlobSize = 8 * 1024 * 1024
)

// Configuration for global config
type Configuration struct {

	// Version is the version which defines the format of the rest of the configuration
	Version string `yaml:"version"`

	// Loglevel is the level at which registry operations are logged.
	LogLevel int `yaml:"loglevel,omitempty"`

	// LogDirectory log dir like "var/logs/"
	LogDirectory string `yaml:"logdirectory,omitempty"`

	// Debuglevel is the level for debugging.
	DebugLevel int `yaml:"debuglevel,omitempty"`

	// HugeBlob true to support huge blob, false to split blobs
	HugeBlob bool `yaml:"hugeblob,omitempty"`

	HTTP struct {
		Addr    string `yaml:"addr,omitempty"`
		APIAddr string `yaml:"apiAddr,omitempty"`
	}

	RootDirectory string `yaml:"rootdirectory,omitempty"`
}

// Conf global configuration
var Conf *Configuration

// GetConfig for global
func GetConfig() *Configuration {
	return Conf
}

// RootDirectory for reposity directory fetching
func RootDirectory() string {
	return Conf.RootDirectory
}

// HugeBlob for huge blob support or not
func HugeBlob() bool {
	return Conf.HugeBlob
}

// ParseConfig from yml
func ParseConfig(in []byte) (*Configuration, error) {

	config := new(Configuration)

	if err := yaml.Unmarshal(in, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// Parse from io.Reader
func Parse(rd io.Reader) (*Configuration, error) {

	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	return ParseConfig(in)
}

// ResolveConfig for config convert from yml
func ResolveConfig(configfile string) (*Configuration, error) {

	var configurationPath string

	if configfile == "" {
		configurationPath = os.Getenv("REGISTRY_CONFIGURATION_PATH")
	} else {
		configurationPath = configfile
	}

	if configurationPath == "" {
		return nil, fmt.Errorf("configuration path unspecified")
	}

	fp, err := os.Open(configurationPath)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	config, err := Parse(fp)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", configurationPath, err)
	}

	Conf = config

	return config, nil
}

func init() {
	Conf = &Configuration{
		LogLevel:      octlog.DebugLevel,
		DebugLevel:    octlog.DebugLevel,
		RootDirectory: "./",
		LogDirectory:  "./var/logs",
	}
}
