# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important
- ALL instructions within this document MUST BE FOLLOWED, these are not optional unless explicitly stated.
- ASK FOR CLARIFICATION If you are uncertain of any of thing within the document.
- DO NOT edit more code than you have to.
- DO NOT WASTE TOKENS, be succinct and concise.
- Use parallel subagents for tasks that can be executed concurrently. Split tasks into smaller, manageable parts, ALWAYS edit file only in one subtask.
- For Go code: always write comments starting with lowercase letters (e.g., `// processData handles...` not `// ProcessData handles...`)


## Architecture Overview

This is a Go-based epoch server for managing epoch-based subsidy distribution in an NFT collection-backed lending protocol. The system orchestrates time-based yield allocation cycles through blockchain smart contracts.

### Core Services Architecture

The system is built around three primary services with clear boundaries:

- **Epoch Service** (`internal/services/epoch/`): Manages epoch lifecycle (start, force-end, earnings calculation)
- **Merkle Service** (`internal/services/merkle/`): Generates cryptographic proofs for subsidy distribution using BadgerDB snapshots
- **Subsidy Service** (`internal/services/subsidy/`): Handles subsidy distribution (interface-based, currently mock implementation)
- **Scheduler Service** (`internal/services/scheduler/`): Orchestrates automated epoch operations at configurable intervals

### Data Flow Pattern

External Subgraph → Business Logic Services → BadgerDB Storage → Blockchain Contracts → API Endpoints

- **Subgraph Integration**: GraphQL client (`internal/infra/subgraph/`) queries historical account/epoch data
- **Storage Layer**: BadgerDB (`internal/infra/storage/`) stores merkle snapshots and processed epoch data
- **Blockchain Client**: Unified client (`internal/infra/blockchain/`) handles all smart contract interactions
- **API Layer**: RESTful endpoints (`internal/api/`) expose operations and data queries

### Contract Integration

The server integrates with multiple smart contracts using generated Go bindings (`pkg/contracts/`):
- **EpochManager**: Start epochs, force-end epochs, get current epoch ID
- **DebtSubsidizer**: Distribute subsidies with merkle roots
- **LendingManager**: Update exchange rates
- **CollectionsVault**: Query vault-specific data

**Contract Locations**:
- **Source contracts**: `/Users/andrey/projects/lend.fam MVP/collection-vault/src/interfaces/`
- **Generated bindings**: `pkg/contracts/`
- **Contract ABIs**: `/Users/andrey/projects/lend.fam MVP/collection-vault/out/`
- **Subgraph**: `/Users/andrey/projects/lend.fam MVP/rewards-subgraph/`

## Code Quality Configuration

### golangci-lint
- **Configuration file**: `.golangci.yml` (comprehensive linting rules)
- **Enabled linters**: 40+ linters including security, style, and performance checks
- **Custom settings**: Go 1.22+ support, project-specific exclusions, security rules
- **Integration**: Runs via `make lint` command

## Common Development Commands

### Building and Running
```bash
# Build server
go build ./cmd/server

# Run server (requires environment variables)
./server -config configs/config.yaml

# Build using Makefile
make build

# Run with development pipeline
make dev
```

### Testing
```bash
# Run all tests
go test ./...

# Run only unit tests
make test-unit

# Run integration tests (requires containers)
make integration-test

# Run with race detection
make test-race

# Generate coverage report
make coverage
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (uses .golangci.yml config)
make lint

# Run security scan (requires gosec)
make security-scan

# Generate mocks (uses moq)
make generate-mocks
```

## Configuration Requirements

### Required Environment Variables
```bash
# Blockchain connection
RPC_URL="https://rpc.example.com"
PRIVATE_KEY="0x..."

# Subgraph endpoint
SUBGRAPH_ENDPOINT="https://subgraph.example.com"

# Contract addresses (all required)
COMPTROLLER_ADDRESS="0x..."
EPOCH_MANAGER_ADDRESS="0x..."
DEBT_SUBSIDIZER_PROXY_ADDRESS="0x..."
LENDING_MANAGER_ADDRESS="0x..."
COLLECTION_REGISTRY_ADDRESS="0x..."
VAULT_ADDRESS="0x..."
```

