package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	mcConfig "github.com/hashcloak/Meson/client/config"
)

type Config struct {
	KSLocation     string
	DefaultChainID int64
	Chain          map[string]struct{ Ticker, Endpoint string }
	Meson          *mcConfig.Config
}

// FixupAndValidate applies defaults to config entries and validates the
// supplied configuration.  Most people should call one of the Load variants
// instead.
func (c *Config) FixupAndValidate() error {
	if info, err := os.Stat(c.KSLocation); err != nil || !info.IsDir() {
		return fmt.Errorf("error key store location \"%s\"", c.KSLocation)
	}
	if _, ok := c.Chain[fmt.Sprint(c.DefaultChainID)]; !ok {
		return fmt.Errorf("missing ticker/endpoint for the default chain id")
	}
	if c.Meson == nil {
		return fmt.Errorf("missing meson config")
	}
	return c.Meson.FixupAndMinimallyValidate()
}

// Load parses and validates the provided buffer b as a config file body and
// returns the Config.
func Load(b []byte) (*Config, error) {
	cfg := new(Config)
	md, err := toml.Decode(string(b), cfg)
	if err != nil {
		return nil, err
	}
	if undecoded := md.Undecoded(); len(undecoded) != 0 {
		return nil, fmt.Errorf("config: Undecoded keys in config file: %v", undecoded)
	}
	if err := cfg.FixupAndValidate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadFile loads, parses, and validates the provided file and returns the
// Config.
func LoadFile(f string) (*Config, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return Load(b)
}

// SaveFile saves the config to the provided file
func (c *Config) SaveFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := toml.NewEncoder(f)
	return enc.Encode(c)
}
