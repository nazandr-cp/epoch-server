package config

import (
	"time"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	// Server configuration
	Server struct {
		Host string `long:"server-host" env:"SERVER_HOST" default:"0.0.0.0" description:"Server host"`
		Port int    `long:"server-port" env:"SERVER_PORT" default:"8080" description:"Server port"`
	} `group:"Server Options" namespace:"server"`

	// Database configuration
	Database struct {
		Type             string `long:"database-type" env:"DATABASE_TYPE" default:"memory" description:"Database type"`
		ConnectionString string `long:"database-connection-string" env:"DATABASE_CONNECTION_STRING" default:"" description:"Database connection string"`
	} `group:"Database Options" namespace:"database"`

	// Logging configuration
	Logging struct {
		Level  string `long:"log-level" env:"LOG_LEVEL" default:"debug" description:"Log level"`
		Format string `long:"log-format" env:"LOG_FORMAT" default:"json" description:"Log format"`
		Output string `long:"log-output" env:"LOG_OUTPUT" default:"stdout" description:"Log output"`
	} `group:"Logging Options" namespace:"logging"`

	// Ethereum configuration
	Ethereum struct {
		RPCURL     string `long:"rpc-url" env:"RPC_URL" required:"true" description:"Ethereum RPC URL"`
		PrivateKey string `long:"private-key" env:"PRIVATE_KEY" required:"true" description:"Ethereum private key"`
		Sender     string `long:"sender" env:"SENDER" description:"Sender address"`
		GasLimit   uint64 `long:"gas-limit" env:"GAS_LIMIT" default:"500000" description:"Gas limit"`
		GasPrice   string `long:"gas-price" env:"GAS_PRICE" default:"20000000000" description:"Gas price"`
	} `group:"Ethereum Options" namespace:"ethereum"`

	// Subgraph configuration
	Subgraph struct {
		Endpoint       string        `long:"subgraph-endpoint" env:"SUBGRAPH_ENDPOINT" required:"true" description:"Subgraph endpoint"`
		Timeout        time.Duration `long:"subgraph-timeout" env:"SUBGRAPH_TIMEOUT" default:"30s" description:"Subgraph timeout"`
		MaxRetries     int           `long:"subgraph-max-retries" env:"SUBGRAPH_MAX_RETRIES" default:"3" description:"Subgraph max retries"`
		PaginationSize int           `long:"subgraph-pagination-size" env:"SUBGRAPH_PAGINATION_SIZE" default:"1000" description:"Subgraph pagination size"`
	} `group:"Subgraph Options" namespace:"subgraph"`

	// Scheduler configuration
	Scheduler struct {
		Interval time.Duration `long:"scheduler-interval" env:"SCHEDULER_INTERVAL" default:"1h" description:"Scheduler interval"`
		Enabled  bool          `long:"scheduler-enabled" env:"SCHEDULER_ENABLED" description:"Enable scheduler"`
		Timezone string        `long:"scheduler-timezone" env:"SCHEDULER_TIMEZONE" default:"UTC" description:"Scheduler timezone"`
	} `group:"Scheduler Options" namespace:"scheduler"`

	// Contract addresses
	Contracts struct {
		Comptroller        string `long:"comptroller-address" env:"COMPTROLLER_ADDRESS" required:"true" description:"Comptroller contract address"`
		EpochManager       string `long:"epoch-manager-address" env:"EPOCH_MANAGER_ADDRESS" required:"true" description:"Epoch manager contract address"`
		DebtSubsidizer     string `long:"debt-subsidizer-address" env:"DEBT_SUBSIDIZER_PROXY_ADDRESS" required:"true" description:"Debt subsidizer contract address"`
		LendingManager     string `long:"lending-manager-address" env:"LENDING_MANAGER_ADDRESS" required:"true" description:"Lending manager contract address"`
		CollectionRegistry string `long:"collection-registry-address" env:"COLLECTION_REGISTRY_ADDRESS" required:"true" description:"Collection registry contract address"`
		CollectionsVault   string `long:"collections-vault-address" env:"VAULT_ADDRESS" required:"true" description:"Collections vault contract address"`
		Asset              string `long:"asset-address" env:"ASSET_ADDRESS" description:"Asset contract address"`
		NFT                string `long:"nft-address" env:"NFT_ADDRESS" description:"NFT contract address"`
		CToken             string `long:"ctoken-address" env:"CTOKEN_ADDRESS" description:"CToken contract address"`
	} `group:"Contract Options" namespace:"contracts"`
}

func Load() (*Config, error) {
	var cfg Config

	parser := flags.NewParser(&cfg, flags.Default)

	// Parse only environment variables and ignore command line arguments
	if _, err := parser.ParseArgs([]string{}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
