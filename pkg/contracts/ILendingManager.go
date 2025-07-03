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

// ILendingManagerMetaData contains all meta data concerning the ILendingManager contract.
var ILendingManagerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addSupportedMarket\",\"inputs\":[{\"name\":\"market\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositToLendingProtocol\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getGlobalCollateralFactor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLiquidationIncentive\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSupportedMarkets\",\"inputs\":[],\"outputs\":[{\"name\":\"markets\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalBorrowVolume\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalLiquidationVolume\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalMarketParticipants\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalSupplyVolume\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordLiquidationVolume\",\"inputs\":[{\"name\":\"liquidationAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"redeemAllCTokens\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amountRedeemed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeSupportedMarket\",\"inputs\":[{\"name\":\"market\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"repayBorrowBehalf\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"repayAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGlobalCollateralFactor\",\"inputs\":[{\"name\":\"_globalCollateralFactor\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLiquidationIncentive\",\"inputs\":[{\"name\":\"_liquidationIncentive\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalPrincipalDeposited\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateMarketParticipants\",\"inputs\":[{\"name\":\"_totalMarketParticipants\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFromLendingProtocol\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BorrowVolumeUpdated\",\"inputs\":[{\"name\":\"totalVolume\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"incrementAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DepositToProtocol\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GlobalCollateralFactorUpdated\",\"inputs\":[{\"name\":\"newFactor\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LendingManagerRoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LendingManagerRoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidationIncentiveUpdated\",\"inputs\":[{\"name\":\"newIncentive\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidationVolumeUpdated\",\"inputs\":[{\"name\":\"totalVolume\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"incrementAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PrincipalReset\",\"inputs\":[{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"trigger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SupplyVolumeUpdated\",\"inputs\":[{\"name\":\"totalVolume\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"incrementAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SupportedMarketAdded\",\"inputs\":[{\"name\":\"market\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SupportedMarketRemoved\",\"inputs\":[{\"name\":\"market\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawFromProtocol\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"YieldTransferred\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"YieldTransferredBatch\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"collections\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotRemoveLastAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientBalanceInProtocol\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LM_CallerNotRewardsController\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LM_CallerNotVault\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenMintFailed\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenMintFailedBytes\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenMintFailedReason\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRedeemFailed\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRedeemFailedBytes\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRedeemFailedReason\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRedeemUnderlyingFailed\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRedeemUnderlyingFailedBytes\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRedeemUnderlyingFailedReason\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRepayBorrowBehalfFailed\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRepayBorrowBehalfFailedBytes\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"LendingManagerCTokenRepayBorrowBehalfFailedReason\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"LendingManager__BalanceCheckFailed\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "ILendingManager",
}

// ILendingManager is an auto generated Go binding around an Ethereum contract.
type ILendingManager struct {
	abi abi.ABI
}

