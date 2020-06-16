package fwatch

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config - config data
type Config struct {
	path    string
	Targets []Target
}

// Target - target
type Target struct {
	Path   string
	Script string
	Type   []string
}

// Load - load config yml file
func Load(path string) (config *Config, err error) {
	config = &Config{}

	// get absolute path
	path, err = filepath.Abs(path)
	print("load config: %s", path)

	// load toml to Config struct
	if _, err = toml.DecodeFile(path, config); err != nil {
		print("error: %v", err)
		return nil, err
	}

	config.path = path

	return config, nil
}

// Register -
func (c Config) Register(p, s string, t []string) error {
	target := Target{p, s, t}
	c.Targets = append(c.Targets, target)

	f, err := os.OpenFile(c.path, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		print("error: %v\n", err)
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	enc.Indent = "  "
	err = enc.Encode(c)

	return err
}

// Unregister -
func (c Config) Unregister(index int) error {
	c.Targets = append(c.Targets[:index], c.Targets[index+1:]...)

	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	err = enc.Encode(c)

	return err
}
