package configuration

type Configuration struct {

	// Version is the version which defines the format of the rest of the configuration
	Version Version `yaml:"version"`

	// Loglevel is the level at which registry operations are logged.
	Loglevel Loglevel `yaml:"loglevel,omitempty"`

	// Debuglevel is the level for debugging.
	Debuglevel Debuglevel `yaml:"debuglevel,omitempty"`

	HTTP struct {
		Addr string `yaml:"addr,omitempty"`
	}
}
