// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// ContractsMetaData contains all meta data concerning the Contracts contract.
var ContractsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"allocateVaultYield\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginEpochProcessingWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"estimatedProcessingTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"finalizeEpochWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"subsidiesDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"processingTimeMs\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCurrentEpochId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"startNewEpochWithParticipants\",\"inputs\":[{\"name\":\"participantCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EpochFinalizedWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"totalYieldAvailable\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalSubsidiesDistributed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"processingTimeMs\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerRoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerRoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochProcessingStartedWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"estimatedProcessingTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochStartedWithParticipants\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	ID:  "Contracts",
}

// Contracts is an auto generated Go binding around an Ethereum contract.
type Contracts struct {
	abi abi.ABI
}

// NewContracts creates a new instance of Contracts.
func NewContracts() *Contracts {
	parsed, err := ContractsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Contracts{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Contracts) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAllocateVaultYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf05ca914.
//
// Solidity: function allocateVaultYield(address vault, uint256 amount) returns()
func (contracts *Contracts) PackAllocateVaultYield(vault common.Address, amount *big.Int) []byte {
	enc, err := contracts.abi.Pack("allocateVaultYield", vault, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBeginEpochProcessingWithMetrics is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51e1382d.
//
// Solidity: function beginEpochProcessingWithMetrics(uint256 epochId, uint256 participantCount, uint256 estimatedProcessingTime) returns()
func (contracts *Contracts) PackBeginEpochProcessingWithMetrics(epochId *big.Int, participantCount *big.Int, estimatedProcessingTime *big.Int) []byte {
	enc, err := contracts.abi.Pack("beginEpochProcessingWithMetrics", epochId, participantCount, estimatedProcessingTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackFinalizeEpochWithMetrics is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfca5c0a2.
//
// Solidity: function finalizeEpochWithMetrics(uint256 epochId, uint256 subsidiesDistributed, uint256 processingTimeMs) returns()
func (contracts *Contracts) PackFinalizeEpochWithMetrics(epochId *big.Int, subsidiesDistributed *big.Int, processingTimeMs *big.Int) []byte {
	enc, err := contracts.abi.Pack("finalizeEpochWithMetrics", epochId, subsidiesDistributed, processingTimeMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetCurrentEpochId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa29a839f.
//
// Solidity: function getCurrentEpochId() view returns(uint256)
func (contracts *Contracts) PackGetCurrentEpochId() []byte {
	enc, err := contracts.abi.Pack("getCurrentEpochId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentEpochId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa29a839f.
//
// Solidity: function getCurrentEpochId() view returns(uint256)
func (contracts *Contracts) UnpackGetCurrentEpochId(data []byte) (*big.Int, error) {
	out, err := contracts.abi.Unpack("getCurrentEpochId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackStartNewEpochWithParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x780e52d9.
//
// Solidity: function startNewEpochWithParticipants(uint256 participantCount) returns()
func (contracts *Contracts) PackStartNewEpochWithParticipants(participantCount *big.Int) []byte {
	enc, err := contracts.abi.Pack("startNewEpochWithParticipants", participantCount)
	if err != nil {
		panic(err)
	}
	return enc
}

// ContractsEpochFinalizedWithMetrics represents a EpochFinalizedWithMetrics event raised by the Contracts contract.
type ContractsEpochFinalizedWithMetrics struct {
	EpochId                   *big.Int
	TotalYieldAvailable       *big.Int
	TotalSubsidiesDistributed *big.Int
	ProcessingTimeMs          *big.Int
	Raw                       *types.Log // Blockchain specific contextual infos
}

const ContractsEpochFinalizedWithMetricsEventName = "EpochFinalizedWithMetrics"

// ContractEventName returns the user-defined event name.
func (ContractsEpochFinalizedWithMetrics) ContractEventName() string {
	return ContractsEpochFinalizedWithMetricsEventName
}

// UnpackEpochFinalizedWithMetricsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochFinalizedWithMetrics(uint256 indexed epochId, uint256 totalYieldAvailable, uint256 totalSubsidiesDistributed, uint256 processingTimeMs)
func (contracts *Contracts) UnpackEpochFinalizedWithMetricsEvent(log *types.Log) (*ContractsEpochFinalizedWithMetrics, error) {
	event := "EpochFinalizedWithMetrics"
	if log.Topics[0] != contracts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractsEpochFinalizedWithMetrics)
	if len(log.Data) > 0 {
		if err := contracts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contracts.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractsEpochManagerRoleGranted represents a EpochManagerRoleGranted event raised by the Contracts contract.
type ContractsEpochManagerRoleGranted struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ContractsEpochManagerRoleGrantedEventName = "EpochManagerRoleGranted"

// ContractEventName returns the user-defined event name.
func (ContractsEpochManagerRoleGranted) ContractEventName() string {
	return ContractsEpochManagerRoleGrantedEventName
}

// UnpackEpochManagerRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochManagerRoleGranted(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (contracts *Contracts) UnpackEpochManagerRoleGrantedEvent(log *types.Log) (*ContractsEpochManagerRoleGranted, error) {
	event := "EpochManagerRoleGranted"
	if log.Topics[0] != contracts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractsEpochManagerRoleGranted)
	if len(log.Data) > 0 {
		if err := contracts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contracts.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractsEpochManagerRoleRevoked represents a EpochManagerRoleRevoked event raised by the Contracts contract.
type ContractsEpochManagerRoleRevoked struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ContractsEpochManagerRoleRevokedEventName = "EpochManagerRoleRevoked"

// ContractEventName returns the user-defined event name.
func (ContractsEpochManagerRoleRevoked) ContractEventName() string {
	return ContractsEpochManagerRoleRevokedEventName
}

// UnpackEpochManagerRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochManagerRoleRevoked(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (contracts *Contracts) UnpackEpochManagerRoleRevokedEvent(log *types.Log) (*ContractsEpochManagerRoleRevoked, error) {
	event := "EpochManagerRoleRevoked"
	if log.Topics[0] != contracts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractsEpochManagerRoleRevoked)
	if len(log.Data) > 0 {
		if err := contracts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contracts.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractsEpochProcessingStartedWithMetrics represents a EpochProcessingStartedWithMetrics event raised by the Contracts contract.
type ContractsEpochProcessingStartedWithMetrics struct {
	EpochId                 *big.Int
	ParticipantCount        *big.Int
	EstimatedProcessingTime *big.Int
	Raw                     *types.Log // Blockchain specific contextual infos
}

const ContractsEpochProcessingStartedWithMetricsEventName = "EpochProcessingStartedWithMetrics"

// ContractEventName returns the user-defined event name.
func (ContractsEpochProcessingStartedWithMetrics) ContractEventName() string {
	return ContractsEpochProcessingStartedWithMetricsEventName
}

// UnpackEpochProcessingStartedWithMetricsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochProcessingStartedWithMetrics(uint256 indexed epochId, uint256 participantCount, uint256 estimatedProcessingTime)
func (contracts *Contracts) UnpackEpochProcessingStartedWithMetricsEvent(log *types.Log) (*ContractsEpochProcessingStartedWithMetrics, error) {
	event := "EpochProcessingStartedWithMetrics"
	if log.Topics[0] != contracts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractsEpochProcessingStartedWithMetrics)
	if len(log.Data) > 0 {
		if err := contracts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contracts.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ContractsEpochStartedWithParticipants represents a EpochStartedWithParticipants event raised by the Contracts contract.
type ContractsEpochStartedWithParticipants struct {
	EpochId          *big.Int
	StartTime        *big.Int
	EndTime          *big.Int
	ParticipantCount *big.Int
	Raw              *types.Log // Blockchain specific contextual infos
}

const ContractsEpochStartedWithParticipantsEventName = "EpochStartedWithParticipants"

// ContractEventName returns the user-defined event name.
func (ContractsEpochStartedWithParticipants) ContractEventName() string {
	return ContractsEpochStartedWithParticipantsEventName
}

// UnpackEpochStartedWithParticipantsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochStartedWithParticipants(uint256 indexed epochId, uint256 startTime, uint256 endTime, uint256 participantCount)
func (contracts *Contracts) UnpackEpochStartedWithParticipantsEvent(log *types.Log) (*ContractsEpochStartedWithParticipants, error) {
	event := "EpochStartedWithParticipants"
	if log.Topics[0] != contracts.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ContractsEpochStartedWithParticipants)
	if len(log.Data) > 0 {
		if err := contracts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contracts.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}
