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

// IDebtSubsidizerClaimData is an auto generated low-level Go binding around an user-defined struct.
type IDebtSubsidizerClaimData struct {
	Recipient   common.Address
	TotalEarned *big.Int
	MerkleProof [][32]byte
}

// IDebtSubsidizerVaultInfo is an auto generated low-level Go binding around an user-defined struct.
type IDebtSubsidizerVaultInfo struct {
	LendingManager common.Address
	CToken         common.Address
}

// IDebtSubsidizerMetaData contains all meta data concerning the IDebtSubsidizer contract.
var IDebtSubsidizerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addVault\",\"inputs\":[{\"name\":\"vaultAddress_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lendingManagerAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimAllSubsidies\",\"inputs\":[{\"name\":\"vaultAddresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"claims\",\"type\":\"tuple[]\",\"internalType\":\"structIDebtSubsidizer.ClaimData[]\",\"components\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"totalEarned\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimSubsidy\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"claim\",\"type\":\"tuple\",\"internalType\":\"structIDebtSubsidizer.ClaimData\",\"components\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"totalEarned\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emergencyValidateAndPause\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getMerkleRoot\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemainingSubsidies\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalClaimedForVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalSubsidies\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalSubsidiesClaimed\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserClaimedTotal\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isCollectionRemoved\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isCollectionWhitelisted\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isVaultRemoved\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeCollection\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeVault\",\"inputs\":[{\"name\":\"vaultAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMerkleRoot\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"totalSubsidiesForEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userSecondsClaimed\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validateVaultClaimsIntegrity\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDebtSubsidizer.VaultInfo\",\"components\":[{\"name\":\"lendingManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"cToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"whitelistCollection\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CollectionRegistryUpdated\",\"inputs\":[{\"name\":\"oldRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DebtSubsidizerRoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DebtSubsidizerRoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MerkleRootUpdated\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"updatedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalSubsidiesForEpoch\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewCollectionWhitelisted\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubsidyClaimed\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAdded\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"cTokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"lendingManagerAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRemoved\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistCollectionRemoved\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotSetSignerToZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CollectionAlreadyExists\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CollectionAlreadyWhitelistedInVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CollectionNotWhitelisted\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CollectionNotWhitelistedInVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientYield\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCollectionInterface\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"InvalidMerkleProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecondsColl\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidYieldSharePercentage\",\"inputs\":[{\"name\":\"totalSharePercentage\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidYieldSlice\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LendingManagerAssetMismatch\",\"inputs\":[{\"name\":\"vaultAsset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lmAsset\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LendingManagerNotSetForVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MerkleRootNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"VaultAlreadyRegistered\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"VaultNotRegistered\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "IDebtSubsidizer",
}

// IDebtSubsidizer is an auto generated Go binding around an Ethereum contract.
type IDebtSubsidizer struct {
	abi abi.ABI
}

