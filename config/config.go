package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// Path to the client config file
	ClientCfgFile string
	// Ticker
	Ticker string
	// Service name
	Service string
	// Private key (in Hex form) used to sign the txn
	PrivKey string
	// Chain ID for specific ETH-based chain
	ChainID int64
	// Ethereum rpc endpoint
	RpcEndpoint string
}

func DefaultConfig() (cfg *Config) {
	cfg = new(Config)
	_ = cfg.FixupAndValidate()
	return
}

// FixupAndValidate applies defaults to config entries and validates the
// supplied configuration.  Most people should call one of the Load variants
// instead.
func (c *Config) FixupAndValidate() error {
	// TODO
	if c.PrivKey == "" {
		return fmt.Errorf("private key is empty")
	}
	return nil
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
func SaveFile(f string, config *Config) error {
	file, err := os.Create(f)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := toml.NewEncoder(file)
	return enc.Encode(config)
}