### Optional Configuration
```bash
# Server settings
SERVER_HOST="0.0.0.0"
SERVER_PORT="8080"

# Database
DATABASE_TYPE="memory"  # or "badger"
DATABASE_CONNECTION_STRING="/path/to/db"

# Scheduler
SCHEDULER_ENABLED="true"
SCHEDULER_INTERVAL="1h"
SCHEDULER_TIMEZONE="UTC"
```

## Development Patterns

### Service Interface Design
- All services implement interfaces for testability
- Dependencies are injected via constructor functions
- Mock implementations generated using moq (`*_mocks.go`)

### Error Handling
- Structured error types for different failure scenarios
- Comprehensive logging with context using `github.com/go-pkgz/lgr`
- Graceful degradation for non-critical operations

### Storage Abstraction
- BadgerDB is abstracted behind a generic storage interface
- Configuration-driven storage type selection
- Integration tests use testcontainers for isolation

### Blockchain Integration
- Unified client handles all contract interactions
- Transaction management with gas estimation and waiting
- Error classification for different blockchain failure types

## Key API Endpoints

```
POST /api/epochs/start              - Start new epoch
POST /api/epochs/force-end          - Force end current epoch  
POST /api/epochs/distribute         - Distribute subsidies
GET /api/users/{address}/total-earned - Get user earnings
GET /api/users/{address}/merkle-proof - Get current merkle proof
GET /api/users/{address}/merkle-proof/epoch/{epochNumber} - Get historical proof
GET /health                         - Health check with service status
GET /swagger/                       - API documentation
```

## Testing Strategy

### Unit Tests
- Mock external dependencies (blockchain, subgraph, storage) using moq
- Test business logic in isolation
- Use testify for assertions

### Integration Tests
- Use testcontainers for BadgerDB testing
- Test cross-service interactions
- Validate storage consistency and performance

### Contract Compatibility Tests
- Verify merkle proof generation matches Solidity implementation
- Test leaf hash calculation compatibility
- Validate proof verification logic

## Common Development Tasks

### Adding New Service
1. **Create service directory structure**:
   ```
   internal/services/{service}/
   ├── {service}.go          # Service interface definition
   ├── {service}_mocks.go    # Generated mocks (via moq)
   ├── model.go             # Data models and types
   ├── errors.go            # Service-specific error types
   └── {service}impl/
       ├── service.go       # Service implementation
       └── store.go         # Data storage logic (if needed)
   ```

2. **Define service interface** (`{service}.go`):
   ```go
   package {service}
   
   import "context"
   
   //go:generate moq -out {service}_mocks.go . Service
   
   type Service interface {
       // Define your service methods here
       ProcessData(ctx context.Context, input *SomeInput) (*SomeOutput, error)
   }
   ```

3. **Create data models** (`model.go`):
   ```go
   package {service}
   
   type SomeInput struct {
       Field1 string `json:"field1"`
       Field2 int    `json:"field2"`
   }
   
   type SomeOutput struct {
       Result string `json:"result"`
   }
   ```

4. **Define service errors** (`errors.go`):
   ```go
   package {service}
   
   import "errors"
   
   var (
       ErrInvalidInput = errors.New("invalid input")
       ErrNotFound     = errors.New("resource not found")
   )
   ```

5. **Implement service** (`{service}impl/service.go`):
   ```go
   package {service}impl
   
   import (
       "context"
       "github.com/andrey/epoch-server/internal/services/{service}"
       "github.com/go-pkgz/lgr"
   )
   
   type Service struct {
       logger lgr.L
       // Add other dependencies as needed
   }
   
   func NewService(logger lgr.L) {service}.Service {
       return &Service{
           logger: logger,
       }
   }
   
   func (s *Service) ProcessData(ctx context.Context, input *{service}.SomeInput) (*{service}.SomeOutput, error) {
       // Implementation here
   }
   ```