// NewIDebtSubsidizer creates a new instance of IDebtSubsidizer.
func NewIDebtSubsidizer() *IDebtSubsidizer {
	parsed, err := IDebtSubsidizerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IDebtSubsidizer{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IDebtSubsidizer) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xec3a7823.
//
// Solidity: function addVault(address vaultAddress_, address lendingManagerAddress_) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackAddVault(vaultAddress common.Address, lendingManagerAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("addVault", vaultAddress, lendingManagerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackClaimAllSubsidies is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfd8f92ae.
//
// Solidity: function claimAllSubsidies(address[] vaultAddresses, (address,uint256,bytes32[])[] claims) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackClaimAllSubsidies(vaultAddresses []common.Address, claims []IDebtSubsidizerClaimData) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("claimAllSubsidies", vaultAddresses, claims)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackClaimSubsidy is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf314ce2d.
//
// Solidity: function claimSubsidy(address vaultAddress, (address,uint256,bytes32[]) claim) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackClaimSubsidy(vaultAddress common.Address, claim IDebtSubsidizerClaimData) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("claimSubsidy", vaultAddress, claim)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEmergencyValidateAndPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x141d2349.
//
// Solidity: function emergencyValidateAndPause(address vaultAddress) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackEmergencyValidateAndPause(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("emergencyValidateAndPause", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetMerkleRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe0536119.
//
// Solidity: function getMerkleRoot(address vaultAddress) view returns(bytes32)
func (iDebtSubsidizer *IDebtSubsidizer) PackGetMerkleRoot(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("getMerkleRoot", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetMerkleRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe0536119.
//
// Solidity: function getMerkleRoot(address vaultAddress) view returns(bytes32)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackGetMerkleRoot(data []byte) ([32]byte, error) {
	out, err := iDebtSubsidizer.abi.Unpack("getMerkleRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackGetRemainingSubsidies is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a3390eb.
//
// Solidity: function getRemainingSubsidies(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackGetRemainingSubsidies(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("getRemainingSubsidies", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetRemainingSubsidies is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a3390eb.
//
// Solidity: function getRemainingSubsidies(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackGetRemainingSubsidies(data []byte) (*big.Int, error) {
	out, err := iDebtSubsidizer.abi.Unpack("getRemainingSubsidies", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalClaimedForVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a494a9f.
//
// Solidity: function getTotalClaimedForVault(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackGetTotalClaimedForVault(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("getTotalClaimedForVault", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalClaimedForVault is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a494a9f.
//
// Solidity: function getTotalClaimedForVault(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackGetTotalClaimedForVault(data []byte) (*big.Int, error) {
	out, err := iDebtSubsidizer.abi.Unpack("getTotalClaimedForVault", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalSubsidies is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb844f62a.
//
// Solidity: function getTotalSubsidies(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackGetTotalSubsidies(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("getTotalSubsidies", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalSubsidies is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb844f62a.
//
// Solidity: function getTotalSubsidies(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackGetTotalSubsidies(data []byte) (*big.Int, error) {
	out, err := iDebtSubsidizer.abi.Unpack("getTotalSubsidies", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalSubsidiesClaimed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x44dea2e7.
//
// Solidity: function getTotalSubsidiesClaimed(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackGetTotalSubsidiesClaimed(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("getTotalSubsidiesClaimed", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalSubsidiesClaimed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x44dea2e7.
//
// Solidity: function getTotalSubsidiesClaimed(address vaultAddress) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackGetTotalSubsidiesClaimed(data []byte) (*big.Int, error) {
	out, err := iDebtSubsidizer.abi.Unpack("getTotalSubsidiesClaimed", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetUserClaimedTotal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e269b4f.
//
// Solidity: function getUserClaimedTotal(address vaultAddress, address user) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackGetUserClaimedTotal(vaultAddress common.Address, user common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("getUserClaimedTotal", vaultAddress, user)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetUserClaimedTotal is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2e269b4f.
//
// Solidity: function getUserClaimedTotal(address vaultAddress, address user) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackGetUserClaimedTotal(data []byte) (*big.Int, error) {
	out, err := iDebtSubsidizer.abi.Unpack("getUserClaimedTotal", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackIsCollectionRemoved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb62d99a0.
//
// Solidity: function isCollectionRemoved(address vaultAddress, address collection) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) PackIsCollectionRemoved(vaultAddress common.Address, collection common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("isCollectionRemoved", vaultAddress, collection)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsCollectionRemoved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb62d99a0.
//
// Solidity: function isCollectionRemoved(address vaultAddress, address collection) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackIsCollectionRemoved(data []byte) (bool, error) {
	out, err := iDebtSubsidizer.abi.Unpack("isCollectionRemoved", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsCollectionWhitelisted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x41afb808.
//
// Solidity: function isCollectionWhitelisted(address vaultAddress, address collectionAddress) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) PackIsCollectionWhitelisted(vaultAddress common.Address, collectionAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("isCollectionWhitelisted", vaultAddress, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsCollectionWhitelisted is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x41afb808.
//
// Solidity: function isCollectionWhitelisted(address vaultAddress, address collectionAddress) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackIsCollectionWhitelisted(data []byte) (bool, error) {
	out, err := iDebtSubsidizer.abi.Unpack("isCollectionWhitelisted", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsVaultRemoved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc21ba7e4.
//
// Solidity: function isVaultRemoved(address vaultAddress) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) PackIsVaultRemoved(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("isVaultRemoved", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsVaultRemoved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc21ba7e4.
//
// Solidity: function isVaultRemoved(address vaultAddress) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackIsVaultRemoved(data []byte) (bool, error) {
	out, err := iDebtSubsidizer.abi.Unpack("isVaultRemoved", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) PackPaused() []byte {
	enc, err := iDebtSubsidizer.abi.Pack("paused")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackPaused(data []byte) (bool, error) {
	out, err := iDebtSubsidizer.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRemoveCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f1bcbad.
//
// Solidity: function removeCollection(address vaultAddress, address collectionAddress) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackRemoveCollection(vaultAddress common.Address, collectionAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("removeCollection", vaultAddress, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xceb68c23.
//
// Solidity: function removeVault(address vaultAddress_) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackRemoveVault(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("removeVault", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateMerkleRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0f4bcfb.
//
// Solidity: function updateMerkleRoot(address vaultAddress, bytes32 merkleRoot, uint256 totalSubsidiesForEpoch) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackUpdateMerkleRoot(vaultAddress common.Address, merkleRoot [32]byte, totalSubsidiesForEpoch *big.Int) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("updateMerkleRoot", vaultAddress, merkleRoot, totalSubsidiesForEpoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUserSecondsClaimed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1cd60f2.
//
// Solidity: function userSecondsClaimed(address user) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackUserSecondsClaimed(user common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("userSecondsClaimed", user)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUserSecondsClaimed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb1cd60f2.
//
// Solidity: function userSecondsClaimed(address user) view returns(uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackUserSecondsClaimed(data []byte) (*big.Int, error) {
	out, err := iDebtSubsidizer.abi.Unpack("userSecondsClaimed", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackValidateVaultClaimsIntegrity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x560dff1d.
//
// Solidity: function validateVaultClaimsIntegrity(address vaultAddress) view returns(bool, uint256, uint256)
func (iDebtSubsidizer *IDebtSubsidizer) PackValidateVaultClaimsIntegrity(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("validateVaultClaimsIntegrity", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// ValidateVaultClaimsIntegrityOutput serves as a container for the return parameters of contract
// method ValidateVaultClaimsIntegrity.
type ValidateVaultClaimsIntegrityOutput struct {
	Arg0 bool
	Arg1 *big.Int
	Arg2 *big.Int
}

// UnpackValidateVaultClaimsIntegrity is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x560dff1d.
//
// Solidity: function validateVaultClaimsIntegrity(address vaultAddress) view returns(bool, uint256, uint256)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackValidateVaultClaimsIntegrity(data []byte) (ValidateVaultClaimsIntegrityOutput, error) {
	out, err := iDebtSubsidizer.abi.Unpack("validateVaultClaimsIntegrity", data)
	outstruct := new(ValidateVaultClaimsIntegrityOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Arg2 = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf815c03d.
//
// Solidity: function vault(address vaultAddress) view returns((address,address))
func (iDebtSubsidizer *IDebtSubsidizer) PackVault(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("vault", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVault is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf815c03d.
//
// Solidity: function vault(address vaultAddress) view returns((address,address))
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVault(data []byte) (IDebtSubsidizerVaultInfo, error) {
	out, err := iDebtSubsidizer.abi.Unpack("vault", data)
	if err != nil {
		return *new(IDebtSubsidizerVaultInfo), err
	}
	out0 := *abi.ConvertType(out[0], new(IDebtSubsidizerVaultInfo)).(*IDebtSubsidizerVaultInfo)
	return out0, err
}

// PackWhitelistCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2526952a.
//
// Solidity: function whitelistCollection(address vaultAddress, address collectionAddress) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackWhitelistCollection(vaultAddress common.Address, collectionAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("whitelistCollection", vaultAddress, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// IDebtSubsidizerCollectionRegistryUpdated represents a CollectionRegistryUpdated event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionRegistryUpdated struct {
	OldRegistry common.Address
	NewRegistry common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerCollectionRegistryUpdatedEventName = "CollectionRegistryUpdated"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerCollectionRegistryUpdated) ContractEventName() string {
	return IDebtSubsidizerCollectionRegistryUpdatedEventName
}

// UnpackCollectionRegistryUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionRegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionRegistryUpdatedEvent(log *types.Log) (*IDebtSubsidizerCollectionRegistryUpdated, error) {
	event := "CollectionRegistryUpdated"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerCollectionRegistryUpdated)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerDebtSubsidizerRoleGranted represents a DebtSubsidizerRoleGranted event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerDebtSubsidizerRoleGranted struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerDebtSubsidizerRoleGrantedEventName = "DebtSubsidizerRoleGranted"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerDebtSubsidizerRoleGranted) ContractEventName() string {
	return IDebtSubsidizerDebtSubsidizerRoleGrantedEventName
}

// UnpackDebtSubsidizerRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DebtSubsidizerRoleGranted(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackDebtSubsidizerRoleGrantedEvent(log *types.Log) (*IDebtSubsidizerDebtSubsidizerRoleGranted, error) {
	event := "DebtSubsidizerRoleGranted"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerDebtSubsidizerRoleGranted)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerDebtSubsidizerRoleRevoked represents a DebtSubsidizerRoleRevoked event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerDebtSubsidizerRoleRevoked struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerDebtSubsidizerRoleRevokedEventName = "DebtSubsidizerRoleRevoked"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerDebtSubsidizerRoleRevoked) ContractEventName() string {
	return IDebtSubsidizerDebtSubsidizerRoleRevokedEventName
}

// UnpackDebtSubsidizerRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DebtSubsidizerRoleRevoked(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackDebtSubsidizerRoleRevokedEvent(log *types.Log) (*IDebtSubsidizerDebtSubsidizerRoleRevoked, error) {
	event := "DebtSubsidizerRoleRevoked"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerDebtSubsidizerRoleRevoked)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerMerkleRootUpdated represents a MerkleRootUpdated event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerMerkleRootUpdated struct {
	VaultAddress           common.Address
	MerkleRoot             [32]byte
	UpdatedBy              common.Address
	TotalSubsidiesForEpoch *big.Int
	Raw                    *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerMerkleRootUpdatedEventName = "MerkleRootUpdated"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerMerkleRootUpdated) ContractEventName() string {
	return IDebtSubsidizerMerkleRootUpdatedEventName
}

// UnpackMerkleRootUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MerkleRootUpdated(address indexed vaultAddress, bytes32 merkleRoot, address indexed updatedBy, uint256 totalSubsidiesForEpoch)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackMerkleRootUpdatedEvent(log *types.Log) (*IDebtSubsidizerMerkleRootUpdated, error) {
	event := "MerkleRootUpdated"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerMerkleRootUpdated)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerNewCollectionWhitelisted represents a NewCollectionWhitelisted event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerNewCollectionWhitelisted struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerNewCollectionWhitelistedEventName = "NewCollectionWhitelisted"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerNewCollectionWhitelisted) ContractEventName() string {
	return IDebtSubsidizerNewCollectionWhitelistedEventName
}

// UnpackNewCollectionWhitelistedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event NewCollectionWhitelisted(address indexed vaultAddress, address indexed collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackNewCollectionWhitelistedEvent(log *types.Log) (*IDebtSubsidizerNewCollectionWhitelisted, error) {
	event := "NewCollectionWhitelisted"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerNewCollectionWhitelisted)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerSubsidyClaimed represents a SubsidyClaimed event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerSubsidyClaimed struct {
	VaultAddress common.Address
	Recipient    common.Address
	Amount       *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerSubsidyClaimedEventName = "SubsidyClaimed"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerSubsidyClaimed) ContractEventName() string {
	return IDebtSubsidizerSubsidyClaimedEventName
}

// UnpackSubsidyClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SubsidyClaimed(address indexed vaultAddress, address indexed recipient, uint256 amount)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackSubsidyClaimedEvent(log *types.Log) (*IDebtSubsidizerSubsidyClaimed, error) {
	event := "SubsidyClaimed"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerSubsidyClaimed)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerVaultAdded represents a VaultAdded event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultAdded struct {
	VaultAddress          common.Address
	CTokenAddress         common.Address
	LendingManagerAddress common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerVaultAddedEventName = "VaultAdded"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerVaultAdded) ContractEventName() string {
	return IDebtSubsidizerVaultAddedEventName
}

// UnpackVaultAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultAdded(address indexed vaultAddress, address indexed cTokenAddress, address indexed lendingManagerAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultAddedEvent(log *types.Log) (*IDebtSubsidizerVaultAdded, error) {
	event := "VaultAdded"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerVaultAdded)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerVaultRemoved represents a VaultRemoved event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultRemoved struct {
	VaultAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerVaultRemovedEventName = "VaultRemoved"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerVaultRemoved) ContractEventName() string {
	return IDebtSubsidizerVaultRemovedEventName
}

// UnpackVaultRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultRemoved(address indexed vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultRemovedEvent(log *types.Log) (*IDebtSubsidizerVaultRemoved, error) {
	event := "VaultRemoved"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerVaultRemoved)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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

// IDebtSubsidizerWhitelistCollectionRemoved represents a WhitelistCollectionRemoved event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerWhitelistCollectionRemoved struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerWhitelistCollectionRemovedEventName = "WhitelistCollectionRemoved"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerWhitelistCollectionRemoved) ContractEventName() string {
	return IDebtSubsidizerWhitelistCollectionRemovedEventName
}

// UnpackWhitelistCollectionRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WhitelistCollectionRemoved(address indexed vaultAddress, address indexed collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackWhitelistCollectionRemovedEvent(log *types.Log) (*IDebtSubsidizerWhitelistCollectionRemoved, error) {
	event := "WhitelistCollectionRemoved"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerWhitelistCollectionRemoved)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
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
func (iDebtSubsidizer *IDebtSubsidizer) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["AddressZero"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackAddressZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["AlreadyClaimed"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackAlreadyClaimedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["ArrayLengthMismatch"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackArrayLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CannotSetSignerToZeroAddress"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCannotSetSignerToZeroAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionAlreadyExists"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionAlreadyWhitelistedInVault"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionAlreadyWhitelistedInVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionNotWhitelisted"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionNotWhitelistedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionNotWhitelistedInVault"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionNotWhitelistedInVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InsufficientYield"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInsufficientYieldError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidCollectionInterface"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidCollectionInterfaceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidMerkleProof"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidMerkleProofError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidSecondsColl"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidSecondsCollError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidYieldSharePercentage"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidYieldSharePercentageError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidYieldSlice"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidYieldSliceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["LendingManagerAssetMismatch"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackLendingManagerAssetMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["LendingManagerNotSetForVault"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackLendingManagerNotSetForVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["MerkleRootNotSet"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackMerkleRootNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["VaultAlreadyRegistered"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackVaultAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["VaultMismatch"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackVaultMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["VaultNotRegistered"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackVaultNotRegisteredError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IDebtSubsidizerAddressZero represents a AddressZero error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerAddressZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressZero()
func IDebtSubsidizerAddressZeroErrorID() common.Hash {
	return common.HexToHash("0x9fabe1c19979afc45ec7efec1bde2c38021c590a0ce42965cf55b3f518197f02")
}

// UnpackAddressZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressZero()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackAddressZeroError(raw []byte) (*IDebtSubsidizerAddressZero, error) {
	out := new(IDebtSubsidizerAddressZero)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "AddressZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerAlreadyClaimed represents a AlreadyClaimed error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerAlreadyClaimed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyClaimed()
func IDebtSubsidizerAlreadyClaimedErrorID() common.Hash {
	return common.HexToHash("0x646cf558a545d59f8a09cbf8a0eb8a9332f1d17834843b20fc8d154839dc46d7")
}

// UnpackAlreadyClaimedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyClaimed()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackAlreadyClaimedError(raw []byte) (*IDebtSubsidizerAlreadyClaimed, error) {
	out := new(IDebtSubsidizerAlreadyClaimed)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "AlreadyClaimed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerArrayLengthMismatch represents a ArrayLengthMismatch error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerArrayLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ArrayLengthMismatch()
func IDebtSubsidizerArrayLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xa24a13a6c9c749fdebc1ced0c54b040f90ec2bad4921a2449a09961f99596abe")
}

// UnpackArrayLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ArrayLengthMismatch()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackArrayLengthMismatchError(raw []byte) (*IDebtSubsidizerArrayLengthMismatch, error) {
	out := new(IDebtSubsidizerArrayLengthMismatch)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "ArrayLengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCannotSetSignerToZeroAddress represents a CannotSetSignerToZeroAddress error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCannotSetSignerToZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotSetSignerToZeroAddress()
func IDebtSubsidizerCannotSetSignerToZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xc4aaac8bb50f06fadf7bb35d9bd7eb2b4ca9c665ea3aaa56621dc042c7b14b3f")
}

// UnpackCannotSetSignerToZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotSetSignerToZeroAddress()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCannotSetSignerToZeroAddressError(raw []byte) (*IDebtSubsidizerCannotSetSignerToZeroAddress, error) {
	out := new(IDebtSubsidizerCannotSetSignerToZeroAddress)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CannotSetSignerToZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionAlreadyExists represents a CollectionAlreadyExists error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionAlreadyExists struct {
	Collection common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionAlreadyExists(address collection)
func IDebtSubsidizerCollectionAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0xf53cf25dd2089419f7b729be10af57f42fce4446592b3e7b63db5e9a7dea1d22")
}

// UnpackCollectionAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionAlreadyExists(address collection)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionAlreadyExistsError(raw []byte) (*IDebtSubsidizerCollectionAlreadyExists, error) {
	out := new(IDebtSubsidizerCollectionAlreadyExists)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionAlreadyWhitelistedInVault represents a CollectionAlreadyWhitelistedInVault error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionAlreadyWhitelistedInVault struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionAlreadyWhitelistedInVault(address vaultAddress, address collectionAddress)
func IDebtSubsidizerCollectionAlreadyWhitelistedInVaultErrorID() common.Hash {
	return common.HexToHash("0x5646f04d83aa0b86b63d2636695294c9fabe78eec635ce30317c8d3f2a7445ee")
}

// UnpackCollectionAlreadyWhitelistedInVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionAlreadyWhitelistedInVault(address vaultAddress, address collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionAlreadyWhitelistedInVaultError(raw []byte) (*IDebtSubsidizerCollectionAlreadyWhitelistedInVault, error) {
	out := new(IDebtSubsidizerCollectionAlreadyWhitelistedInVault)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionAlreadyWhitelistedInVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionNotWhitelisted represents a CollectionNotWhitelisted error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionNotWhitelisted struct {
	Collection common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionNotWhitelisted(address collection)
func IDebtSubsidizerCollectionNotWhitelistedErrorID() common.Hash {
	return common.HexToHash("0xd76f0d44a12a3e65ee33e72c9198ae0d715d644e6916f257dfab83e4fe49142d")
}

// UnpackCollectionNotWhitelistedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionNotWhitelisted(address collection)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionNotWhitelistedError(raw []byte) (*IDebtSubsidizerCollectionNotWhitelisted, error) {
	out := new(IDebtSubsidizerCollectionNotWhitelisted)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionNotWhitelisted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionNotWhitelistedInVault represents a CollectionNotWhitelistedInVault error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionNotWhitelistedInVault struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionNotWhitelistedInVault(address vaultAddress, address collectionAddress)
func IDebtSubsidizerCollectionNotWhitelistedInVaultErrorID() common.Hash {
	return common.HexToHash("0x063b14e113f34053abd31ac480d3cf89b7e139870716916231508db61ef4537f")
}

// UnpackCollectionNotWhitelistedInVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionNotWhitelistedInVault(address vaultAddress, address collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionNotWhitelistedInVaultError(raw []byte) (*IDebtSubsidizerCollectionNotWhitelistedInVault, error) {
	out := new(IDebtSubsidizerCollectionNotWhitelistedInVault)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionNotWhitelistedInVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInsufficientYield represents a InsufficientYield error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInsufficientYield struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientYield()
func IDebtSubsidizerInsufficientYieldErrorID() common.Hash {
	return common.HexToHash("0xfa3505cfeaa6d565a857daf3cb6754e23c09dd9555dc3027e00c0110bf42262e")
}

// UnpackInsufficientYieldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientYield()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInsufficientYieldError(raw []byte) (*IDebtSubsidizerInsufficientYield, error) {
	out := new(IDebtSubsidizerInsufficientYield)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InsufficientYield", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidCollectionInterface represents a InvalidCollectionInterface error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidCollectionInterface struct {
	CollectionAddress common.Address
	InterfaceId       [4]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidCollectionInterface(address collectionAddress, bytes4 interfaceId)
func IDebtSubsidizerInvalidCollectionInterfaceErrorID() common.Hash {
	return common.HexToHash("0x77ec721aca2ae0444bced9a0d7a5f5188a2be2a2efd5a50c57eebbebdf070e76")
}

// UnpackInvalidCollectionInterfaceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidCollectionInterface(address collectionAddress, bytes4 interfaceId)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidCollectionInterfaceError(raw []byte) (*IDebtSubsidizerInvalidCollectionInterface, error) {
	out := new(IDebtSubsidizerInvalidCollectionInterface)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidCollectionInterface", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidMerkleProof represents a InvalidMerkleProof error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidMerkleProof struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidMerkleProof()
func IDebtSubsidizerInvalidMerkleProofErrorID() common.Hash {
	return common.HexToHash("0xb05e92facfa0fd6ba9338977017107202232768a12b21141a91a36a56212ad1e")
}

// UnpackInvalidMerkleProofError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidMerkleProof()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidMerkleProofError(raw []byte) (*IDebtSubsidizerInvalidMerkleProof, error) {
	out := new(IDebtSubsidizerInvalidMerkleProof)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidMerkleProof", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidSecondsColl represents a InvalidSecondsColl error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidSecondsColl struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSecondsColl()
func IDebtSubsidizerInvalidSecondsCollErrorID() common.Hash {
	return common.HexToHash("0x8d50ad91ce4aa4f6a33fbdf3db5deb828c1f770a89b6a63e157eb9d7ae851df9")
}

// UnpackInvalidSecondsCollError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSecondsColl()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidSecondsCollError(raw []byte) (*IDebtSubsidizerInvalidSecondsColl, error) {
	out := new(IDebtSubsidizerInvalidSecondsColl)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidSecondsColl", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidSignature represents a InvalidSignature error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func IDebtSubsidizerInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidSignatureError(raw []byte) (*IDebtSubsidizerInvalidSignature, error) {
	out := new(IDebtSubsidizerInvalidSignature)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidYieldSharePercentage represents a InvalidYieldSharePercentage error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidYieldSharePercentage struct {
	TotalSharePercentage *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidYieldSharePercentage(uint256 totalSharePercentage)
func IDebtSubsidizerInvalidYieldSharePercentageErrorID() common.Hash {
	return common.HexToHash("0xbb366faa5bb0b720039d1f82c07bbb07db4fa7da184ae4181de9d51a39d893d4")
}

// UnpackInvalidYieldSharePercentageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidYieldSharePercentage(uint256 totalSharePercentage)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidYieldSharePercentageError(raw []byte) (*IDebtSubsidizerInvalidYieldSharePercentage, error) {
	out := new(IDebtSubsidizerInvalidYieldSharePercentage)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidYieldSharePercentage", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidYieldSlice represents a InvalidYieldSlice error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidYieldSlice struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidYieldSlice()
func IDebtSubsidizerInvalidYieldSliceErrorID() common.Hash {
	return common.HexToHash("0x1d18dd9358b8bac99c04d0a1149ec9fca50265639c372e50e9468ac996d78ee1")
}

// UnpackInvalidYieldSliceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidYieldSlice()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidYieldSliceError(raw []byte) (*IDebtSubsidizerInvalidYieldSlice, error) {
	out := new(IDebtSubsidizerInvalidYieldSlice)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidYieldSlice", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerLendingManagerAssetMismatch represents a LendingManagerAssetMismatch error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerLendingManagerAssetMismatch struct {
	VaultAsset common.Address
	LmAsset    common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerAssetMismatch(address vaultAsset, address lmAsset)
func IDebtSubsidizerLendingManagerAssetMismatchErrorID() common.Hash {
	return common.HexToHash("0x42a981d74553ce955a75472e22ff8184e6a5612f06aec91a457e357383d769f0")
}

// UnpackLendingManagerAssetMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerAssetMismatch(address vaultAsset, address lmAsset)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackLendingManagerAssetMismatchError(raw []byte) (*IDebtSubsidizerLendingManagerAssetMismatch, error) {
	out := new(IDebtSubsidizerLendingManagerAssetMismatch)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "LendingManagerAssetMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerLendingManagerNotSetForVault represents a LendingManagerNotSetForVault error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerLendingManagerNotSetForVault struct {
	VaultAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerNotSetForVault(address vaultAddress)
func IDebtSubsidizerLendingManagerNotSetForVaultErrorID() common.Hash {
	return common.HexToHash("0x1930e66911b345fdaa39a639549dc16c2105f0f4bc3ae1c1e39867527469205e")
}

// UnpackLendingManagerNotSetForVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerNotSetForVault(address vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackLendingManagerNotSetForVaultError(raw []byte) (*IDebtSubsidizerLendingManagerNotSetForVault, error) {
	out := new(IDebtSubsidizerLendingManagerNotSetForVault)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "LendingManagerNotSetForVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerMerkleRootNotSet represents a MerkleRootNotSet error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerMerkleRootNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MerkleRootNotSet()
func IDebtSubsidizerMerkleRootNotSetErrorID() common.Hash {
	return common.HexToHash("0x9f8a28f2b6ecfc78c926d70ccbdbc1af50319e4764903fa1ee9bb94dc558b392")
}

// UnpackMerkleRootNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MerkleRootNotSet()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackMerkleRootNotSetError(raw []byte) (*IDebtSubsidizerMerkleRootNotSet, error) {
	out := new(IDebtSubsidizerMerkleRootNotSet)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "MerkleRootNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerVaultAlreadyRegistered represents a VaultAlreadyRegistered error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultAlreadyRegistered struct {
	VaultAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VaultAlreadyRegistered(address vaultAddress)
func IDebtSubsidizerVaultAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x38bfcc166bcba4b5532c5572994cadc9ff436e7c322bb1d41b75caf05248d852")
}

// UnpackVaultAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VaultAlreadyRegistered(address vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultAlreadyRegisteredError(raw []byte) (*IDebtSubsidizerVaultAlreadyRegistered, error) {
	out := new(IDebtSubsidizerVaultAlreadyRegistered)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "VaultAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerVaultMismatch represents a VaultMismatch error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VaultMismatch()
func IDebtSubsidizerVaultMismatchErrorID() common.Hash {
	return common.HexToHash("0xc1faacc51457ff2b9dc86cf152d0b03efa014271f93627a43af5c22f0bb7c5ec")
}

// UnpackVaultMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VaultMismatch()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultMismatchError(raw []byte) (*IDebtSubsidizerVaultMismatch, error) {
	out := new(IDebtSubsidizerVaultMismatch)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "VaultMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerVaultNotRegistered represents a VaultNotRegistered error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultNotRegistered struct {
	VaultAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VaultNotRegistered(address vaultAddress)
func IDebtSubsidizerVaultNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x299f3425f9f512f17e69ec652b8e51ec669fa2236d04369d83965289ca6a5afb")
}

// UnpackVaultNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VaultNotRegistered(address vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultNotRegisteredError(raw []byte) (*IDebtSubsidizerVaultNotRegistered, error) {
	out := new(IDebtSubsidizerVaultNotRegistered)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "VaultNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}
