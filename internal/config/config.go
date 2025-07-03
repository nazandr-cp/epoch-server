package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Type             string `yaml:"type"`
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`

	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
	} `yaml:"logging"`

	Ethereum struct {
		RPCURL     string `yaml:"rpc_url"`
		PrivateKey string `yaml:"private_key"`
		GasLimit   uint64 `yaml:"gas_limit"`
		GasPrice   string `yaml:"gas_price"`
	} `yaml:"ethereum"`

	Subgraph struct {
		Endpoint       string        `yaml:"endpoint"`
		Timeout        time.Duration `yaml:"timeout"`
		MaxRetries     int           `yaml:"max_retries"`
		PaginationSize int           `yaml:"pagination_size"`
	} `yaml:"subgraph"`

	Scheduler struct {
		Interval time.Duration `yaml:"interval"`
		Enabled  bool          `yaml:"enabled"`
		Timezone string        `yaml:"timezone"`
	} `yaml:"scheduler"`

	Contracts struct {
		Comptroller        string `yaml:"comptroller"`
		EpochManager       string `yaml:"epoch_manager"`
		DebtSubsidizer     string `yaml:"debt_subsidizer"`
		LendingManager     string `yaml:"lending_manager"`
		CollectionRegistry string `yaml:"collection_registry"`
	} `yaml:"contracts"`

	Features struct {
		DryRun        bool `yaml:"dry_run"`
		EnableMetrics bool `yaml:"enable_metrics"`
		DebugMode     bool `yaml:"debug_mode"`
	} `yaml:"features"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	if privateKey := os.Getenv("ETHEREUM_PRIVATE_KEY"); privateKey != "" {
		cfg.Ethereum.PrivateKey = privateKey
	}

	return &cfg, nil
}
