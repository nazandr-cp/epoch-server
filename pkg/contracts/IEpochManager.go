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
	ABI: "[{\"type\":\"function\",\"name\":\"allocateVaultYield\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endEpochWithSubsidies\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"subsidiesDistributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forceEndEpochWithZeroYield\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCurrentEpochId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVaultYieldForEpoch\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantVaultRole\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeVaultRole\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAutomatedSystem\",\"inputs\":[{\"name\":\"newAutomatedSystem\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDebtSubsidizer\",\"inputs\":[{\"name\":\"newDebtSubsidizer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AutomatedSystemUpdated\",\"inputs\":[{\"name\":\"newAutomatedSystem\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DebtSubsidizerUpdated\",\"inputs\":[{\"name\":\"newDebtSubsidizer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochFailed\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochFinalized\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"totalYieldAvailable\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalSubsidiesDistributed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerRoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerRoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochProcessingStarted\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochStarted\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProcessingFailed\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProcessingStarted\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultYieldAllocated\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EpochManager__EpochNotEnded\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EpochManager__EpochStillActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EpochManager__InvalidEpochDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EpochManager__InvalidEpochId\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EpochManager__InvalidEpochStatus\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"currentStatus\",\"type\":\"uint8\",\"internalType\":\"enumIEpochManager.EpochStatus\"},{\"name\":\"expectedStatus\",\"type\":\"uint8\",\"internalType\":\"enumIEpochManager.EpochStatus\"}]},{\"type\":\"error\",\"name\":\"EpochManager__Unauthorized\",\"inputs\":[]}]",
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