// NewILendingManager creates a new instance of ILendingManager.
func NewILendingManager() *ILendingManager {
	parsed, err := ILendingManagerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ILendingManager{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ILendingManager) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddSupportedMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x16a2d41f.
//
// Solidity: function addSupportedMarket(address market) returns()
func (iLendingManager *ILendingManager) PackAddSupportedMarket(market common.Address) []byte {
	enc, err := iLendingManager.abi.Pack("addSupportedMarket", market)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (iLendingManager *ILendingManager) PackAsset() []byte {
	enc, err := iLendingManager.abi.Pack("asset")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (iLendingManager *ILendingManager) UnpackAsset(data []byte) (common.Address, error) {
	out, err := iLendingManager.abi.Unpack("asset", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackCToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x69e527da.
//
// Solidity: function cToken() view returns(address)
func (iLendingManager *ILendingManager) PackCToken() []byte {
	enc, err := iLendingManager.abi.Pack("cToken")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x69e527da.
//
// Solidity: function cToken() view returns(address)
func (iLendingManager *ILendingManager) UnpackCToken(data []byte) (common.Address, error) {
	out, err := iLendingManager.abi.Unpack("cToken", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDepositToLendingProtocol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2188582.
//
// Solidity: function depositToLendingProtocol(uint256 amount) returns(bool success)
func (iLendingManager *ILendingManager) PackDepositToLendingProtocol(amount *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("depositToLendingProtocol", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositToLendingProtocol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd2188582.
//
// Solidity: function depositToLendingProtocol(uint256 amount) returns(bool success)
func (iLendingManager *ILendingManager) UnpackDepositToLendingProtocol(data []byte) (bool, error) {
	out, err := iLendingManager.abi.Unpack("depositToLendingProtocol", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackGetGlobalCollateralFactor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7216c451.
//
// Solidity: function getGlobalCollateralFactor() view returns(uint256)
func (iLendingManager *ILendingManager) PackGetGlobalCollateralFactor() []byte {
	enc, err := iLendingManager.abi.Pack("getGlobalCollateralFactor")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetGlobalCollateralFactor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7216c451.
//
// Solidity: function getGlobalCollateralFactor() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackGetGlobalCollateralFactor(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("getGlobalCollateralFactor", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetLiquidationIncentive is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc29a1d58.
//
// Solidity: function getLiquidationIncentive() view returns(uint256)
func (iLendingManager *ILendingManager) PackGetLiquidationIncentive() []byte {
	enc, err := iLendingManager.abi.Pack("getLiquidationIncentive")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetLiquidationIncentive is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc29a1d58.
//
// Solidity: function getLiquidationIncentive() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackGetLiquidationIncentive(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("getLiquidationIncentive", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetSupportedMarkets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf9476290.
//
// Solidity: function getSupportedMarkets() view returns(address[] markets)
func (iLendingManager *ILendingManager) PackGetSupportedMarkets() []byte {
	enc, err := iLendingManager.abi.Pack("getSupportedMarkets")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetSupportedMarkets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf9476290.
//
// Solidity: function getSupportedMarkets() view returns(address[] markets)
func (iLendingManager *ILendingManager) UnpackGetSupportedMarkets(data []byte) ([]common.Address, error) {
	out, err := iLendingManager.abi.Unpack("getSupportedMarkets", data)
	if err != nil {
		return *new([]common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, err
}

// PackGetTotalBorrowVolume is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d6d529c.
//
// Solidity: function getTotalBorrowVolume() view returns(uint256)
func (iLendingManager *ILendingManager) PackGetTotalBorrowVolume() []byte {
	enc, err := iLendingManager.abi.Pack("getTotalBorrowVolume")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalBorrowVolume is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1d6d529c.
//
// Solidity: function getTotalBorrowVolume() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackGetTotalBorrowVolume(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("getTotalBorrowVolume", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalLiquidationVolume is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1aeab45b.
//
// Solidity: function getTotalLiquidationVolume() view returns(uint256)
func (iLendingManager *ILendingManager) PackGetTotalLiquidationVolume() []byte {
	enc, err := iLendingManager.abi.Pack("getTotalLiquidationVolume")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalLiquidationVolume is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1aeab45b.
//
// Solidity: function getTotalLiquidationVolume() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackGetTotalLiquidationVolume(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("getTotalLiquidationVolume", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalMarketParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x948a092d.
//
// Solidity: function getTotalMarketParticipants() view returns(uint256)
func (iLendingManager *ILendingManager) PackGetTotalMarketParticipants() []byte {
	enc, err := iLendingManager.abi.Pack("getTotalMarketParticipants")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalMarketParticipants is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x948a092d.
//
// Solidity: function getTotalMarketParticipants() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackGetTotalMarketParticipants(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("getTotalMarketParticipants", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalSupplyVolume is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x919a2e5d.
//
// Solidity: function getTotalSupplyVolume() view returns(uint256)
func (iLendingManager *ILendingManager) PackGetTotalSupplyVolume() []byte {
	enc, err := iLendingManager.abi.Pack("getTotalSupplyVolume")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalSupplyVolume is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x919a2e5d.
//
// Solidity: function getTotalSupplyVolume() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackGetTotalSupplyVolume(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("getTotalSupplyVolume", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackRecordLiquidationVolume is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe372c8c4.
//
// Solidity: function recordLiquidationVolume(uint256 liquidationAmount) returns()
func (iLendingManager *ILendingManager) PackRecordLiquidationVolume(liquidationAmount *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("recordLiquidationVolume", liquidationAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRedeemAllCTokens is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5eb8fa6f.
//
// Solidity: function redeemAllCTokens(address recipient) returns(uint256 amountRedeemed)
func (iLendingManager *ILendingManager) PackRedeemAllCTokens(recipient common.Address) []byte {
	enc, err := iLendingManager.abi.Pack("redeemAllCTokens", recipient)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRedeemAllCTokens is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5eb8fa6f.
//
// Solidity: function redeemAllCTokens(address recipient) returns(uint256 amountRedeemed)
func (iLendingManager *ILendingManager) UnpackRedeemAllCTokens(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("redeemAllCTokens", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackRemoveSupportedMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0652320.
//
// Solidity: function removeSupportedMarket(address market) returns()
func (iLendingManager *ILendingManager) PackRemoveSupportedMarket(market common.Address) []byte {
	enc, err := iLendingManager.abi.Pack("removeSupportedMarket", market)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRepayBorrowBehalf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2608f818.
//
// Solidity: function repayBorrowBehalf(address borrower, uint256 repayAmount) returns(uint256)
func (iLendingManager *ILendingManager) PackRepayBorrowBehalf(borrower common.Address, repayAmount *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("repayBorrowBehalf", borrower, repayAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRepayBorrowBehalf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2608f818.
//
// Solidity: function repayBorrowBehalf(address borrower, uint256 repayAmount) returns(uint256)
func (iLendingManager *ILendingManager) UnpackRepayBorrowBehalf(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("repayBorrowBehalf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackSetGlobalCollateralFactor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfb603d2d.
//
// Solidity: function setGlobalCollateralFactor(uint256 _globalCollateralFactor) returns()
func (iLendingManager *ILendingManager) PackSetGlobalCollateralFactor(globalCollateralFactor *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("setGlobalCollateralFactor", globalCollateralFactor)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetLiquidationIncentive is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa8431081.
//
// Solidity: function setLiquidationIncentive(uint256 _liquidationIncentive) returns()
func (iLendingManager *ILendingManager) PackSetLiquidationIncentive(liquidationIncentive *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("setLiquidationIncentive", liquidationIncentive)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (iLendingManager *ILendingManager) PackTotalAssets() []byte {
	enc, err := iLendingManager.abi.Pack("totalAssets")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackTotalAssets(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("totalAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalPrincipalDeposited is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8070bb28.
//
// Solidity: function totalPrincipalDeposited() view returns(uint256)
func (iLendingManager *ILendingManager) PackTotalPrincipalDeposited() []byte {
	enc, err := iLendingManager.abi.Pack("totalPrincipalDeposited")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalPrincipalDeposited is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8070bb28.
//
// Solidity: function totalPrincipalDeposited() view returns(uint256)
func (iLendingManager *ILendingManager) UnpackTotalPrincipalDeposited(data []byte) (*big.Int, error) {
	out, err := iLendingManager.abi.Unpack("totalPrincipalDeposited", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackUpdateMarketParticipants is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cf42a78.
//
// Solidity: function updateMarketParticipants(uint256 _totalMarketParticipants) returns()
func (iLendingManager *ILendingManager) PackUpdateMarketParticipants(totalMarketParticipants *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("updateMarketParticipants", totalMarketParticipants)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackWithdrawFromLendingProtocol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0942d2e.
//
// Solidity: function withdrawFromLendingProtocol(uint256 amount) returns(bool success)
func (iLendingManager *ILendingManager) PackWithdrawFromLendingProtocol(amount *big.Int) []byte {
	enc, err := iLendingManager.abi.Pack("withdrawFromLendingProtocol", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdrawFromLendingProtocol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd0942d2e.
//
// Solidity: function withdrawFromLendingProtocol(uint256 amount) returns(bool success)
func (iLendingManager *ILendingManager) UnpackWithdrawFromLendingProtocol(data []byte) (bool, error) {
	out, err := iLendingManager.abi.Unpack("withdrawFromLendingProtocol", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// ILendingManagerBorrowVolumeUpdated represents a BorrowVolumeUpdated event raised by the ILendingManager contract.
type ILendingManagerBorrowVolumeUpdated struct {
	TotalVolume     *big.Int
	IncrementAmount *big.Int
	Timestamp       *big.Int
	Raw             *types.Log // Blockchain specific contextual infos
}

const ILendingManagerBorrowVolumeUpdatedEventName = "BorrowVolumeUpdated"

// ContractEventName returns the user-defined event name.
func (ILendingManagerBorrowVolumeUpdated) ContractEventName() string {
	return ILendingManagerBorrowVolumeUpdatedEventName
}

// UnpackBorrowVolumeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BorrowVolumeUpdated(uint256 indexed totalVolume, uint256 indexed incrementAmount, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackBorrowVolumeUpdatedEvent(log *types.Log) (*ILendingManagerBorrowVolumeUpdated, error) {
	event := "BorrowVolumeUpdated"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerBorrowVolumeUpdated)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerDepositToProtocol represents a DepositToProtocol event raised by the ILendingManager contract.
type ILendingManagerDepositToProtocol struct {
	Caller common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ILendingManagerDepositToProtocolEventName = "DepositToProtocol"

// ContractEventName returns the user-defined event name.
func (ILendingManagerDepositToProtocol) ContractEventName() string {
	return ILendingManagerDepositToProtocolEventName
}

// UnpackDepositToProtocolEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DepositToProtocol(address indexed caller, uint256 amount)
func (iLendingManager *ILendingManager) UnpackDepositToProtocolEvent(log *types.Log) (*ILendingManagerDepositToProtocol, error) {
	event := "DepositToProtocol"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerDepositToProtocol)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerGlobalCollateralFactorUpdated represents a GlobalCollateralFactorUpdated event raised by the ILendingManager contract.
type ILendingManagerGlobalCollateralFactorUpdated struct {
	NewFactor *big.Int
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILendingManagerGlobalCollateralFactorUpdatedEventName = "GlobalCollateralFactorUpdated"

// ContractEventName returns the user-defined event name.
func (ILendingManagerGlobalCollateralFactorUpdated) ContractEventName() string {
	return ILendingManagerGlobalCollateralFactorUpdatedEventName
}

// UnpackGlobalCollateralFactorUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GlobalCollateralFactorUpdated(uint256 indexed newFactor, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackGlobalCollateralFactorUpdatedEvent(log *types.Log) (*ILendingManagerGlobalCollateralFactorUpdated, error) {
	event := "GlobalCollateralFactorUpdated"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerGlobalCollateralFactorUpdated)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerLendingManagerRoleGranted represents a LendingManagerRoleGranted event raised by the ILendingManager contract.
type ILendingManagerLendingManagerRoleGranted struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILendingManagerLendingManagerRoleGrantedEventName = "LendingManagerRoleGranted"

// ContractEventName returns the user-defined event name.
func (ILendingManagerLendingManagerRoleGranted) ContractEventName() string {
	return ILendingManagerLendingManagerRoleGrantedEventName
}

// UnpackLendingManagerRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LendingManagerRoleGranted(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (iLendingManager *ILendingManager) UnpackLendingManagerRoleGrantedEvent(log *types.Log) (*ILendingManagerLendingManagerRoleGranted, error) {
	event := "LendingManagerRoleGranted"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerLendingManagerRoleGranted)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerLendingManagerRoleRevoked represents a LendingManagerRoleRevoked event raised by the ILendingManager contract.
type ILendingManagerLendingManagerRoleRevoked struct {
	Role      [32]byte
	Account   common.Address
	Sender    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILendingManagerLendingManagerRoleRevokedEventName = "LendingManagerRoleRevoked"

// ContractEventName returns the user-defined event name.
func (ILendingManagerLendingManagerRoleRevoked) ContractEventName() string {
	return ILendingManagerLendingManagerRoleRevokedEventName
}

// UnpackLendingManagerRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LendingManagerRoleRevoked(bytes32 indexed role, address indexed account, address sender, uint256 timestamp)
func (iLendingManager *ILendingManager) UnpackLendingManagerRoleRevokedEvent(log *types.Log) (*ILendingManagerLendingManagerRoleRevoked, error) {
	event := "LendingManagerRoleRevoked"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerLendingManagerRoleRevoked)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerLiquidationIncentiveUpdated represents a LiquidationIncentiveUpdated event raised by the ILendingManager contract.
type ILendingManagerLiquidationIncentiveUpdated struct {
	NewIncentive *big.Int
	Timestamp    *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const ILendingManagerLiquidationIncentiveUpdatedEventName = "LiquidationIncentiveUpdated"

// ContractEventName returns the user-defined event name.
func (ILendingManagerLiquidationIncentiveUpdated) ContractEventName() string {
	return ILendingManagerLiquidationIncentiveUpdatedEventName
}

// UnpackLiquidationIncentiveUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LiquidationIncentiveUpdated(uint256 indexed newIncentive, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackLiquidationIncentiveUpdatedEvent(log *types.Log) (*ILendingManagerLiquidationIncentiveUpdated, error) {
	event := "LiquidationIncentiveUpdated"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerLiquidationIncentiveUpdated)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerLiquidationVolumeUpdated represents a LiquidationVolumeUpdated event raised by the ILendingManager contract.
type ILendingManagerLiquidationVolumeUpdated struct {
	TotalVolume     *big.Int
	IncrementAmount *big.Int
	Timestamp       *big.Int
	Raw             *types.Log // Blockchain specific contextual infos
}

const ILendingManagerLiquidationVolumeUpdatedEventName = "LiquidationVolumeUpdated"

// ContractEventName returns the user-defined event name.
func (ILendingManagerLiquidationVolumeUpdated) ContractEventName() string {
	return ILendingManagerLiquidationVolumeUpdatedEventName
}

// UnpackLiquidationVolumeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LiquidationVolumeUpdated(uint256 indexed totalVolume, uint256 indexed incrementAmount, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackLiquidationVolumeUpdatedEvent(log *types.Log) (*ILendingManagerLiquidationVolumeUpdated, error) {
	event := "LiquidationVolumeUpdated"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerLiquidationVolumeUpdated)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerPrincipalReset represents a PrincipalReset event raised by the ILendingManager contract.
type ILendingManagerPrincipalReset struct {
	OldValue *big.Int
	Trigger  common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const ILendingManagerPrincipalResetEventName = "PrincipalReset"

// ContractEventName returns the user-defined event name.
func (ILendingManagerPrincipalReset) ContractEventName() string {
	return ILendingManagerPrincipalResetEventName
}

// UnpackPrincipalResetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PrincipalReset(uint256 oldValue, address indexed trigger)
func (iLendingManager *ILendingManager) UnpackPrincipalResetEvent(log *types.Log) (*ILendingManagerPrincipalReset, error) {
	event := "PrincipalReset"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerPrincipalReset)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerSupplyVolumeUpdated represents a SupplyVolumeUpdated event raised by the ILendingManager contract.
type ILendingManagerSupplyVolumeUpdated struct {
	TotalVolume     *big.Int
	IncrementAmount *big.Int
	Timestamp       *big.Int
	Raw             *types.Log // Blockchain specific contextual infos
}

const ILendingManagerSupplyVolumeUpdatedEventName = "SupplyVolumeUpdated"

// ContractEventName returns the user-defined event name.
func (ILendingManagerSupplyVolumeUpdated) ContractEventName() string {
	return ILendingManagerSupplyVolumeUpdatedEventName
}

// UnpackSupplyVolumeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SupplyVolumeUpdated(uint256 indexed totalVolume, uint256 indexed incrementAmount, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackSupplyVolumeUpdatedEvent(log *types.Log) (*ILendingManagerSupplyVolumeUpdated, error) {
	event := "SupplyVolumeUpdated"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerSupplyVolumeUpdated)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerSupportedMarketAdded represents a SupportedMarketAdded event raised by the ILendingManager contract.
type ILendingManagerSupportedMarketAdded struct {
	Market    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILendingManagerSupportedMarketAddedEventName = "SupportedMarketAdded"

// ContractEventName returns the user-defined event name.
func (ILendingManagerSupportedMarketAdded) ContractEventName() string {
	return ILendingManagerSupportedMarketAddedEventName
}

// UnpackSupportedMarketAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SupportedMarketAdded(address indexed market, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackSupportedMarketAddedEvent(log *types.Log) (*ILendingManagerSupportedMarketAdded, error) {
	event := "SupportedMarketAdded"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerSupportedMarketAdded)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerSupportedMarketRemoved represents a SupportedMarketRemoved event raised by the ILendingManager contract.
type ILendingManagerSupportedMarketRemoved struct {
	Market    common.Address
	Timestamp *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILendingManagerSupportedMarketRemovedEventName = "SupportedMarketRemoved"

// ContractEventName returns the user-defined event name.
func (ILendingManagerSupportedMarketRemoved) ContractEventName() string {
	return ILendingManagerSupportedMarketRemovedEventName
}

// UnpackSupportedMarketRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SupportedMarketRemoved(address indexed market, uint256 indexed timestamp)
func (iLendingManager *ILendingManager) UnpackSupportedMarketRemovedEvent(log *types.Log) (*ILendingManagerSupportedMarketRemoved, error) {
	event := "SupportedMarketRemoved"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerSupportedMarketRemoved)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerWithdrawFromProtocol represents a WithdrawFromProtocol event raised by the ILendingManager contract.
type ILendingManagerWithdrawFromProtocol struct {
	Caller common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ILendingManagerWithdrawFromProtocolEventName = "WithdrawFromProtocol"

// ContractEventName returns the user-defined event name.
func (ILendingManagerWithdrawFromProtocol) ContractEventName() string {
	return ILendingManagerWithdrawFromProtocolEventName
}

// UnpackWithdrawFromProtocolEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WithdrawFromProtocol(address indexed caller, uint256 amount)
func (iLendingManager *ILendingManager) UnpackWithdrawFromProtocolEvent(log *types.Log) (*ILendingManagerWithdrawFromProtocol, error) {
	event := "WithdrawFromProtocol"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerWithdrawFromProtocol)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerYieldTransferred represents a YieldTransferred event raised by the ILendingManager contract.
type ILendingManagerYieldTransferred struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ILendingManagerYieldTransferredEventName = "YieldTransferred"

// ContractEventName returns the user-defined event name.
func (ILendingManagerYieldTransferred) ContractEventName() string {
	return ILendingManagerYieldTransferredEventName
}

// UnpackYieldTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event YieldTransferred(address indexed recipient, uint256 amount)
func (iLendingManager *ILendingManager) UnpackYieldTransferredEvent(log *types.Log) (*ILendingManagerYieldTransferred, error) {
	event := "YieldTransferred"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerYieldTransferred)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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

// ILendingManagerYieldTransferredBatch represents a YieldTransferredBatch event raised by the ILendingManager contract.
type ILendingManagerYieldTransferredBatch struct {
	Recipient   common.Address
	TotalAmount *big.Int
	Collections []common.Address
	Amounts     []*big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const ILendingManagerYieldTransferredBatchEventName = "YieldTransferredBatch"

// ContractEventName returns the user-defined event name.
func (ILendingManagerYieldTransferredBatch) ContractEventName() string {
	return ILendingManagerYieldTransferredBatchEventName
}

// UnpackYieldTransferredBatchEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event YieldTransferredBatch(address indexed recipient, uint256 totalAmount, address[] collections, uint256[] amounts)
func (iLendingManager *ILendingManager) UnpackYieldTransferredBatchEvent(log *types.Log) (*ILendingManagerYieldTransferredBatch, error) {
	event := "YieldTransferredBatch"
	if log.Topics[0] != iLendingManager.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ILendingManagerYieldTransferredBatch)
	if len(log.Data) > 0 {
		if err := iLendingManager.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iLendingManager.abi.Events[event].Inputs {
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
func (iLendingManager *ILendingManager) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["AddressZero"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackAddressZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["CannotRemoveLastAdmin"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackCannotRemoveLastAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["InsufficientBalanceInProtocol"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackInsufficientBalanceInProtocolError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LMCallerNotRewardsController"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLMCallerNotRewardsControllerError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LMCallerNotVault"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLMCallerNotVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenMintFailed"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenMintFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenMintFailedBytes"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenMintFailedBytesError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenMintFailedReason"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenMintFailedReasonError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRedeemFailed"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRedeemFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRedeemFailedBytes"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRedeemFailedBytesError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRedeemFailedReason"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRedeemFailedReasonError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRedeemUnderlyingFailed"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRedeemUnderlyingFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRedeemUnderlyingFailedBytes"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRedeemUnderlyingFailedBytesError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRedeemUnderlyingFailedReason"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRedeemUnderlyingFailedReasonError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRepayBorrowBehalfFailed"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRepayBorrowBehalfFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRepayBorrowBehalfFailedBytes"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRepayBorrowBehalfFailedBytesError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerCTokenRepayBorrowBehalfFailedReason"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerCTokenRepayBorrowBehalfFailedReasonError(raw[4:])
	}
	if bytes.Equal(raw[:4], iLendingManager.abi.Errors["LendingManagerBalanceCheckFailed"].ID.Bytes()[:4]) {
		return iLendingManager.UnpackLendingManagerBalanceCheckFailedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ILendingManagerAddressZero represents a AddressZero error raised by the ILendingManager contract.
type ILendingManagerAddressZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressZero()
func ILendingManagerAddressZeroErrorID() common.Hash {
	return common.HexToHash("0x9fabe1c19979afc45ec7efec1bde2c38021c590a0ce42965cf55b3f518197f02")
}

// UnpackAddressZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressZero()
func (iLendingManager *ILendingManager) UnpackAddressZeroError(raw []byte) (*ILendingManagerAddressZero, error) {
	out := new(ILendingManagerAddressZero)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "AddressZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerCannotRemoveLastAdmin represents a CannotRemoveLastAdmin error raised by the ILendingManager contract.
type ILendingManagerCannotRemoveLastAdmin struct {
	Role [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotRemoveLastAdmin(bytes32 role)
func ILendingManagerCannotRemoveLastAdminErrorID() common.Hash {
	return common.HexToHash("0xed12cc3613424688424c7a5fa4c692ad52518e0161af68235ae9d449ea5487e4")
}

// UnpackCannotRemoveLastAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotRemoveLastAdmin(bytes32 role)
func (iLendingManager *ILendingManager) UnpackCannotRemoveLastAdminError(raw []byte) (*ILendingManagerCannotRemoveLastAdmin, error) {
	out := new(ILendingManagerCannotRemoveLastAdmin)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "CannotRemoveLastAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerInsufficientBalanceInProtocol represents a InsufficientBalanceInProtocol error raised by the ILendingManager contract.
type ILendingManagerInsufficientBalanceInProtocol struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientBalanceInProtocol()
func ILendingManagerInsufficientBalanceInProtocolErrorID() common.Hash {
	return common.HexToHash("0x7e65ad2ca86fa6ae58f772c31968382bed67947d6803d528a3c28871cbee7948")
}

// UnpackInsufficientBalanceInProtocolError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientBalanceInProtocol()
func (iLendingManager *ILendingManager) UnpackInsufficientBalanceInProtocolError(raw []byte) (*ILendingManagerInsufficientBalanceInProtocol, error) {
	out := new(ILendingManagerInsufficientBalanceInProtocol)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "InsufficientBalanceInProtocol", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLMCallerNotRewardsController represents a LM_CallerNotRewardsController error raised by the ILendingManager contract.
type ILendingManagerLMCallerNotRewardsController struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LM_CallerNotRewardsController(address caller)
func ILendingManagerLMCallerNotRewardsControllerErrorID() common.Hash {
	return common.HexToHash("0x933221863ef0abe30ee597955dd2f95920f6a7aefbbbb243f29e245667017f88")
}

// UnpackLMCallerNotRewardsControllerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LM_CallerNotRewardsController(address caller)
func (iLendingManager *ILendingManager) UnpackLMCallerNotRewardsControllerError(raw []byte) (*ILendingManagerLMCallerNotRewardsController, error) {
	out := new(ILendingManagerLMCallerNotRewardsController)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LMCallerNotRewardsController", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLMCallerNotVault represents a LM_CallerNotVault error raised by the ILendingManager contract.
type ILendingManagerLMCallerNotVault struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LM_CallerNotVault(address caller)
func ILendingManagerLMCallerNotVaultErrorID() common.Hash {
	return common.HexToHash("0x4383e864424f8a9091b3c7994a810b69c877cba41f26eb3aa1c3c938b57ec5ea")
}

// UnpackLMCallerNotVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LM_CallerNotVault(address caller)
func (iLendingManager *ILendingManager) UnpackLMCallerNotVaultError(raw []byte) (*ILendingManagerLMCallerNotVault, error) {
	out := new(ILendingManagerLMCallerNotVault)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LMCallerNotVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenMintFailed represents a LendingManagerCTokenMintFailed error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenMintFailed struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenMintFailed(uint256 errorCode)
func ILendingManagerLendingManagerCTokenMintFailedErrorID() common.Hash {
	return common.HexToHash("0xa1c750daf7d7446389741abb667b1e85bbdea86d4e5fbc0e2d4ad2cbc0acda13")
}

// UnpackLendingManagerCTokenMintFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenMintFailed(uint256 errorCode)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenMintFailedError(raw []byte) (*ILendingManagerLendingManagerCTokenMintFailed, error) {
	out := new(ILendingManagerLendingManagerCTokenMintFailed)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenMintFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenMintFailedBytes represents a LendingManagerCTokenMintFailedBytes error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenMintFailedBytes struct {
	Data []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenMintFailedBytes(bytes data)
func ILendingManagerLendingManagerCTokenMintFailedBytesErrorID() common.Hash {
	return common.HexToHash("0x1da159508578ed9bcab37ab7dcaf273e7c130e17ae6c43c4a490803eff4b2896")
}

// UnpackLendingManagerCTokenMintFailedBytesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenMintFailedBytes(bytes data)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenMintFailedBytesError(raw []byte) (*ILendingManagerLendingManagerCTokenMintFailedBytes, error) {
	out := new(ILendingManagerLendingManagerCTokenMintFailedBytes)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenMintFailedBytes", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenMintFailedReason represents a LendingManagerCTokenMintFailedReason error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenMintFailedReason struct {
	Reason string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenMintFailedReason(string reason)
func ILendingManagerLendingManagerCTokenMintFailedReasonErrorID() common.Hash {
	return common.HexToHash("0x05709e52763623f8375c7e270d5625504bdca2c76ff7ac0d068887cdf765beb2")
}

// UnpackLendingManagerCTokenMintFailedReasonError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenMintFailedReason(string reason)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenMintFailedReasonError(raw []byte) (*ILendingManagerLendingManagerCTokenMintFailedReason, error) {
	out := new(ILendingManagerLendingManagerCTokenMintFailedReason)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenMintFailedReason", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRedeemFailed represents a LendingManagerCTokenRedeemFailed error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRedeemFailed struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRedeemFailed(uint256 errorCode)
func ILendingManagerLendingManagerCTokenRedeemFailedErrorID() common.Hash {
	return common.HexToHash("0x515ce46375f243d8d3041b9439c428f9020cce5953daa9f3b478fd80ac38d848")
}

// UnpackLendingManagerCTokenRedeemFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRedeemFailed(uint256 errorCode)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRedeemFailedError(raw []byte) (*ILendingManagerLendingManagerCTokenRedeemFailed, error) {
	out := new(ILendingManagerLendingManagerCTokenRedeemFailed)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRedeemFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRedeemFailedBytes represents a LendingManagerCTokenRedeemFailedBytes error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRedeemFailedBytes struct {
	Data []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRedeemFailedBytes(bytes data)
func ILendingManagerLendingManagerCTokenRedeemFailedBytesErrorID() common.Hash {
	return common.HexToHash("0x450178724d31f7431d4be46e8aa13faed9da2fd371e6978016b4a257017fcb04")
}

// UnpackLendingManagerCTokenRedeemFailedBytesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRedeemFailedBytes(bytes data)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRedeemFailedBytesError(raw []byte) (*ILendingManagerLendingManagerCTokenRedeemFailedBytes, error) {
	out := new(ILendingManagerLendingManagerCTokenRedeemFailedBytes)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRedeemFailedBytes", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRedeemFailedReason represents a LendingManagerCTokenRedeemFailedReason error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRedeemFailedReason struct {
	Reason string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRedeemFailedReason(string reason)
func ILendingManagerLendingManagerCTokenRedeemFailedReasonErrorID() common.Hash {
	return common.HexToHash("0xdc358670a17b9c7a3139020a4248459f8f1815353e1b5db2453066755f142137")
}

// UnpackLendingManagerCTokenRedeemFailedReasonError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRedeemFailedReason(string reason)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRedeemFailedReasonError(raw []byte) (*ILendingManagerLendingManagerCTokenRedeemFailedReason, error) {
	out := new(ILendingManagerLendingManagerCTokenRedeemFailedReason)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRedeemFailedReason", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRedeemUnderlyingFailed represents a LendingManagerCTokenRedeemUnderlyingFailed error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRedeemUnderlyingFailed struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRedeemUnderlyingFailed(uint256 errorCode)
func ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedErrorID() common.Hash {
	return common.HexToHash("0xbb73192e5be6149d31755d6ff8a10f8044576e5ec1bd0c5a089acb69efdc2109")
}

// UnpackLendingManagerCTokenRedeemUnderlyingFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRedeemUnderlyingFailed(uint256 errorCode)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRedeemUnderlyingFailedError(raw []byte) (*ILendingManagerLendingManagerCTokenRedeemUnderlyingFailed, error) {
	out := new(ILendingManagerLendingManagerCTokenRedeemUnderlyingFailed)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRedeemUnderlyingFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedBytes represents a LendingManagerCTokenRedeemUnderlyingFailedBytes error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedBytes struct {
	Data []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRedeemUnderlyingFailedBytes(bytes data)
func ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedBytesErrorID() common.Hash {
	return common.HexToHash("0x2d5607502b6ce219699a3ab0d525a18175364d0eb472eb9c6ce04f94d287c3c8")
}

// UnpackLendingManagerCTokenRedeemUnderlyingFailedBytesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRedeemUnderlyingFailedBytes(bytes data)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRedeemUnderlyingFailedBytesError(raw []byte) (*ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedBytes, error) {
	out := new(ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedBytes)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRedeemUnderlyingFailedBytes", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedReason represents a LendingManagerCTokenRedeemUnderlyingFailedReason error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedReason struct {
	Reason string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRedeemUnderlyingFailedReason(string reason)
func ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedReasonErrorID() common.Hash {
	return common.HexToHash("0xf299ba20be26876d1bef3a9c70f20b9ad9c5f6617b6393dad1dd80eaccf664fc")
}

// UnpackLendingManagerCTokenRedeemUnderlyingFailedReasonError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRedeemUnderlyingFailedReason(string reason)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRedeemUnderlyingFailedReasonError(raw []byte) (*ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedReason, error) {
	out := new(ILendingManagerLendingManagerCTokenRedeemUnderlyingFailedReason)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRedeemUnderlyingFailedReason", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailed represents a LendingManagerCTokenRepayBorrowBehalfFailed error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailed struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRepayBorrowBehalfFailed(uint256 errorCode)
func ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedErrorID() common.Hash {
	return common.HexToHash("0x442a3ba33016ff23619388e5f6d7afc5446e27b6eec44dbb3651428c32644832")
}

// UnpackLendingManagerCTokenRepayBorrowBehalfFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRepayBorrowBehalfFailed(uint256 errorCode)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRepayBorrowBehalfFailedError(raw []byte) (*ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailed, error) {
	out := new(ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailed)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRepayBorrowBehalfFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedBytes represents a LendingManagerCTokenRepayBorrowBehalfFailedBytes error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedBytes struct {
	Data []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRepayBorrowBehalfFailedBytes(bytes data)
func ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedBytesErrorID() common.Hash {
	return common.HexToHash("0x88aa8830fab014f1fa077f0b4f1cd15545b1d8164a2906c1e4d009a1d73da5eb")
}

// UnpackLendingManagerCTokenRepayBorrowBehalfFailedBytesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRepayBorrowBehalfFailedBytes(bytes data)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRepayBorrowBehalfFailedBytesError(raw []byte) (*ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedBytes, error) {
	out := new(ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedBytes)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRepayBorrowBehalfFailedBytes", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedReason represents a LendingManagerCTokenRepayBorrowBehalfFailedReason error raised by the ILendingManager contract.
type ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedReason struct {
	Reason string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerCTokenRepayBorrowBehalfFailedReason(string reason)
func ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedReasonErrorID() common.Hash {
	return common.HexToHash("0x4543929a1c3b442001f4d52746255af38c1e2ffbd24d8277baff8960b637a9f7")
}

// UnpackLendingManagerCTokenRepayBorrowBehalfFailedReasonError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerCTokenRepayBorrowBehalfFailedReason(string reason)
func (iLendingManager *ILendingManager) UnpackLendingManagerCTokenRepayBorrowBehalfFailedReasonError(raw []byte) (*ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedReason, error) {
	out := new(ILendingManagerLendingManagerCTokenRepayBorrowBehalfFailedReason)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerCTokenRepayBorrowBehalfFailedReason", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ILendingManagerLendingManagerBalanceCheckFailed represents a LendingManager__BalanceCheckFailed error raised by the ILendingManager contract.
type ILendingManagerLendingManagerBalanceCheckFailed struct {
	Reason   string
	Expected *big.Int
	Actual   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManager__BalanceCheckFailed(string reason, uint256 expected, uint256 actual)
func ILendingManagerLendingManagerBalanceCheckFailedErrorID() common.Hash {
	return common.HexToHash("0x3e62af153c5bf3b4b5f5f7717894822a1e4c857b7e26a3e1312c614ae52d6d7a")
}

// UnpackLendingManagerBalanceCheckFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManager__BalanceCheckFailed(string reason, uint256 expected, uint256 actual)
func (iLendingManager *ILendingManager) UnpackLendingManagerBalanceCheckFailedError(raw []byte) (*ILendingManagerLendingManagerBalanceCheckFailed, error) {
	out := new(ILendingManagerLendingManagerBalanceCheckFailed)
	if err := iLendingManager.abi.UnpackIntoInterface(out, "LendingManagerBalanceCheckFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
