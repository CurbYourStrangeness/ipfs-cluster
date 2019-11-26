package numpin

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ipfs/ipfs-cluster/config"
	"github.com/kelseyhightower/envconfig"
)

const configKey = "numpin"
const envConfigKey = "cluster_numpin"

// These are the default values for a Config.
const (
	DefaultMetricTTL = 10 * time.Second
)

// Config allows to initialize an Informer.
type Config struct {
	config.Saver

	MetricTTL time.Duration
}

type jsonConfig struct {
	MetricTTL string `json:"metric_ttl"`
}

// ConfigKey returns a human-friendly identifier for this
// Config's type.
func (cfg *Config) ConfigKey() string {
	return configKey
}

// Default initializes this Config with sensible values.
func (cfg *Config) Default() error {
	cfg.MetricTTL = DefaultMetricTTL
	return nil
}

// ApplyEnvVars fills in any Config fields found
// as environment variables.
func (cfg *Config) ApplyEnvVars() error {
	jcfg := cfg.toJSONConfig()

	err := envconfig.Process(envConfigKey, jcfg)
	if err != nil {
		return err
	}

	return cfg.applyJSONConfig(jcfg)
}

// Validate checks that the fields of this configuration have
// sensible values.
func (cfg *Config) Validate() error {
	if cfg.MetricTTL <= 0 {
		return errors.New("disk.metric_ttl is invalid")
	}

	return nil
}

// LoadJSON parses a raw JSON byte-slice as generated by ToJSON().
func (cfg *Config) LoadJSON(raw []byte) error {
	jcfg := &jsonConfig{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		return err
	}

	cfg.Default()

	return cfg.applyJSONConfig(jcfg)
}

func (cfg *Config) applyJSONConfig(jcfg *jsonConfig) error {
	t, _ := time.ParseDuration(jcfg.MetricTTL)
	cfg.MetricTTL = t

	return cfg.Validate()
}

// ToJSON generates a human-friendly JSON representation of this Config.
func (cfg *Config) ToJSON() ([]byte, error) {
	jcfg := cfg.toJSONConfig()

	return config.DefaultJSONMarshal(jcfg)
}

func (cfg *Config) toJSONConfig() *jsonConfig {
	return &jsonConfig{
		MetricTTL: cfg.MetricTTL.String(),
	}
}

func (cfg *Config) String() (string, error) {
	bytes, err := config.DefaultJSONMarshalWithoutHiddenFields(*cfg.toJSONConfig())
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