// PackEndEpochWithSubsidies is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x680843cc.
//
// Solidity: function endEpochWithSubsidies(uint256 epochId, address vaultAddress, bytes32 merkleRoot, uint256 subsidiesDistributed) returns()
func (iEpochManager *IEpochManager) PackEndEpochWithSubsidies(epochId *big.Int, vaultAddress common.Address, merkleRoot [32]byte, subsidiesDistributed *big.Int) []byte {
	enc, err := iEpochManager.abi.Pack("endEpochWithSubsidies", epochId, vaultAddress, merkleRoot, subsidiesDistributed)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackForceEndEpochWithZeroYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf78e6ce1.
//
// Solidity: function forceEndEpochWithZeroYield(uint256 epochId, address vaultAddress) returns()
func (iEpochManager *IEpochManager) PackForceEndEpochWithZeroYield(epochId *big.Int, vaultAddress common.Address) []byte {
	enc, err := iEpochManager.abi.Pack("forceEndEpochWithZeroYield", epochId, vaultAddress)
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

// PackGetVaultYieldForEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa74a003.
//
// Solidity: function getVaultYieldForEpoch(uint256 epochId, address vault) view returns(uint256)
func (iEpochManager *IEpochManager) PackGetVaultYieldForEpoch(epochId *big.Int, vault common.Address) []byte {
	enc, err := iEpochManager.abi.Pack("getVaultYieldForEpoch", epochId, vault)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetVaultYieldForEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa74a003.
//
// Solidity: function getVaultYieldForEpoch(uint256 epochId, address vault) view returns(uint256)
func (iEpochManager *IEpochManager) UnpackGetVaultYieldForEpoch(data []byte) (*big.Int, error) {
	out, err := iEpochManager.abi.Unpack("getVaultYieldForEpoch", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGrantVaultRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x60698943.
//
// Solidity: function grantVaultRole(address vault) returns()
func (iEpochManager *IEpochManager) PackGrantVaultRole(vault common.Address) []byte {
	enc, err := iEpochManager.abi.Pack("grantVaultRole", vault)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRevokeVaultRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4fe3d970.
//
// Solidity: function revokeVaultRole(address vault) returns()
func (iEpochManager *IEpochManager) PackRevokeVaultRole(vault common.Address) []byte {
	enc, err := iEpochManager.abi.Pack("revokeVaultRole", vault)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetAutomatedSystem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x29f4e5ea.
//
// Solidity: function setAutomatedSystem(address newAutomatedSystem) returns()
func (iEpochManager *IEpochManager) PackSetAutomatedSystem(newAutomatedSystem common.Address) []byte {
	enc, err := iEpochManager.abi.Pack("setAutomatedSystem", newAutomatedSystem)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetDebtSubsidizer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb94305aa.
//
// Solidity: function setDebtSubsidizer(address newDebtSubsidizer) returns()
func (iEpochManager *IEpochManager) PackSetDebtSubsidizer(newDebtSubsidizer common.Address) []byte {
	enc, err := iEpochManager.abi.Pack("setDebtSubsidizer", newDebtSubsidizer)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackStartEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2c8b177.
//
// Solidity: function startEpoch() returns(uint256 epochId)
func (iEpochManager *IEpochManager) PackStartEpoch() []byte {
	enc, err := iEpochManager.abi.Pack("startEpoch")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackStartEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa2c8b177.
//
// Solidity: function startEpoch() returns(uint256 epochId)
func (iEpochManager *IEpochManager) UnpackStartEpoch(data []byte) (*big.Int, error) {
	out, err := iEpochManager.abi.Unpack("startEpoch", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// IEpochManagerAutomatedSystemUpdated represents a AutomatedSystemUpdated event raised by the IEpochManager contract.
type IEpochManagerAutomatedSystemUpdated struct {
	NewAutomatedSystem common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const IEpochManagerAutomatedSystemUpdatedEventName = "AutomatedSystemUpdated"

// ContractEventName returns the user-defined event name.
func (IEpochManagerAutomatedSystemUpdated) ContractEventName() string {
	return IEpochManagerAutomatedSystemUpdatedEventName
}

// UnpackAutomatedSystemUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AutomatedSystemUpdated(address indexed newAutomatedSystem)
func (iEpochManager *IEpochManager) UnpackAutomatedSystemUpdatedEvent(log *types.Log) (*IEpochManagerAutomatedSystemUpdated, error) {
	event := "AutomatedSystemUpdated"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerAutomatedSystemUpdated)
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

// IEpochManagerDebtSubsidizerUpdated represents a DebtSubsidizerUpdated event raised by the IEpochManager contract.
type IEpochManagerDebtSubsidizerUpdated struct {
	NewDebtSubsidizer common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const IEpochManagerDebtSubsidizerUpdatedEventName = "DebtSubsidizerUpdated"

// ContractEventName returns the user-defined event name.
func (IEpochManagerDebtSubsidizerUpdated) ContractEventName() string {
	return IEpochManagerDebtSubsidizerUpdatedEventName
}

// UnpackDebtSubsidizerUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DebtSubsidizerUpdated(address indexed newDebtSubsidizer)
func (iEpochManager *IEpochManager) UnpackDebtSubsidizerUpdatedEvent(log *types.Log) (*IEpochManagerDebtSubsidizerUpdated, error) {
	event := "DebtSubsidizerUpdated"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerDebtSubsidizerUpdated)
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

// IEpochManagerEpochFailed represents a EpochFailed event raised by the IEpochManager contract.
type IEpochManagerEpochFailed struct {
	EpochId *big.Int
	Reason  string
	Raw     *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochFailedEventName = "EpochFailed"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochFailed) ContractEventName() string {
	return IEpochManagerEpochFailedEventName
}

// UnpackEpochFailedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochFailed(uint256 indexed epochId, string reason)
func (iEpochManager *IEpochManager) UnpackEpochFailedEvent(log *types.Log) (*IEpochManagerEpochFailed, error) {
	event := "EpochFailed"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochFailed)
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

// IEpochManagerEpochFinalized represents a EpochFinalized event raised by the IEpochManager contract.
type IEpochManagerEpochFinalized struct {
	EpochId                   *big.Int
	TotalYieldAvailable       *big.Int
	TotalSubsidiesDistributed *big.Int
	Raw                       *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochFinalizedEventName = "EpochFinalized"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochFinalized) ContractEventName() string {
	return IEpochManagerEpochFinalizedEventName
}

// UnpackEpochFinalizedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochFinalized(uint256 indexed epochId, uint256 totalYieldAvailable, uint256 totalSubsidiesDistributed)
func (iEpochManager *IEpochManager) UnpackEpochFinalizedEvent(log *types.Log) (*IEpochManagerEpochFinalized, error) {
	event := "EpochFinalized"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochFinalized)
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

// IEpochManagerEpochProcessingStarted represents a EpochProcessingStarted event raised by the IEpochManager contract.
type IEpochManagerEpochProcessingStarted struct {
	EpochId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochProcessingStartedEventName = "EpochProcessingStarted"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochProcessingStarted) ContractEventName() string {
	return IEpochManagerEpochProcessingStartedEventName
}

// UnpackEpochProcessingStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochProcessingStarted(uint256 indexed epochId)
func (iEpochManager *IEpochManager) UnpackEpochProcessingStartedEvent(log *types.Log) (*IEpochManagerEpochProcessingStarted, error) {
	event := "EpochProcessingStarted"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochProcessingStarted)
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

// IEpochManagerEpochStarted represents a EpochStarted event raised by the IEpochManager contract.
type IEpochManagerEpochStarted struct {
	EpochId   *big.Int
	StartTime *big.Int
	EndTime   *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const IEpochManagerEpochStartedEventName = "EpochStarted"

// ContractEventName returns the user-defined event name.
func (IEpochManagerEpochStarted) ContractEventName() string {
	return IEpochManagerEpochStartedEventName
}

// UnpackEpochStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochStarted(uint256 indexed epochId, uint256 startTime, uint256 endTime)
func (iEpochManager *IEpochManager) UnpackEpochStartedEvent(log *types.Log) (*IEpochManagerEpochStarted, error) {
	event := "EpochStarted"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerEpochStarted)
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

// IEpochManagerProcessingFailed represents a ProcessingFailed event raised by the IEpochManager contract.
type IEpochManagerProcessingFailed struct {
	EpochId *big.Int
	Reason  string
	Raw     *types.Log // Blockchain specific contextual infos
}

const IEpochManagerProcessingFailedEventName = "ProcessingFailed"

// ContractEventName returns the user-defined event name.
func (IEpochManagerProcessingFailed) ContractEventName() string {
	return IEpochManagerProcessingFailedEventName
}

// UnpackProcessingFailedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ProcessingFailed(uint256 indexed epochId, string reason)
func (iEpochManager *IEpochManager) UnpackProcessingFailedEvent(log *types.Log) (*IEpochManagerProcessingFailed, error) {
	event := "ProcessingFailed"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerProcessingFailed)
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

// IEpochManagerProcessingStarted represents a ProcessingStarted event raised by the IEpochManager contract.
type IEpochManagerProcessingStarted struct {
	EpochId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const IEpochManagerProcessingStartedEventName = "ProcessingStarted"

// ContractEventName returns the user-defined event name.
func (IEpochManagerProcessingStarted) ContractEventName() string {
	return IEpochManagerProcessingStartedEventName
}

// UnpackProcessingStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ProcessingStarted(uint256 indexed epochId)
func (iEpochManager *IEpochManager) UnpackProcessingStartedEvent(log *types.Log) (*IEpochManagerProcessingStarted, error) {
	event := "ProcessingStarted"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerProcessingStarted)
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

// IEpochManagerVaultYieldAllocated represents a VaultYieldAllocated event raised by the IEpochManager contract.
type IEpochManagerVaultYieldAllocated struct {
	EpochId *big.Int
	Vault   common.Address
	Amount  *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const IEpochManagerVaultYieldAllocatedEventName = "VaultYieldAllocated"

// ContractEventName returns the user-defined event name.
func (IEpochManagerVaultYieldAllocated) ContractEventName() string {
	return IEpochManagerVaultYieldAllocatedEventName
}

// UnpackVaultYieldAllocatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultYieldAllocated(uint256 indexed epochId, address indexed vault, uint256 amount)
func (iEpochManager *IEpochManager) UnpackVaultYieldAllocatedEvent(log *types.Log) (*IEpochManagerVaultYieldAllocated, error) {
	event := "VaultYieldAllocated"
	if log.Topics[0] != iEpochManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEpochManagerVaultYieldAllocated)
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

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (iEpochManager *IEpochManager) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iEpochManager.abi.Errors["EpochManagerEpochNotEnded"].ID.Bytes()[:4]) {
		return iEpochManager.UnpackEpochManagerEpochNotEndedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEpochManager.abi.Errors["EpochManagerEpochStillActive"].ID.Bytes()[:4]) {
		return iEpochManager.UnpackEpochManagerEpochStillActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEpochManager.abi.Errors["EpochManagerInvalidEpochDuration"].ID.Bytes()[:4]) {
		return iEpochManager.UnpackEpochManagerInvalidEpochDurationError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEpochManager.abi.Errors["EpochManagerInvalidEpochId"].ID.Bytes()[:4]) {
		return iEpochManager.UnpackEpochManagerInvalidEpochIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEpochManager.abi.Errors["EpochManagerInvalidEpochStatus"].ID.Bytes()[:4]) {
		return iEpochManager.UnpackEpochManagerInvalidEpochStatusError(raw[4:])
	}
	if bytes.Equal(raw[:4], iEpochManager.abi.Errors["EpochManagerUnauthorized"].ID.Bytes()[:4]) {
		return iEpochManager.UnpackEpochManagerUnauthorizedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IEpochManagerEpochManagerEpochNotEnded represents a EpochManager__EpochNotEnded error raised by the IEpochManager contract.
type IEpochManagerEpochManagerEpochNotEnded struct {
	EpochId *big.Int
	EndTime *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManager__EpochNotEnded(uint256 epochId, uint256 endTime)
func IEpochManagerEpochManagerEpochNotEndedErrorID() common.Hash {
	return common.HexToHash("0xc009067a49e9244045f824f460d5c9e37aa8931e947e569627d62abde4046ce7")
}

// UnpackEpochManagerEpochNotEndedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManager__EpochNotEnded(uint256 epochId, uint256 endTime)
func (iEpochManager *IEpochManager) UnpackEpochManagerEpochNotEndedError(raw []byte) (*IEpochManagerEpochManagerEpochNotEnded, error) {
	out := new(IEpochManagerEpochManagerEpochNotEnded)
	if err := iEpochManager.abi.UnpackIntoInterface(out, "EpochManagerEpochNotEnded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEpochManagerEpochManagerEpochStillActive represents a EpochManager__EpochStillActive error raised by the IEpochManager contract.
type IEpochManagerEpochManagerEpochStillActive struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManager__EpochStillActive()
func IEpochManagerEpochManagerEpochStillActiveErrorID() common.Hash {
	return common.HexToHash("0x077ec33b943c16accc85fe49cc56a01d22af6dfd9231ddd8fb552c2b376cbad4")
}

// UnpackEpochManagerEpochStillActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManager__EpochStillActive()
func (iEpochManager *IEpochManager) UnpackEpochManagerEpochStillActiveError(raw []byte) (*IEpochManagerEpochManagerEpochStillActive, error) {
	out := new(IEpochManagerEpochManagerEpochStillActive)
	if err := iEpochManager.abi.UnpackIntoInterface(out, "EpochManagerEpochStillActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEpochManagerEpochManagerInvalidEpochDuration represents a EpochManager__InvalidEpochDuration error raised by the IEpochManager contract.
type IEpochManagerEpochManagerInvalidEpochDuration struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManager__InvalidEpochDuration()
func IEpochManagerEpochManagerInvalidEpochDurationErrorID() common.Hash {
	return common.HexToHash("0xaa74abf80f0b84904bb802d0dcafc6255baef7e480abdc4a8ff730fca2ec1b9f")
}

// UnpackEpochManagerInvalidEpochDurationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManager__InvalidEpochDuration()
func (iEpochManager *IEpochManager) UnpackEpochManagerInvalidEpochDurationError(raw []byte) (*IEpochManagerEpochManagerInvalidEpochDuration, error) {
	out := new(IEpochManagerEpochManagerInvalidEpochDuration)
	if err := iEpochManager.abi.UnpackIntoInterface(out, "EpochManagerInvalidEpochDuration", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEpochManagerEpochManagerInvalidEpochId represents a EpochManager__InvalidEpochId error raised by the IEpochManager contract.
type IEpochManagerEpochManagerInvalidEpochId struct {
	EpochId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManager__InvalidEpochId(uint256 epochId)
func IEpochManagerEpochManagerInvalidEpochIdErrorID() common.Hash {
	return common.HexToHash("0xbdd09f9d4371ddcbd142b10ccccca9e84423a3217f301f28b9c796fc55840548")
}

// UnpackEpochManagerInvalidEpochIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManager__InvalidEpochId(uint256 epochId)
func (iEpochManager *IEpochManager) UnpackEpochManagerInvalidEpochIdError(raw []byte) (*IEpochManagerEpochManagerInvalidEpochId, error) {
	out := new(IEpochManagerEpochManagerInvalidEpochId)
	if err := iEpochManager.abi.UnpackIntoInterface(out, "EpochManagerInvalidEpochId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEpochManagerEpochManagerInvalidEpochStatus represents a EpochManager__InvalidEpochStatus error raised by the IEpochManager contract.
type IEpochManagerEpochManagerInvalidEpochStatus struct {
	EpochId        *big.Int
	CurrentStatus  uint8
	ExpectedStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManager__InvalidEpochStatus(uint256 epochId, uint8 currentStatus, uint8 expectedStatus)
func IEpochManagerEpochManagerInvalidEpochStatusErrorID() common.Hash {
	return common.HexToHash("0xda6bfe957aeff47e28b2cfd09fb355c612a32c9f853b7cc93c4354c99f066077")
}

// UnpackEpochManagerInvalidEpochStatusError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManager__InvalidEpochStatus(uint256 epochId, uint8 currentStatus, uint8 expectedStatus)
func (iEpochManager *IEpochManager) UnpackEpochManagerInvalidEpochStatusError(raw []byte) (*IEpochManagerEpochManagerInvalidEpochStatus, error) {
	out := new(IEpochManagerEpochManagerInvalidEpochStatus)
	if err := iEpochManager.abi.UnpackIntoInterface(out, "EpochManagerInvalidEpochStatus", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IEpochManagerEpochManagerUnauthorized represents a EpochManager__Unauthorized error raised by the IEpochManager contract.
type IEpochManagerEpochManagerUnauthorized struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManager__Unauthorized()
func IEpochManagerEpochManagerUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0x29b1e89e9302a5b7f3f67df88766b30c0c14ca4186069116a6f98517ca4cd157")
}

// UnpackEpochManagerUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManager__Unauthorized()
func (iEpochManager *IEpochManager) UnpackEpochManagerUnauthorizedError(raw []byte) (*IEpochManagerEpochManagerUnauthorized, error) {
	out := new(IEpochManagerEpochManagerUnauthorized)
	if err := iEpochManager.abi.UnpackIntoInterface(out, "EpochManagerUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}
