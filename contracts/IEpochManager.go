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

// IEpochManagerMetaData contains all meta data concerning the IEpochManager contract.
var IEpochManagerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"allocateVaultYield\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginEpochProcessingWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"estimatedProcessingTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"finalizeEpochWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"subsidiesDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"processingTimeMs\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCurrentEpochId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"startNewEpochWithParticipants\",\"inputs\":[{\"name\":\"participantCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EpochFinalizedWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"totalYieldAvailable\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalSubsidiesDistributed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"processingTimeMs\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerRoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerRoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochProcessingStartedWithMetrics\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"estimatedProcessingTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochStartedWithParticipants\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	ID:  "IEpochManager",
}

// IEpochManager is an auto generated Go binding around an Ethereum contract.
type IEpochManager struct {
	abi abi.ABI
}

// NewIEpochManager creates a new instance of IEpochManager.
func NewIEpochManager() *IEpochManager {
	parsed, err := IEpochManagerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IEpochManager{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IEpochManager) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAllocateVaultYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf05ca914.
//
// Solidity: function allocateVaultYield(address vault, uint256 amount) returns()
func (iEpochManager *IEpochManager) PackAllocateVaultYield(vault common.Address, amount *big.Int) []byte {
	enc, err := iEpochManager.abi.Pack("allocateVaultYield", vault, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBeginEpochProcessingWithMetrics is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x51e1382d.
//
// Solidity: function beginEpochProcessingWithMetrics(uint256 epochId, uint256 participantCount, uint256 estimatedProcessingTime) returns()
func (iEpochManager *IEpochManager) PackBeginEpochProcessingWithMetrics(epochId *big.Int, participantCount *big.Int, estimatedProcessingTime *big.Int) []byte {
	enc, err := iEpochManager.abi.Pack("beginEpochProcessingWithMetrics", epochId, participantCount, estimatedProcessingTime)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackFinalizeEpochWithMetrics is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfca5c0a2.
//
// Solidity: function finalizeEpochWithMetrics(uint256 epochId, uint256 subsidiesDistributed, uint256 processingTimeMs) returns()
func (iEpochManager *IEpochManager) PackFinalizeEpochWithMetrics(epochId *big.Int, subsidiesDistributed *big.Int, processingTimeMs *big.Int) []byte {
	enc, err := iEpochManager.abi.Pack("finalizeEpochWithMetrics", epochId, subsidiesDistributed, processingTimeMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetCurrentEpochId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa29a839f.
//
// Solidity: function getCurrentEpochId() view returns(uint256)
func (iEpochManager *IEpochManager) PackGetCurrentEpochId() []byte {
	enc, err := iEpochManager.abi.Pack("getCurrentEpochId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentEpochId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa29a839f.
//
// Solidity: function getCurrentEpochId() view returns(uint256)
func (iEpochManager *IEpochManager) UnpackGetCurrentEpochId(data []byte) (*big.Int, error) {
	out, err := iEpochManager.abi.Unpack("getCurrentEpochId", data)
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
func (iEpochManager *IEpochManager) PackStartNewEpochWithParticipants(participantCount *big.Int) []byte {
	enc, err := iEpochManager.abi.Pack("startNewEpochWithParticipants", participantCount)
	if err != nil {
		panic(err)
	}
	return enc
}

// IEpochManagerEpochFinalizedWithMetrics represents a EpochFinalizedWithMetrics event raised by the IEpochManager contract.
type IEpochManagerEpochFinalizedWithMetrics struct {
	EpochId                   *big.Int
	TotalYieldAvailable       *big.Int
	TotalSubsidiesDistributed *big.Int
	ProcessingTimeMs          *big.Int
	Raw                       *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochFinalizedWithMetricsEventName = "EpochFinalizedWithMetrics"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochFinalizedWithMetrics) ContractEventName() string {
	return IEpochManagerEpochFinalizedWithMetricsEventName
}

// UnpackEpochFinalizedWithMetricsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochFinalizedWithMetrics(uint256 indexed epochId, uint256 totalYieldAvailable, uint256 totalSubsidiesDistributed, uint256 processingTimeMs)
func (iEpochManager *IEpochManager) UnpackEpochFinalizedWithMetricsEvent(log *types.Log) (*IEpochManagerEpochFinalizedWithMetrics, error) {
	event := "EpochFinalizedWithMetrics"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochFinalizedWithMetrics)
	if len(log.Data) > 0 {
		if err := iEpochManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEpochManager.abi.Events[event].Inputs {
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

// IEpochManagerEpochManagerRoleGranted represents a EpochManagerRoleGranted event raised by the IEpochManager contract.
type IEpochManagerEpochManagerRoleGranted struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochManagerRoleGrantedEventName = "EpochManagerRoleGranted"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochManagerRoleGranted) ContractEventName() string {
	return IEpochManagerEpochManagerRoleGrantedEventName
}

// UnpackEpochManagerRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochManagerRoleGranted(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (iEpochManager *IEpochManager) UnpackEpochManagerRoleGrantedEvent(log *types.Log) (*IEpochManagerEpochManagerRoleGranted, error) {
	event := "EpochManagerRoleGranted"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochManagerRoleGranted)
	if len(log.Data) > 0 {
		if err := iEpochManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEpochManager.abi.Events[event].Inputs {
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

// IEpochManagerEpochManagerRoleRevoked represents a EpochManagerRoleRevoked event raised by the IEpochManager contract.
type IEpochManagerEpochManagerRoleRevoked struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochManagerRoleRevokedEventName = "EpochManagerRoleRevoked"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochManagerRoleRevoked) ContractEventName() string {
	return IEpochManagerEpochManagerRoleRevokedEventName
}

// UnpackEpochManagerRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochManagerRoleRevoked(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (iEpochManager *IEpochManager) UnpackEpochManagerRoleRevokedEvent(log *types.Log) (*IEpochManagerEpochManagerRoleRevoked, error) {
	event := "EpochManagerRoleRevoked"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochManagerRoleRevoked)
	if len(log.Data) > 0 {
		if err := iEpochManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEpochManager.abi.Events[event].Inputs {
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

// IEpochManagerEpochProcessingStartedWithMetrics represents a EpochProcessingStartedWithMetrics event raised by the IEpochManager contract.
type IEpochManagerEpochProcessingStartedWithMetrics struct {
	EpochId                 *big.Int
	ParticipantCount        *big.Int
	EstimatedProcessingTime *big.Int
	Raw                     *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochProcessingStartedWithMetricsEventName = "EpochProcessingStartedWithMetrics"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochProcessingStartedWithMetrics) ContractEventName() string {
	return IEpochManagerEpochProcessingStartedWithMetricsEventName
}

// UnpackEpochProcessingStartedWithMetricsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochProcessingStartedWithMetrics(uint256 indexed epochId, uint256 participantCount, uint256 estimatedProcessingTime)
func (iEpochManager *IEpochManager) UnpackEpochProcessingStartedWithMetricsEvent(log *types.Log) (*IEpochManagerEpochProcessingStartedWithMetrics, error) {
	event := "EpochProcessingStartedWithMetrics"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochProcessingStartedWithMetrics)
	if len(log.Data) > 0 {
		if err := iEpochManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEpochManager.abi.Events[event].Inputs {
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

// IEpochManagerEpochStartedWithParticipants represents a EpochStartedWithParticipants event raised by the IEpochManager contract.
type IEpochManagerEpochStartedWithParticipants struct {
	EpochId          *big.Int
	StartTime        *big.Int
	EndTime          *big.Int
	ParticipantCount *big.Int
	Raw              *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochStartedWithParticipantsEventName = "EpochStartedWithParticipants"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochStartedWithParticipants) ContractEventName() string {
	return IEpochManagerEpochStartedWithParticipantsEventName
}

// UnpackEpochStartedWithParticipantsEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochStartedWithParticipants(uint256 indexed epochId, uint256 startTime, uint256 endTime, uint256 participantCount)
func (iEpochManager *IEpochManager) UnpackEpochStartedWithParticipantsEvent(log *types.Log) (*IEpochManagerEpochStartedWithParticipants, error) {
	event := "EpochStartedWithParticipants"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochStartedWithParticipants)
	if len(log.Data) > 0 {
		if err := iEpochManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEpochManager.abi.Events[event].Inputs {
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