6. **Generate mocks**:
   ```bash
   go generate ./internal/services/{service}/
   ```

7. **Register service in main server** (`cmd/server/main.go`):
   ```go
   // Add service initialization
   {service}Service := {service}impl.NewService(logger)
   
   // Pass to server constructor
   server := api.NewServer(
       epochService,
       subsidyService,
       merkleService,
       {service}Service, // Add here
       logger,
       cfg,
   )
   ```

8. **Update server struct** (`internal/api/server.go`):
   ```go
   type Server struct {
       // existing services...
       {service}Service {service}.Service
   }
   
   func NewServer(
       // existing parameters...
       {service}Service {service}.Service,
       logger lgr.L,
       cfg *config.Config,
   ) *Server {
       return &Server{
           // existing assignments...
           {service}Service: {service}Service,
       }
   }
   ```

9. **Create API handlers** (`internal/api/handlers/{service}.go`):
   ```go
   package handlers
   
   import (
       "net/http"
       "github.com/andrey/epoch-server/internal/services/{service}"
       "github.com/go-pkgz/lgr"
       "github.com/go-pkgz/rest"
   )
   
   type {Service}Handler struct {
       {service}Service {service}.Service
       logger          lgr.L
   }
   
   func New{Service}Handler({service}Service {service}.Service, logger lgr.L) *{Service}Handler {
       return &{Service}Handler{
           {service}Service: {service}Service,
           logger:          logger,
       }
   }
   
   // @Summary Process data
   // @Description Process data using the service
   // @Tags {service}
   // @Accept json
   // @Produce json
   // @Param input body {service}.SomeInput true "Input data"
   // @Success 200 {object} {service}.SomeOutput
   // @Failure 400 {object} handlers.ErrorResponse
   // @Router /api/{service}/process [post]
   func (h *{Service}Handler) ProcessData(w http.ResponseWriter, r *http.Request) {
       // Handler implementation
   }
   ```

10. **Register routes** (`internal/api/server.go` in `SetupRoutes`):
    ```go
    // Create handler
    {service}Handler := handlers.New{Service}Handler(s.{service}Service, s.logger)
    
    // Register routes
    api.Route("/api/{service}", func(r chi.Router) {
        r.Post("/process", {service}Handler.ProcessData)
    })
    ```

11. **Write tests** (`{service}impl/service_test.go`):
    ```go
    package {service}impl
    
    import (
        "context"
        "testing"
        "github.com/stretchr/testify/assert"
        "github.com/andrey/epoch-server/internal/services/{service}"
    )
    
    func TestService_ProcessData(t *testing.T) {
        // Test implementation using generated mocks
    }
    ```

### Updating Contract Bindings
1. **Generate bindings from collection-vault contracts**:
   ```bash
   # Script automatically builds contracts and generates Go bindings
   ./scripts/generate_bindings.sh
   ```
   
2. **Generated bindings** (`pkg/contracts/`):
   - `ICollectionRegistry.go` - Collection registry interface
   - `ICollectionsVault.go` - Main vault interface  
   - `IDebtSubsidizer.go` - Subsidy distribution interface
   - `IEpochManager.go` - Epoch lifecycle management interface
   - `ILendingManager.go` - Compound integration interface
   - `IERC20.go` - Standard ERC20 interface
   
3. **Contract source location**: `/Users/andrey/projects/lend.fam MVP/collection-vault/src/interfaces/`

4. **Update contract addresses** in configuration after deployment

5. **Test integration** with new contract methods

### Adding New Storage Backend
1. Implement storage interface in `internal/infra/storage/`
2. Add configuration options for new backend
3. Update storage factory to support new type
4. Add integration tests for new backend

### Performance Optimization
- Use `make profile-cpu` and `make profile-mem` for profiling
- Run `make benchmark` for performance testing
- Monitor BadgerDB performance with integration benchmarks