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

// ICollectionRegistryCollection is an auto generated low-level Go binding around an user-defined struct.
type ICollectionRegistryCollection struct {
	CollectionAddress    common.Address
	CollectionType       uint8
	WeightFunction       ICollectionRegistryWeightFunction
	P1                   *big.Int
	P2                   *big.Int
	YieldSharePercentage uint16
	Vaults               []common.Address
}

// ICollectionRegistryWeightFunction is an auto generated low-level Go binding around an user-defined struct.
type ICollectionRegistryWeightFunction struct {
	FnType uint8
	P1     *big.Int
	P2     *big.Int
}

// ICollectionRegistryMetaData contains all meta data concerning the ICollectionRegistry contract.
var ICollectionRegistryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addVaultToCollection\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allCollections\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCollection\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structICollectionRegistry.Collection\",\"components\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionType\",\"type\":\"uint8\",\"internalType\":\"enumICollectionRegistry.CollectionType\"},{\"name\":\"weightFunction\",\"type\":\"tuple\",\"internalType\":\"structICollectionRegistry.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumICollectionRegistry.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"yieldSharePercentage\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerCollection\",\"inputs\":[{\"name\":\"collection\",\"type\":\"tuple\",\"internalType\":\"structICollectionRegistry.Collection\",\"components\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionType\",\"type\":\"uint8\",\"internalType\":\"enumICollectionRegistry.CollectionType\"},{\"name\":\"weightFunction\",\"type\":\"tuple\",\"internalType\":\"structICollectionRegistry.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumICollectionRegistry.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"yieldSharePercentage\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeVaultFromCollection\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWeightFunction\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"weightFunction\",\"type\":\"tuple\",\"internalType\":\"structICollectionRegistry.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumICollectionRegistry.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setYieldShare\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"share\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	ID:  "ICollectionRegistry",
}

// ICollectionRegistry is an auto generated Go binding around an Ethereum contract.
type ICollectionRegistry struct {
	abi abi.ABI
}

// NewICollectionRegistry creates a new instance of ICollectionRegistry.
func NewICollectionRegistry() *ICollectionRegistry {
	parsed, err := ICollectionRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ICollectionRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ICollectionRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddVaultToCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1a68399f.
//
// Solidity: function addVaultToCollection(address collection, address vault) returns()
func (iCollectionRegistry *ICollectionRegistry) PackAddVaultToCollection(collection common.Address, vault common.Address) []byte {
	enc, err := iCollectionRegistry.abi.Pack("addVaultToCollection", collection, vault)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAllCollections is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe40d4544.
//
// Solidity: function allCollections() view returns(address[])
func (iCollectionRegistry *ICollectionRegistry) PackAllCollections() []byte {
	enc, err := iCollectionRegistry.abi.Pack("allCollections")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAllCollections is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe40d4544.
//
// Solidity: function allCollections() view returns(address[])
func (iCollectionRegistry *ICollectionRegistry) UnpackAllCollections(data []byte) ([]common.Address, error) {
	out, err := iCollectionRegistry.abi.Unpack("allCollections", data)
	if err != nil {
		return *new([]common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, err
}

// PackGetCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe40de887.
//
// Solidity: function getCollection(address collection) view returns((address,uint8,(uint8,int256,int256),int256,int256,uint16,address[]))
func (iCollectionRegistry *ICollectionRegistry) PackGetCollection(collection common.Address) []byte {
	enc, err := iCollectionRegistry.abi.Pack("getCollection", collection)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCollection is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe40de887.
//
// Solidity: function getCollection(address collection) view returns((address,uint8,(uint8,int256,int256),int256,int256,uint16,address[]))
func (iCollectionRegistry *ICollectionRegistry) UnpackGetCollection(data []byte) (ICollectionRegistryCollection, error) {
	out, err := iCollectionRegistry.abi.Unpack("getCollection", data)
	if err != nil {
		return *new(ICollectionRegistryCollection), err
	}
	out0 := *abi.ConvertType(out[0], new(ICollectionRegistryCollection)).(*ICollectionRegistryCollection)
	return out0, err
}

// PackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3c5a547.
//
// Solidity: function isRegistered(address collection) view returns(bool)
func (iCollectionRegistry *ICollectionRegistry) PackIsRegistered(collection common.Address) []byte {
	enc, err := iCollectionRegistry.abi.Pack("isRegistered", collection)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsRegistered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3c5a547.
//
// Solidity: function isRegistered(address collection) view returns(bool)
func (iCollectionRegistry *ICollectionRegistry) UnpackIsRegistered(data []byte) (bool, error) {
	out, err := iCollectionRegistry.abi.Unpack("isRegistered", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRegisterCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf9b8794.
//
// Solidity: function registerCollection((address,uint8,(uint8,int256,int256),int256,int256,uint16,address[]) collection) returns()
func (iCollectionRegistry *ICollectionRegistry) PackRegisterCollection(collection ICollectionRegistryCollection) []byte {
	enc, err := iCollectionRegistry.abi.Pack("registerCollection", collection)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveVaultFromCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x189f2859.
//
// Solidity: function removeVaultFromCollection(address collection, address vault) returns()
func (iCollectionRegistry *ICollectionRegistry) PackRemoveVaultFromCollection(collection common.Address, vault common.Address) []byte {
	enc, err := iCollectionRegistry.abi.Pack("removeVaultFromCollection", collection, vault)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetWeightFunction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0546fb17.
//
// Solidity: function setWeightFunction(address collection, (uint8,int256,int256) weightFunction, int256 p1, int256 p2) returns()
func (iCollectionRegistry *ICollectionRegistry) PackSetWeightFunction(collection common.Address, weightFunction ICollectionRegistryWeightFunction, p1 *big.Int, p2 *big.Int) []byte {
	enc, err := iCollectionRegistry.abi.Pack("setWeightFunction", collection, weightFunction, p1, p2)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetYieldShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4ebd6611.
//
// Solidity: function setYieldShare(address collection, uint16 share) returns()
func (iCollectionRegistry *ICollectionRegistry) PackSetYieldShare(collection common.Address, share uint16) []byte {
	enc, err := iCollectionRegistry.abi.Pack("setYieldShare", collection, share)
	if err != nil {
		panic(err)
	}
	return enc
}
