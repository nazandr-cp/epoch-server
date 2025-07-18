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

// ICollectionsVaultMetaData contains all meta data concerning the ICollectionsVault contract.
var ICollectionsVaultMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEBT_SUBSIDIZER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allocateCumulativeYieldToEpoch\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocateEpochYield\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocateYieldToEpoch\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"applyCollectionYieldForEpoch\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[],\"outputs\":[{\"name\":\"assetTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collectionTotalAssetsDeposited\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToAssets\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToShares\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositForCollection\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEpochManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCollectionPerformanceScore\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCollectionTotalBorrowVolume\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCollectionTotalYieldGenerated\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentEpochYield\",\"inputs\":[{\"name\":\"includeNonShared\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"availableYield\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochYieldAllocated\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRemainingCumulativeYield\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalAvailableYield\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"indexCollectionsDeposits\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isCollectionOperator\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lendingManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractILendingManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"maxAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"maxShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"maxShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"maxAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintForCollection\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewDeposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewMint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRedeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWithdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordCollectionBorrowVolume\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"borrowAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"redeemForCollection\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"repayBorrowBehalf\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"borrower\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"repayBorrowBehalfBatch\",\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"borrowers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resetEpochCollectionYieldFlags\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"collections\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCollectionRegistry\",\"inputs\":[{\"name\":\"_collectionRegistryAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDebtSubsidizer\",\"inputs\":[{\"name\":\"_debtSubsidizerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEpochManager\",\"inputs\":[{\"name\":\"_epochManagerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLendingManager\",\"inputs\":[{\"name\":\"_lendingManagerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"totalManagedAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalAssetsDeposited\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalCollectionYieldShareBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalYieldAllocated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalYieldReserved\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferForCollection\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"underlying\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateCollectionPerformanceScore\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"score\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validateCumulativeClaims\",\"inputs\":[{\"name\":\"totalClaimedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawForCollection\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionAccessGranted\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionAccessRevoked\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionBorrowVolumeUpdated\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalVolume\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"incrementAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionDeposit\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cTokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionPerformanceUpdated\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"performanceScore\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionRegistryUpdated\",\"inputs\":[{\"name\":\"oldRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionTransfer\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionWithdraw\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cTokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionYieldAccrued\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"yieldAccrued\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newTotalDeposits\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"globalIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"previousCollectionIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionYieldAppliedForEpoch\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"collection\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"yieldSharePercentage\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"yieldAdded\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newTotalDeposits\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionYieldGenerated\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"yieldAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectionYieldIndexed\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cTokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerCallFailed\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochManagerCallUnavailable\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"functionName\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LendingManagerCallFailed\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operation\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LendingManagerChanged\",\"inputs\":[{\"name\":\"oldLendingManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newLendingManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"changedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultYieldAllocatedToEpoch\",\"inputs\":[{\"name\":\"epochId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"YieldBatchRepaid\",\"inputs\":[{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllocExceedsAvail\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllocationAmountZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AllocationUnderflow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BatchSizeExceedsLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CollectionInsufficientBalance\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"CollectionNotRegistered\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EpochManagerAllocationFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EpochManagerNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EpochManagerUnavailable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExcessiveYieldAmount\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAllowed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"FunctionDisabledUse\",\"inputs\":[{\"name\":\"functionName\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InsufficientBalanceInProtocol\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidEpochId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LendingManagerDepositFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LendingManagerMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LendingManagerWithdrawFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoActiveEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoCumulativeYield\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PerfScoreExceedsMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RedeemRoundsToZero\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RepayBorrowFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReqAmountExceeds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ShareBalanceUnderflow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TotalAllocExceeds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedCollectionAccess\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Vault_InsufficientBalancePostLMWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"YieldAlreadyApplied\",\"inputs\":[]}]",
	ID:  "ICollectionsVault",
}

// ICollectionsVault is an auto generated Go binding around an Ethereum contract.
type ICollectionsVault struct {
	abi abi.ABI
}

// NewICollectionsVault creates a new instance of ICollectionsVault.
func NewICollectionsVault() *ICollectionsVault {
	parsed, err := ICollectionsVaultMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ICollectionsVault{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ICollectionsVault) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (iCollectionsVault *ICollectionsVault) PackADMINROLE() []byte {
	enc, err := iCollectionsVault.abi.Pack("ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (iCollectionsVault *ICollectionsVault) UnpackADMINROLE(data []byte) ([32]byte, error) {
	out, err := iCollectionsVault.abi.Unpack("ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackDEBTSUBSIDIZERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x96a86a0c.
//
// Solidity: function DEBT_SUBSIDIZER_ROLE() view returns(bytes32)
func (iCollectionsVault *ICollectionsVault) PackDEBTSUBSIDIZERROLE() []byte {
	enc, err := iCollectionsVault.abi.Pack("DEBT_SUBSIDIZER_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDEBTSUBSIDIZERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x96a86a0c.
//
// Solidity: function DEBT_SUBSIDIZER_ROLE() view returns(bytes32)
func (iCollectionsVault *ICollectionsVault) UnpackDEBTSUBSIDIZERROLE(data []byte) ([32]byte, error) {
	out, err := iCollectionsVault.abi.Unpack("DEBT_SUBSIDIZER_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackAllocateCumulativeYieldToEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1b7f212.
//
// Solidity: function allocateCumulativeYieldToEpoch(uint256 epochId, uint256 amount) returns()
func (iCollectionsVault *ICollectionsVault) PackAllocateCumulativeYieldToEpoch(epochId *big.Int, amount *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("allocateCumulativeYieldToEpoch", epochId, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAllocateEpochYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf30bd63f.
//
// Solidity: function allocateEpochYield(uint256 amount) returns()
func (iCollectionsVault *ICollectionsVault) PackAllocateEpochYield(amount *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("allocateEpochYield", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAllocateYieldToEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x476f4840.
//
// Solidity: function allocateYieldToEpoch(uint256 epochId) returns()
func (iCollectionsVault *ICollectionsVault) PackAllocateYieldToEpoch(epochId *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("allocateYieldToEpoch", epochId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackApplyCollectionYieldForEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5780b963.
//
// Solidity: function applyCollectionYieldForEpoch(address collection, uint256 epochId) returns()
func (iCollectionsVault *ICollectionsVault) PackApplyCollectionYieldForEpoch(collection common.Address, epochId *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("applyCollectionYieldForEpoch", collection, epochId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (iCollectionsVault *ICollectionsVault) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (iCollectionsVault *ICollectionsVault) UnpackApprove(data []byte) (bool, error) {
	out, err := iCollectionsVault.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address assetTokenAddress)
func (iCollectionsVault *ICollectionsVault) PackAsset() []byte {
	enc, err := iCollectionsVault.abi.Pack("asset")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address assetTokenAddress)
func (iCollectionsVault *ICollectionsVault) UnpackAsset(data []byte) (common.Address, error) {
	out, err := iCollectionsVault.abi.Unpack("asset", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackBalanceOf(account common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackCollectionTotalAssetsDeposited is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa40589f8.
//
// Solidity: function collectionTotalAssetsDeposited(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackCollectionTotalAssetsDeposited(collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("collectionTotalAssetsDeposited", collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCollectionTotalAssetsDeposited is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa40589f8.
//
// Solidity: function collectionTotalAssetsDeposited(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionTotalAssetsDeposited(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("collectionTotalAssetsDeposited", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackConvertToAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackConvertToAssets(shares *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("convertToAssets", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackConvertToAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackConvertToAssets(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("convertToAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackConvertToShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackConvertToShares(assets *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("convertToShares", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackConvertToShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackConvertToShares(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("convertToShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (iCollectionsVault *ICollectionsVault) PackDecimals() []byte {
	enc, err := iCollectionsVault.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (iCollectionsVault *ICollectionsVault) UnpackDecimals(data []byte) (uint8, error) {
	out, err := iCollectionsVault.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackDeposit(assets *big.Int, receiver common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("deposit", assets, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackDeposit(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("deposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackDepositForCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x44061d90.
//
// Solidity: function depositForCollection(uint256 assets, address receiver, address collectionAddress) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackDepositForCollection(assets *big.Int, receiver common.Address, collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("depositForCollection", assets, receiver, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositForCollection is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x44061d90.
//
// Solidity: function depositForCollection(uint256 assets, address receiver, address collectionAddress) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackDepositForCollection(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("depositForCollection", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackEpochManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2d2bfe3.
//
// Solidity: function epochManager() view returns(address)
func (iCollectionsVault *ICollectionsVault) PackEpochManager() []byte {
	enc, err := iCollectionsVault.abi.Pack("epochManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEpochManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe2d2bfe3.
//
// Solidity: function epochManager() view returns(address)
func (iCollectionsVault *ICollectionsVault) UnpackEpochManager(data []byte) (common.Address, error) {
	out, err := iCollectionsVault.abi.Unpack("epochManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetCollectionPerformanceScore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4507ff6b.
//
// Solidity: function getCollectionPerformanceScore(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackGetCollectionPerformanceScore(collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("getCollectionPerformanceScore", collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCollectionPerformanceScore is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4507ff6b.
//
// Solidity: function getCollectionPerformanceScore(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackGetCollectionPerformanceScore(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getCollectionPerformanceScore", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetCollectionTotalBorrowVolume is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8ec4ceff.
//
// Solidity: function getCollectionTotalBorrowVolume(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackGetCollectionTotalBorrowVolume(collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("getCollectionTotalBorrowVolume", collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCollectionTotalBorrowVolume is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8ec4ceff.
//
// Solidity: function getCollectionTotalBorrowVolume(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackGetCollectionTotalBorrowVolume(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getCollectionTotalBorrowVolume", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetCollectionTotalYieldGenerated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3c73ecfd.
//
// Solidity: function getCollectionTotalYieldGenerated(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackGetCollectionTotalYieldGenerated(collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("getCollectionTotalYieldGenerated", collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCollectionTotalYieldGenerated is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3c73ecfd.
//
// Solidity: function getCollectionTotalYieldGenerated(address collectionAddress) view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackGetCollectionTotalYieldGenerated(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getCollectionTotalYieldGenerated", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetCurrentEpochYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6aebd973.
//
// Solidity: function getCurrentEpochYield(bool includeNonShared) view returns(uint256 availableYield)
func (iCollectionsVault *ICollectionsVault) PackGetCurrentEpochYield(includeNonShared bool) []byte {
	enc, err := iCollectionsVault.abi.Pack("getCurrentEpochYield", includeNonShared)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentEpochYield is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6aebd973.
//
// Solidity: function getCurrentEpochYield(bool includeNonShared) view returns(uint256 availableYield)
func (iCollectionsVault *ICollectionsVault) UnpackGetCurrentEpochYield(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getCurrentEpochYield", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetEpochYieldAllocated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb777451c.
//
// Solidity: function getEpochYieldAllocated(uint256 epochId) view returns(uint256 amount)
func (iCollectionsVault *ICollectionsVault) PackGetEpochYieldAllocated(epochId *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("getEpochYieldAllocated", epochId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEpochYieldAllocated is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb777451c.
//
// Solidity: function getEpochYieldAllocated(uint256 epochId) view returns(uint256 amount)
func (iCollectionsVault *ICollectionsVault) UnpackGetEpochYieldAllocated(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getEpochYieldAllocated", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetRemainingCumulativeYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8722fe09.
//
// Solidity: function getRemainingCumulativeYield() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackGetRemainingCumulativeYield() []byte {
	enc, err := iCollectionsVault.abi.Pack("getRemainingCumulativeYield")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetRemainingCumulativeYield is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8722fe09.
//
// Solidity: function getRemainingCumulativeYield() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackGetRemainingCumulativeYield(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getRemainingCumulativeYield", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalAvailableYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf73e9f5e.
//
// Solidity: function getTotalAvailableYield() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackGetTotalAvailableYield() []byte {
	enc, err := iCollectionsVault.abi.Pack("getTotalAvailableYield")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalAvailableYield is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf73e9f5e.
//
// Solidity: function getTotalAvailableYield() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackGetTotalAvailableYield(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("getTotalAvailableYield", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackIndexCollectionsDeposits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe5890f20.
//
// Solidity: function indexCollectionsDeposits() returns()
func (iCollectionsVault *ICollectionsVault) PackIndexCollectionsDeposits() []byte {
	enc, err := iCollectionsVault.abi.Pack("indexCollectionsDeposits")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackIsCollectionOperator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27205380.
//
// Solidity: function isCollectionOperator(address collectionAddress, address operator) view returns(bool)
func (iCollectionsVault *ICollectionsVault) PackIsCollectionOperator(collectionAddress common.Address, operator common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("isCollectionOperator", collectionAddress, operator)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsCollectionOperator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x27205380.
//
// Solidity: function isCollectionOperator(address collectionAddress, address operator) view returns(bool)
func (iCollectionsVault *ICollectionsVault) UnpackIsCollectionOperator(data []byte) (bool, error) {
	out, err := iCollectionsVault.abi.Unpack("isCollectionOperator", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackLendingManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9e3b77af.
//
// Solidity: function lendingManager() view returns(address)
func (iCollectionsVault *ICollectionsVault) PackLendingManager() []byte {
	enc, err := iCollectionsVault.abi.Pack("lendingManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLendingManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9e3b77af.
//
// Solidity: function lendingManager() view returns(address)
func (iCollectionsVault *ICollectionsVault) UnpackLendingManager(data []byte) (common.Address, error) {
	out, err := iCollectionsVault.abi.Unpack("lendingManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256 maxAssets)
func (iCollectionsVault *ICollectionsVault) PackMaxDeposit(receiver common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("maxDeposit", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMaxDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x402d267d.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256 maxAssets)
func (iCollectionsVault *ICollectionsVault) UnpackMaxDeposit(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("maxDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMaxMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc63d75b6.
//
// Solidity: function maxMint(address receiver) view returns(uint256 maxShares)
func (iCollectionsVault *ICollectionsVault) PackMaxMint(receiver common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("maxMint", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMaxMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc63d75b6.
//
// Solidity: function maxMint(address receiver) view returns(uint256 maxShares)
func (iCollectionsVault *ICollectionsVault) UnpackMaxMint(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("maxMint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMaxRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905777e.
//
// Solidity: function maxRedeem(address owner) view returns(uint256 maxShares)
func (iCollectionsVault *ICollectionsVault) PackMaxRedeem(owner common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("maxRedeem", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMaxRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd905777e.
//
// Solidity: function maxRedeem(address owner) view returns(uint256 maxShares)
func (iCollectionsVault *ICollectionsVault) UnpackMaxRedeem(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("maxRedeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMaxWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce96cb77.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256 maxAssets)
func (iCollectionsVault *ICollectionsVault) PackMaxWithdraw(owner common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("maxWithdraw", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMaxWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce96cb77.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256 maxAssets)
func (iCollectionsVault *ICollectionsVault) UnpackMaxWithdraw(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("maxWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackMint(shares *big.Int, receiver common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("mint", shares, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackMint(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("mint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackMintForCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84f1bc62.
//
// Solidity: function mintForCollection(uint256 shares, address receiver, address collectionAddress) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackMintForCollection(shares *big.Int, receiver common.Address, collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("mintForCollection", shares, receiver, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMintForCollection is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84f1bc62.
//
// Solidity: function mintForCollection(uint256 shares, address receiver, address collectionAddress) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackMintForCollection(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("mintForCollection", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (iCollectionsVault *ICollectionsVault) PackName() []byte {
	enc, err := iCollectionsVault.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (iCollectionsVault *ICollectionsVault) UnpackName(data []byte) (string, error) {
	out, err := iCollectionsVault.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackPreviewDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackPreviewDeposit(assets *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("previewDeposit", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPreviewDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackPreviewDeposit(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("previewDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackPreviewMint(shares *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("previewMint", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPreviewMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackPreviewMint(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("previewMint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackPreviewRedeem(shares *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("previewRedeem", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPreviewRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackPreviewRedeem(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("previewRedeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackPreviewWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackPreviewWithdraw(assets *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("previewWithdraw", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPreviewWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackPreviewWithdraw(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("previewWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackRecordCollectionBorrowVolume is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6fd298f4.
//
// Solidity: function recordCollectionBorrowVolume(address collectionAddress, uint256 borrowAmount) returns()
func (iCollectionsVault *ICollectionsVault) PackRecordCollectionBorrowVolume(collectionAddress common.Address, borrowAmount *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("recordCollectionBorrowVolume", collectionAddress, borrowAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackRedeem(shares *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("redeem", shares, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackRedeem(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("redeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackRedeemForCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb839eaf8.
//
// Solidity: function redeemForCollection(uint256 shares, address receiver, address owner, address collectionAddress) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) PackRedeemForCollection(shares *big.Int, receiver common.Address, owner common.Address, collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("redeemForCollection", shares, receiver, owner, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRedeemForCollection is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb839eaf8.
//
// Solidity: function redeemForCollection(uint256 shares, address receiver, address owner, address collectionAddress) returns(uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackRedeemForCollection(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("redeemForCollection", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackRepayBorrowBehalf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa553b61.
//
// Solidity: function repayBorrowBehalf(uint256 amount, address borrower) returns()
func (iCollectionsVault *ICollectionsVault) PackRepayBorrowBehalf(amount *big.Int, borrower common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("repayBorrowBehalf", amount, borrower)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRepayBorrowBehalfBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11839451.
//
// Solidity: function repayBorrowBehalfBatch(uint256[] amounts, address[] borrowers, uint256 totalAmount) returns()
func (iCollectionsVault *ICollectionsVault) PackRepayBorrowBehalfBatch(amounts []*big.Int, borrowers []common.Address, totalAmount *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("repayBorrowBehalfBatch", amounts, borrowers, totalAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResetEpochCollectionYieldFlags is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x105804b1.
//
// Solidity: function resetEpochCollectionYieldFlags(uint256 epochId, address[] collections) returns()
func (iCollectionsVault *ICollectionsVault) PackResetEpochCollectionYieldFlags(epochId *big.Int, collections []common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("resetEpochCollectionYieldFlags", epochId, collections)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetCollectionRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaff437a1.
//
// Solidity: function setCollectionRegistry(address _collectionRegistryAddress) returns()
func (iCollectionsVault *ICollectionsVault) PackSetCollectionRegistry(collectionRegistryAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("setCollectionRegistry", collectionRegistryAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetDebtSubsidizer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb94305aa.
//
// Solidity: function setDebtSubsidizer(address _debtSubsidizerAddress) returns()
func (iCollectionsVault *ICollectionsVault) PackSetDebtSubsidizer(debtSubsidizerAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("setDebtSubsidizer", debtSubsidizerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetEpochManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4d59edb.
//
// Solidity: function setEpochManager(address _epochManagerAddress) returns()
func (iCollectionsVault *ICollectionsVault) PackSetEpochManager(epochManagerAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("setEpochManager", epochManagerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetLendingManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc869d0ed.
//
// Solidity: function setLendingManager(address _lendingManagerAddress) returns()
func (iCollectionsVault *ICollectionsVault) PackSetLendingManager(lendingManagerAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("setLendingManager", lendingManagerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (iCollectionsVault *ICollectionsVault) PackSymbol() []byte {
	enc, err := iCollectionsVault.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (iCollectionsVault *ICollectionsVault) UnpackSymbol(data []byte) (string, error) {
	out, err := iCollectionsVault.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256 totalManagedAssets)
func (iCollectionsVault *ICollectionsVault) PackTotalAssets() []byte {
	enc, err := iCollectionsVault.abi.Pack("totalAssets")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256 totalManagedAssets)
func (iCollectionsVault *ICollectionsVault) UnpackTotalAssets(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("totalAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalAssetsDeposited is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c63c2da.
//
// Solidity: function totalAssetsDeposited() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackTotalAssetsDeposited() []byte {
	enc, err := iCollectionsVault.abi.Pack("totalAssetsDeposited")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalAssetsDeposited is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6c63c2da.
//
// Solidity: function totalAssetsDeposited() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackTotalAssetsDeposited(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("totalAssetsDeposited", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalCollectionYieldShareBps is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfbc0d706.
//
// Solidity: function totalCollectionYieldShareBps() view returns(uint16)
func (iCollectionsVault *ICollectionsVault) PackTotalCollectionYieldShareBps() []byte {
	enc, err := iCollectionsVault.abi.Pack("totalCollectionYieldShareBps")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalCollectionYieldShareBps is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfbc0d706.
//
// Solidity: function totalCollectionYieldShareBps() view returns(uint16)
func (iCollectionsVault *ICollectionsVault) UnpackTotalCollectionYieldShareBps(data []byte) (uint16, error) {
	out, err := iCollectionsVault.abi.Unpack("totalCollectionYieldShareBps", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, err
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackTotalSupply() []byte {
	enc, err := iCollectionsVault.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalYieldAllocated is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x160e5d9a.
//
// Solidity: function totalYieldAllocated() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackTotalYieldAllocated() []byte {
	enc, err := iCollectionsVault.abi.Pack("totalYieldAllocated")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalYieldAllocated is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x160e5d9a.
//
// Solidity: function totalYieldAllocated() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackTotalYieldAllocated(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("totalYieldAllocated", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalYieldReserved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b2373f5.
//
// Solidity: function totalYieldReserved() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) PackTotalYieldReserved() []byte {
	enc, err := iCollectionsVault.abi.Pack("totalYieldReserved")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalYieldReserved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4b2373f5.
//
// Solidity: function totalYieldReserved() view returns(uint256)
func (iCollectionsVault *ICollectionsVault) UnpackTotalYieldReserved(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("totalYieldReserved", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (iCollectionsVault *ICollectionsVault) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (iCollectionsVault *ICollectionsVault) UnpackTransfer(data []byte) (bool, error) {
	out, err := iCollectionsVault.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTransferForCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa071036f.
//
// Solidity: function transferForCollection(address collectionAddress, address to, uint256 assets) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackTransferForCollection(collectionAddress common.Address, to common.Address, assets *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("transferForCollection", collectionAddress, to, assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransferForCollection is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa071036f.
//
// Solidity: function transferForCollection(address collectionAddress, address to, uint256 assets) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackTransferForCollection(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("transferForCollection", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (iCollectionsVault *ICollectionsVault) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (iCollectionsVault *ICollectionsVault) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := iCollectionsVault.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackUnderlying is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (iCollectionsVault *ICollectionsVault) PackUnderlying() []byte {
	enc, err := iCollectionsVault.abi.Pack("underlying")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUnderlying is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (iCollectionsVault *ICollectionsVault) UnpackUnderlying(data []byte) (common.Address, error) {
	out, err := iCollectionsVault.abi.Unpack("underlying", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUpdateCollectionPerformanceScore is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11104786.
//
// Solidity: function updateCollectionPerformanceScore(address collectionAddress, uint256 score) returns()
func (iCollectionsVault *ICollectionsVault) PackUpdateCollectionPerformanceScore(collectionAddress common.Address, score *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("updateCollectionPerformanceScore", collectionAddress, score)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackValidateCumulativeClaims is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cb9b5b6.
//
// Solidity: function validateCumulativeClaims(uint256 totalClaimedAmount) view returns(bool)
func (iCollectionsVault *ICollectionsVault) PackValidateCumulativeClaims(totalClaimedAmount *big.Int) []byte {
	enc, err := iCollectionsVault.abi.Pack("validateCumulativeClaims", totalClaimedAmount)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackValidateCumulativeClaims is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cb9b5b6.
//
// Solidity: function validateCumulativeClaims(uint256 totalClaimedAmount) view returns(bool)
func (iCollectionsVault *ICollectionsVault) UnpackValidateCumulativeClaims(data []byte) (bool, error) {
	out, err := iCollectionsVault.abi.Unpack("validateCumulativeClaims", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("withdraw", assets, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb460af94.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackWithdraw(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("withdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackWithdrawForCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xda721076.
//
// Solidity: function withdrawForCollection(uint256 assets, address receiver, address owner, address collectionAddress) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) PackWithdrawForCollection(assets *big.Int, receiver common.Address, owner common.Address, collectionAddress common.Address) []byte {
	enc, err := iCollectionsVault.abi.Pack("withdrawForCollection", assets, receiver, owner, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdrawForCollection is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xda721076.
//
// Solidity: function withdrawForCollection(uint256 assets, address receiver, address owner, address collectionAddress) returns(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackWithdrawForCollection(data []byte) (*big.Int, error) {
	out, err := iCollectionsVault.abi.Unpack("withdrawForCollection", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// ICollectionsVaultApproval represents a Approval event raised by the ICollectionsVault contract.
type ICollectionsVaultApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultApproval) ContractEventName() string {
	return ICollectionsVaultApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (iCollectionsVault *ICollectionsVault) UnpackApprovalEvent(log *types.Log) (*ICollectionsVaultApproval, error) {
	event := "Approval"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultApproval)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionAccessGranted represents a CollectionAccessGranted event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionAccessGranted struct {
	Collection common.Address
	Operator   common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionAccessGrantedEventName = "CollectionAccessGranted"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionAccessGranted) ContractEventName() string {
	return ICollectionsVaultCollectionAccessGrantedEventName
}

// UnpackCollectionAccessGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionAccessGranted(address indexed collection, address indexed operator)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionAccessGrantedEvent(log *types.Log) (*ICollectionsVaultCollectionAccessGranted, error) {
	event := "CollectionAccessGranted"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionAccessGranted)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionAccessRevoked represents a CollectionAccessRevoked event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionAccessRevoked struct {
	Collection common.Address
	Operator   common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionAccessRevokedEventName = "CollectionAccessRevoked"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionAccessRevoked) ContractEventName() string {
	return ICollectionsVaultCollectionAccessRevokedEventName
}

// UnpackCollectionAccessRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionAccessRevoked(address indexed collection, address indexed operator)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionAccessRevokedEvent(log *types.Log) (*ICollectionsVaultCollectionAccessRevoked, error) {
	event := "CollectionAccessRevoked"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionAccessRevoked)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionBorrowVolumeUpdated represents a CollectionBorrowVolumeUpdated event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionBorrowVolumeUpdated struct {
	CollectionAddress common.Address
	TotalVolume       *big.Int
	IncrementAmount   *big.Int
	Timestamp         *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionBorrowVolumeUpdatedEventName = "CollectionBorrowVolumeUpdated"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionBorrowVolumeUpdated) ContractEventName() string {
	return ICollectionsVaultCollectionBorrowVolumeUpdatedEventName
}

// UnpackCollectionBorrowVolumeUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionBorrowVolumeUpdated(address indexed collectionAddress, uint256 indexed totalVolume, uint256 indexed incrementAmount, uint256 timestamp)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionBorrowVolumeUpdatedEvent(log *types.Log) (*ICollectionsVaultCollectionBorrowVolumeUpdated, error) {
	event := "CollectionBorrowVolumeUpdated"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionBorrowVolumeUpdated)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionDeposit represents a CollectionDeposit event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionDeposit struct {
	CollectionAddress common.Address
	Caller            common.Address
	Receiver          common.Address
	Assets            *big.Int
	Shares            *big.Int
	CTokenAmount      *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionDepositEventName = "CollectionDeposit"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionDeposit) ContractEventName() string {
	return ICollectionsVaultCollectionDepositEventName
}

// UnpackCollectionDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionDeposit(address indexed collectionAddress, address indexed caller, address indexed receiver, uint256 assets, uint256 shares, uint256 cTokenAmount)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionDepositEvent(log *types.Log) (*ICollectionsVaultCollectionDeposit, error) {
	event := "CollectionDeposit"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionDeposit)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionPerformanceUpdated represents a CollectionPerformanceUpdated event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionPerformanceUpdated struct {
	CollectionAddress common.Address
	PerformanceScore  *big.Int
	Timestamp         *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionPerformanceUpdatedEventName = "CollectionPerformanceUpdated"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionPerformanceUpdated) ContractEventName() string {
	return ICollectionsVaultCollectionPerformanceUpdatedEventName
}

// UnpackCollectionPerformanceUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionPerformanceUpdated(address indexed collectionAddress, uint256 indexed performanceScore, uint256 timestamp)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionPerformanceUpdatedEvent(log *types.Log) (*ICollectionsVaultCollectionPerformanceUpdated, error) {
	event := "CollectionPerformanceUpdated"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionPerformanceUpdated)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionRegistryUpdated represents a CollectionRegistryUpdated event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionRegistryUpdated struct {
	OldRegistry common.Address
	NewRegistry common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionRegistryUpdatedEventName = "CollectionRegistryUpdated"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionRegistryUpdated) ContractEventName() string {
	return ICollectionsVaultCollectionRegistryUpdatedEventName
}

// UnpackCollectionRegistryUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionRegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionRegistryUpdatedEvent(log *types.Log) (*ICollectionsVaultCollectionRegistryUpdated, error) {
	event := "CollectionRegistryUpdated"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionRegistryUpdated)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionTransfer represents a CollectionTransfer event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionTransfer struct {
	CollectionAddress common.Address
	From              common.Address
	To                common.Address
	Assets            *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionTransferEventName = "CollectionTransfer"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionTransfer) ContractEventName() string {
	return ICollectionsVaultCollectionTransferEventName
}

// UnpackCollectionTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionTransfer(address indexed collectionAddress, address indexed from, address indexed to, uint256 assets)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionTransferEvent(log *types.Log) (*ICollectionsVaultCollectionTransfer, error) {
	event := "CollectionTransfer"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionTransfer)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionWithdraw represents a CollectionWithdraw event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionWithdraw struct {
	CollectionAddress common.Address
	Caller            common.Address
	Receiver          common.Address
	Assets            *big.Int
	Shares            *big.Int
	CTokenAmount      *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionWithdrawEventName = "CollectionWithdraw"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionWithdraw) ContractEventName() string {
	return ICollectionsVaultCollectionWithdrawEventName
}

// UnpackCollectionWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionWithdraw(address indexed collectionAddress, address indexed caller, address indexed receiver, uint256 assets, uint256 shares, uint256 cTokenAmount)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionWithdrawEvent(log *types.Log) (*ICollectionsVaultCollectionWithdraw, error) {
	event := "CollectionWithdraw"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionWithdraw)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionYieldAccrued represents a CollectionYieldAccrued event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionYieldAccrued struct {
	CollectionAddress       common.Address
	YieldAccrued            *big.Int
	NewTotalDeposits        *big.Int
	GlobalIndex             *big.Int
	PreviousCollectionIndex *big.Int
	Raw                     *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionYieldAccruedEventName = "CollectionYieldAccrued"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionYieldAccrued) ContractEventName() string {
	return ICollectionsVaultCollectionYieldAccruedEventName
}

// UnpackCollectionYieldAccruedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionYieldAccrued(address indexed collectionAddress, uint256 yieldAccrued, uint256 newTotalDeposits, uint256 globalIndex, uint256 previousCollectionIndex)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionYieldAccruedEvent(log *types.Log) (*ICollectionsVaultCollectionYieldAccrued, error) {
	event := "CollectionYieldAccrued"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionYieldAccrued)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionYieldAppliedForEpoch represents a CollectionYieldAppliedForEpoch event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionYieldAppliedForEpoch struct {
	EpochId              *big.Int
	Collection           common.Address
	YieldSharePercentage uint16
	YieldAdded           *big.Int
	NewTotalDeposits     *big.Int
	Raw                  *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionYieldAppliedForEpochEventName = "CollectionYieldAppliedForEpoch"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionYieldAppliedForEpoch) ContractEventName() string {
	return ICollectionsVaultCollectionYieldAppliedForEpochEventName
}

// UnpackCollectionYieldAppliedForEpochEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionYieldAppliedForEpoch(uint256 indexed epochId, address indexed collection, uint16 yieldSharePercentage, uint256 yieldAdded, uint256 newTotalDeposits)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionYieldAppliedForEpochEvent(log *types.Log) (*ICollectionsVaultCollectionYieldAppliedForEpoch, error) {
	event := "CollectionYieldAppliedForEpoch"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionYieldAppliedForEpoch)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionYieldGenerated represents a CollectionYieldGenerated event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionYieldGenerated struct {
	CollectionAddress common.Address
	YieldAmount       *big.Int
	Timestamp         *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionYieldGeneratedEventName = "CollectionYieldGenerated"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionYieldGenerated) ContractEventName() string {
	return ICollectionsVaultCollectionYieldGeneratedEventName
}

// UnpackCollectionYieldGeneratedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionYieldGenerated(address indexed collectionAddress, uint256 indexed yieldAmount, uint256 indexed timestamp)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionYieldGeneratedEvent(log *types.Log) (*ICollectionsVaultCollectionYieldGenerated, error) {
	event := "CollectionYieldGenerated"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionYieldGenerated)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultCollectionYieldIndexed represents a CollectionYieldIndexed event raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionYieldIndexed struct {
	CollectionAddress common.Address
	EpochId           *big.Int
	Assets            *big.Int
	Shares            *big.Int
	CTokenAmount      *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultCollectionYieldIndexedEventName = "CollectionYieldIndexed"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultCollectionYieldIndexed) ContractEventName() string {
	return ICollectionsVaultCollectionYieldIndexedEventName
}

// UnpackCollectionYieldIndexedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionYieldIndexed(address indexed collectionAddress, uint256 indexed epochId, uint256 assets, uint256 shares, uint256 cTokenAmount)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionYieldIndexedEvent(log *types.Log) (*ICollectionsVaultCollectionYieldIndexed, error) {
	event := "CollectionYieldIndexed"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultCollectionYieldIndexed)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultDeposit represents a Deposit event raised by the ICollectionsVault contract.
type ICollectionsVaultDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultDeposit) ContractEventName() string {
	return ICollectionsVaultDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackDepositEvent(log *types.Log) (*ICollectionsVaultDeposit, error) {
	event := "Deposit"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultDeposit)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultEpochManagerCallFailed represents a EpochManagerCallFailed event raised by the ICollectionsVault contract.
type ICollectionsVaultEpochManagerCallFailed struct {
	VaultAddress common.Address
	EpochId      *big.Int
	Amount       *big.Int
	Reason       string
	Raw          *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultEpochManagerCallFailedEventName = "EpochManagerCallFailed"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultEpochManagerCallFailed) ContractEventName() string {
	return ICollectionsVaultEpochManagerCallFailedEventName
}

// UnpackEpochManagerCallFailedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochManagerCallFailed(address indexed vaultAddress, uint256 indexed epochId, uint256 indexed amount, string reason)
func (iCollectionsVault *ICollectionsVault) UnpackEpochManagerCallFailedEvent(log *types.Log) (*ICollectionsVaultEpochManagerCallFailed, error) {
	event := "EpochManagerCallFailed"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultEpochManagerCallFailed)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultEpochManagerCallUnavailable represents a EpochManagerCallUnavailable event raised by the ICollectionsVault contract.
type ICollectionsVaultEpochManagerCallUnavailable struct {
	VaultAddress common.Address
	FunctionName string
	Reason       string
	Raw          *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultEpochManagerCallUnavailableEventName = "EpochManagerCallUnavailable"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultEpochManagerCallUnavailable) ContractEventName() string {
	return ICollectionsVaultEpochManagerCallUnavailableEventName
}

// UnpackEpochManagerCallUnavailableEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EpochManagerCallUnavailable(address indexed vaultAddress, string functionName, string reason)
func (iCollectionsVault *ICollectionsVault) UnpackEpochManagerCallUnavailableEvent(log *types.Log) (*ICollectionsVaultEpochManagerCallUnavailable, error) {
	event := "EpochManagerCallUnavailable"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultEpochManagerCallUnavailable)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultLendingManagerCallFailed represents a LendingManagerCallFailed event raised by the ICollectionsVault contract.
type ICollectionsVaultLendingManagerCallFailed struct {
	VaultAddress common.Address
	Operation    string
	Amount       *big.Int
	Reason       string
	Raw          *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultLendingManagerCallFailedEventName = "LendingManagerCallFailed"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultLendingManagerCallFailed) ContractEventName() string {
	return ICollectionsVaultLendingManagerCallFailedEventName
}

// UnpackLendingManagerCallFailedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LendingManagerCallFailed(address indexed vaultAddress, string operation, uint256 amount, string reason)
func (iCollectionsVault *ICollectionsVault) UnpackLendingManagerCallFailedEvent(log *types.Log) (*ICollectionsVaultLendingManagerCallFailed, error) {
	event := "LendingManagerCallFailed"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultLendingManagerCallFailed)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultLendingManagerChanged represents a LendingManagerChanged event raised by the ICollectionsVault contract.
type ICollectionsVaultLendingManagerChanged struct {
	OldLendingManager common.Address
	NewLendingManager common.Address
	ChangedBy         common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultLendingManagerChangedEventName = "LendingManagerChanged"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultLendingManagerChanged) ContractEventName() string {
	return ICollectionsVaultLendingManagerChangedEventName
}

// UnpackLendingManagerChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LendingManagerChanged(address indexed oldLendingManager, address indexed newLendingManager, address indexed changedBy)
func (iCollectionsVault *ICollectionsVault) UnpackLendingManagerChangedEvent(log *types.Log) (*ICollectionsVaultLendingManagerChanged, error) {
	event := "LendingManagerChanged"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultLendingManagerChanged)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultTransfer represents a Transfer event raised by the ICollectionsVault contract.
type ICollectionsVaultTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultTransfer) ContractEventName() string {
	return ICollectionsVaultTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (iCollectionsVault *ICollectionsVault) UnpackTransferEvent(log *types.Log) (*ICollectionsVaultTransfer, error) {
	event := "Transfer"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultTransfer)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultVaultYieldAllocatedToEpoch represents a VaultYieldAllocatedToEpoch event raised by the ICollectionsVault contract.
type ICollectionsVaultVaultYieldAllocatedToEpoch struct {
	EpochId *big.Int
	Amount  *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultVaultYieldAllocatedToEpochEventName = "VaultYieldAllocatedToEpoch"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultVaultYieldAllocatedToEpoch) ContractEventName() string {
	return ICollectionsVaultVaultYieldAllocatedToEpochEventName
}

// UnpackVaultYieldAllocatedToEpochEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultYieldAllocatedToEpoch(uint256 indexed epochId, uint256 amount)
func (iCollectionsVault *ICollectionsVault) UnpackVaultYieldAllocatedToEpochEvent(log *types.Log) (*ICollectionsVaultVaultYieldAllocatedToEpoch, error) {
	event := "VaultYieldAllocatedToEpoch"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultVaultYieldAllocatedToEpoch)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultWithdraw represents a Withdraw event raised by the ICollectionsVault contract.
type ICollectionsVaultWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultWithdraw) ContractEventName() string {
	return ICollectionsVaultWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackWithdrawEvent(log *types.Log) (*ICollectionsVaultWithdraw, error) {
	event := "Withdraw"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultWithdraw)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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

// ICollectionsVaultYieldBatchRepaid represents a YieldBatchRepaid event raised by the ICollectionsVault contract.
type ICollectionsVaultYieldBatchRepaid struct {
	TotalAmount *big.Int
	Recipient   common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const ICollectionsVaultYieldBatchRepaidEventName = "YieldBatchRepaid"

// ContractEventName returns the user-defined event name.
func (ICollectionsVaultYieldBatchRepaid) ContractEventName() string {
	return ICollectionsVaultYieldBatchRepaidEventName
}

// UnpackYieldBatchRepaidEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event YieldBatchRepaid(uint256 totalAmount, address indexed recipient)
func (iCollectionsVault *ICollectionsVault) UnpackYieldBatchRepaidEvent(log *types.Log) (*ICollectionsVaultYieldBatchRepaid, error) {
	event := "YieldBatchRepaid"
	if log.Topics[0] != iCollectionsVault.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ICollectionsVaultYieldBatchRepaid)
	if len(log.Data) > 0 {
		if err := iCollectionsVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iCollectionsVault.abi.Events[event].Inputs {
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
func (iCollectionsVault *ICollectionsVault) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["AddressZero"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackAddressZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["AllocExceedsAvail"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackAllocExceedsAvailError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["AllocationAmountZero"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackAllocationAmountZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["AllocationUnderflow"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackAllocationUnderflowError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["BatchSizeExceedsLimit"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackBatchSizeExceedsLimitError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["CollectionInsufficientBalance"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackCollectionInsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["CollectionNotRegistered"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackCollectionNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["EpochManagerAllocationFailed"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackEpochManagerAllocationFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["EpochManagerNotSet"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackEpochManagerNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["EpochManagerUnavailable"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackEpochManagerUnavailableError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["ExcessiveYieldAmount"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackExcessiveYieldAmountError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["FunctionDisabledUse"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackFunctionDisabledUseError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["InsufficientBalanceInProtocol"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackInsufficientBalanceInProtocolError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["InvalidEpochId"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackInvalidEpochIdError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["LendingManagerDepositFailed"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackLendingManagerDepositFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["LendingManagerMismatch"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackLendingManagerMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["LendingManagerWithdrawFailed"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackLendingManagerWithdrawFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["NoActiveEpoch"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackNoActiveEpochError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["NoCumulativeYield"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackNoCumulativeYieldError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["PerfScoreExceedsMax"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackPerfScoreExceedsMaxError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["RedeemRoundsToZero"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackRedeemRoundsToZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["RepayBorrowFailed"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackRepayBorrowFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["ReqAmountExceeds"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackReqAmountExceedsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["ShareBalanceUnderflow"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackShareBalanceUnderflowError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["TotalAllocExceeds"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackTotalAllocExceedsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["UnauthorizedCollectionAccess"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackUnauthorizedCollectionAccessError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["VaultInsufficientBalancePostLMWithdraw"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackVaultInsufficientBalancePostLMWithdrawError(raw[4:])
	}
	if bytes.Equal(raw[:4], iCollectionsVault.abi.Errors["YieldAlreadyApplied"].ID.Bytes()[:4]) {
		return iCollectionsVault.UnpackYieldAlreadyAppliedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ICollectionsVaultAddressZero represents a AddressZero error raised by the ICollectionsVault contract.
type ICollectionsVaultAddressZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressZero()
func ICollectionsVaultAddressZeroErrorID() common.Hash {
	return common.HexToHash("0x9fabe1c19979afc45ec7efec1bde2c38021c590a0ce42965cf55b3f518197f02")
}

// UnpackAddressZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressZero()
func (iCollectionsVault *ICollectionsVault) UnpackAddressZeroError(raw []byte) (*ICollectionsVaultAddressZero, error) {
	out := new(ICollectionsVaultAddressZero)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "AddressZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultAllocExceedsAvail represents a AllocExceedsAvail error raised by the ICollectionsVault contract.
type ICollectionsVaultAllocExceedsAvail struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AllocExceedsAvail()
func ICollectionsVaultAllocExceedsAvailErrorID() common.Hash {
	return common.HexToHash("0xfd02ddcdacbc8fb269d600350badffeadcaabda459d26dc37061f86efe1efeff")
}

// UnpackAllocExceedsAvailError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AllocExceedsAvail()
func (iCollectionsVault *ICollectionsVault) UnpackAllocExceedsAvailError(raw []byte) (*ICollectionsVaultAllocExceedsAvail, error) {
	out := new(ICollectionsVaultAllocExceedsAvail)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "AllocExceedsAvail", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultAllocationAmountZero represents a AllocationAmountZero error raised by the ICollectionsVault contract.
type ICollectionsVaultAllocationAmountZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AllocationAmountZero()
func ICollectionsVaultAllocationAmountZeroErrorID() common.Hash {
	return common.HexToHash("0x5d24732f466b44890097ac2f5498fed70abc9c6e69c48c607890a76ffc295b8d")
}

// UnpackAllocationAmountZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AllocationAmountZero()
func (iCollectionsVault *ICollectionsVault) UnpackAllocationAmountZeroError(raw []byte) (*ICollectionsVaultAllocationAmountZero, error) {
	out := new(ICollectionsVaultAllocationAmountZero)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "AllocationAmountZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultAllocationUnderflow represents a AllocationUnderflow error raised by the ICollectionsVault contract.
type ICollectionsVaultAllocationUnderflow struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AllocationUnderflow()
func ICollectionsVaultAllocationUnderflowErrorID() common.Hash {
	return common.HexToHash("0x3e104230e5dbc3350a5e942a223874ee444e1c73ea7e4224eca3ae3accc35e1d")
}

// UnpackAllocationUnderflowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AllocationUnderflow()
func (iCollectionsVault *ICollectionsVault) UnpackAllocationUnderflowError(raw []byte) (*ICollectionsVaultAllocationUnderflow, error) {
	out := new(ICollectionsVaultAllocationUnderflow)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "AllocationUnderflow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultBatchSizeExceedsLimit represents a BatchSizeExceedsLimit error raised by the ICollectionsVault contract.
type ICollectionsVaultBatchSizeExceedsLimit struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchSizeExceedsLimit()
func ICollectionsVaultBatchSizeExceedsLimitErrorID() common.Hash {
	return common.HexToHash("0xaa19ded49172c7d8053f3214b858bdc31179af88ff4cbaa0007f7dc3853f6dee")
}

// UnpackBatchSizeExceedsLimitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchSizeExceedsLimit()
func (iCollectionsVault *ICollectionsVault) UnpackBatchSizeExceedsLimitError(raw []byte) (*ICollectionsVaultBatchSizeExceedsLimit, error) {
	out := new(ICollectionsVaultBatchSizeExceedsLimit)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "BatchSizeExceedsLimit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultCollectionInsufficientBalance represents a CollectionInsufficientBalance error raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionInsufficientBalance struct {
	CollectionAddress common.Address
	Requested         *big.Int
	Available         *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionInsufficientBalance(address collectionAddress, uint256 requested, uint256 available)
func ICollectionsVaultCollectionInsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xae86887122d5896697736925442a7e04aa3165774f62365d5a4a37489d9196e9")
}

// UnpackCollectionInsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionInsufficientBalance(address collectionAddress, uint256 requested, uint256 available)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionInsufficientBalanceError(raw []byte) (*ICollectionsVaultCollectionInsufficientBalance, error) {
	out := new(ICollectionsVaultCollectionInsufficientBalance)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "CollectionInsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultCollectionNotRegistered represents a CollectionNotRegistered error raised by the ICollectionsVault contract.
type ICollectionsVaultCollectionNotRegistered struct {
	CollectionAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionNotRegistered(address collectionAddress)
func ICollectionsVaultCollectionNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0xb156116221c5df2e93ed39baa0dfb8cb70feafbc6c39d9ed9c95e8f716fc7483")
}

// UnpackCollectionNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionNotRegistered(address collectionAddress)
func (iCollectionsVault *ICollectionsVault) UnpackCollectionNotRegisteredError(raw []byte) (*ICollectionsVaultCollectionNotRegistered, error) {
	out := new(ICollectionsVaultCollectionNotRegistered)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "CollectionNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultEpochManagerAllocationFailed represents a EpochManagerAllocationFailed error raised by the ICollectionsVault contract.
type ICollectionsVaultEpochManagerAllocationFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManagerAllocationFailed()
func ICollectionsVaultEpochManagerAllocationFailedErrorID() common.Hash {
	return common.HexToHash("0x6d9a1215b679d84d849a2dbd0a516e65016fc85fce97d0afb5930203b77f4cca")
}

// UnpackEpochManagerAllocationFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManagerAllocationFailed()
func (iCollectionsVault *ICollectionsVault) UnpackEpochManagerAllocationFailedError(raw []byte) (*ICollectionsVaultEpochManagerAllocationFailed, error) {
	out := new(ICollectionsVaultEpochManagerAllocationFailed)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "EpochManagerAllocationFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultEpochManagerNotSet represents a EpochManagerNotSet error raised by the ICollectionsVault contract.
type ICollectionsVaultEpochManagerNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManagerNotSet()
func ICollectionsVaultEpochManagerNotSetErrorID() common.Hash {
	return common.HexToHash("0x23afc25f012f52235fc3f433f9d218e0731b259ee389f073a698239f7e1738ba")
}

// UnpackEpochManagerNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManagerNotSet()
func (iCollectionsVault *ICollectionsVault) UnpackEpochManagerNotSetError(raw []byte) (*ICollectionsVaultEpochManagerNotSet, error) {
	out := new(ICollectionsVaultEpochManagerNotSet)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "EpochManagerNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultEpochManagerUnavailable represents a EpochManagerUnavailable error raised by the ICollectionsVault contract.
type ICollectionsVaultEpochManagerUnavailable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EpochManagerUnavailable()
func ICollectionsVaultEpochManagerUnavailableErrorID() common.Hash {
	return common.HexToHash("0x874e4703918695585e263eaf32fbf262ca190bba9ea93cdcadd2653368b09fb9")
}

// UnpackEpochManagerUnavailableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EpochManagerUnavailable()
func (iCollectionsVault *ICollectionsVault) UnpackEpochManagerUnavailableError(raw []byte) (*ICollectionsVaultEpochManagerUnavailable, error) {
	out := new(ICollectionsVaultEpochManagerUnavailable)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "EpochManagerUnavailable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultExcessiveYieldAmount represents a ExcessiveYieldAmount error raised by the ICollectionsVault contract.
type ICollectionsVaultExcessiveYieldAmount struct {
	Collection common.Address
	Requested  *big.Int
	MaxAllowed *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ExcessiveYieldAmount(address collection, uint256 requested, uint256 maxAllowed)
func ICollectionsVaultExcessiveYieldAmountErrorID() common.Hash {
	return common.HexToHash("0x6b53545e228d6c01653606073e1ad3bdda9f2825e37964146a02602a397ce116")
}

// UnpackExcessiveYieldAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ExcessiveYieldAmount(address collection, uint256 requested, uint256 maxAllowed)
func (iCollectionsVault *ICollectionsVault) UnpackExcessiveYieldAmountError(raw []byte) (*ICollectionsVaultExcessiveYieldAmount, error) {
	out := new(ICollectionsVaultExcessiveYieldAmount)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "ExcessiveYieldAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultFunctionDisabledUse represents a FunctionDisabledUse error raised by the ICollectionsVault contract.
type ICollectionsVaultFunctionDisabledUse struct {
	FunctionName string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FunctionDisabledUse(string functionName)
func ICollectionsVaultFunctionDisabledUseErrorID() common.Hash {
	return common.HexToHash("0xf18e438ed85786dda853be37960a83db2d5c40387bc83bd253b4e4b4d2c120a0")
}

// UnpackFunctionDisabledUseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FunctionDisabledUse(string functionName)
func (iCollectionsVault *ICollectionsVault) UnpackFunctionDisabledUseError(raw []byte) (*ICollectionsVaultFunctionDisabledUse, error) {
	out := new(ICollectionsVaultFunctionDisabledUse)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "FunctionDisabledUse", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultInsufficientBalanceInProtocol represents a InsufficientBalanceInProtocol error raised by the ICollectionsVault contract.
type ICollectionsVaultInsufficientBalanceInProtocol struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientBalanceInProtocol()
func ICollectionsVaultInsufficientBalanceInProtocolErrorID() common.Hash {
	return common.HexToHash("0x7e65ad2ca86fa6ae58f772c31968382bed67947d6803d528a3c28871cbee7948")
}

// UnpackInsufficientBalanceInProtocolError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientBalanceInProtocol()
func (iCollectionsVault *ICollectionsVault) UnpackInsufficientBalanceInProtocolError(raw []byte) (*ICollectionsVaultInsufficientBalanceInProtocol, error) {
	out := new(ICollectionsVaultInsufficientBalanceInProtocol)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "InsufficientBalanceInProtocol", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultInvalidEpochId represents a InvalidEpochId error raised by the ICollectionsVault contract.
type ICollectionsVaultInvalidEpochId struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidEpochId()
func ICollectionsVaultInvalidEpochIdErrorID() common.Hash {
	return common.HexToHash("0x395a9c30368af6a0448306d5fa8d20c3d0f1d0460103edce8e5b5cd2d69fc4e6")
}

// UnpackInvalidEpochIdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidEpochId()
func (iCollectionsVault *ICollectionsVault) UnpackInvalidEpochIdError(raw []byte) (*ICollectionsVaultInvalidEpochId, error) {
	out := new(ICollectionsVaultInvalidEpochId)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "InvalidEpochId", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultLendingManagerDepositFailed represents a LendingManagerDepositFailed error raised by the ICollectionsVault contract.
type ICollectionsVaultLendingManagerDepositFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerDepositFailed()
func ICollectionsVaultLendingManagerDepositFailedErrorID() common.Hash {
	return common.HexToHash("0xd8449a590b53d0baaf9e162be52915834b7a3c2387847ce9014fb288e92e2f9d")
}

// UnpackLendingManagerDepositFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerDepositFailed()
func (iCollectionsVault *ICollectionsVault) UnpackLendingManagerDepositFailedError(raw []byte) (*ICollectionsVaultLendingManagerDepositFailed, error) {
	out := new(ICollectionsVaultLendingManagerDepositFailed)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "LendingManagerDepositFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultLendingManagerMismatch represents a LendingManagerMismatch error raised by the ICollectionsVault contract.
type ICollectionsVaultLendingManagerMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerMismatch()
func ICollectionsVaultLendingManagerMismatchErrorID() common.Hash {
	return common.HexToHash("0x75de435863935cbccbbd808ec2d19cb84aab15cf7da8b2c7e92461f6b3cf0ecd")
}

// UnpackLendingManagerMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerMismatch()
func (iCollectionsVault *ICollectionsVault) UnpackLendingManagerMismatchError(raw []byte) (*ICollectionsVaultLendingManagerMismatch, error) {
	out := new(ICollectionsVaultLendingManagerMismatch)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "LendingManagerMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultLendingManagerWithdrawFailed represents a LendingManagerWithdrawFailed error raised by the ICollectionsVault contract.
type ICollectionsVaultLendingManagerWithdrawFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerWithdrawFailed()
func ICollectionsVaultLendingManagerWithdrawFailedErrorID() common.Hash {
	return common.HexToHash("0x8b78daf7fd34c46828095f2405bfd3ccd4d1c5f8be85c8e6403cc296f4c2c783")
}

// UnpackLendingManagerWithdrawFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerWithdrawFailed()
func (iCollectionsVault *ICollectionsVault) UnpackLendingManagerWithdrawFailedError(raw []byte) (*ICollectionsVaultLendingManagerWithdrawFailed, error) {
	out := new(ICollectionsVaultLendingManagerWithdrawFailed)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "LendingManagerWithdrawFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultNoActiveEpoch represents a NoActiveEpoch error raised by the ICollectionsVault contract.
type ICollectionsVaultNoActiveEpoch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoActiveEpoch()
func ICollectionsVaultNoActiveEpochErrorID() common.Hash {
	return common.HexToHash("0xc7bbad56162f02d20da7b0e54734998a45df7dc8b1ac35366d549639b5b7cbe6")
}

// UnpackNoActiveEpochError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoActiveEpoch()
func (iCollectionsVault *ICollectionsVault) UnpackNoActiveEpochError(raw []byte) (*ICollectionsVaultNoActiveEpoch, error) {
	out := new(ICollectionsVaultNoActiveEpoch)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "NoActiveEpoch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultNoCumulativeYield represents a NoCumulativeYield error raised by the ICollectionsVault contract.
type ICollectionsVaultNoCumulativeYield struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoCumulativeYield()
func ICollectionsVaultNoCumulativeYieldErrorID() common.Hash {
	return common.HexToHash("0x7c1f1861e48b0ad868c396608e1fc5cfc64bb60676742b7af522b96931fb216e")
}

// UnpackNoCumulativeYieldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoCumulativeYield()
func (iCollectionsVault *ICollectionsVault) UnpackNoCumulativeYieldError(raw []byte) (*ICollectionsVaultNoCumulativeYield, error) {
	out := new(ICollectionsVaultNoCumulativeYield)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "NoCumulativeYield", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultPerfScoreExceedsMax represents a PerfScoreExceedsMax error raised by the ICollectionsVault contract.
type ICollectionsVaultPerfScoreExceedsMax struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PerfScoreExceedsMax()
func ICollectionsVaultPerfScoreExceedsMaxErrorID() common.Hash {
	return common.HexToHash("0x8b3f501f6eee0688b6d4ba392f4885a26d8de48e4db18b4ec70c0059c2ff11ea")
}

// UnpackPerfScoreExceedsMaxError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PerfScoreExceedsMax()
func (iCollectionsVault *ICollectionsVault) UnpackPerfScoreExceedsMaxError(raw []byte) (*ICollectionsVaultPerfScoreExceedsMax, error) {
	out := new(ICollectionsVaultPerfScoreExceedsMax)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "PerfScoreExceedsMax", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultRedeemRoundsToZero represents a RedeemRoundsToZero error raised by the ICollectionsVault contract.
type ICollectionsVaultRedeemRoundsToZero struct {
	Shares *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RedeemRoundsToZero(uint256 shares)
func ICollectionsVaultRedeemRoundsToZeroErrorID() common.Hash {
	return common.HexToHash("0x7312b5e6b907e0d1bf8d8b40a749f8dbd20150581ac561685d69810c49b2c97b")
}

// UnpackRedeemRoundsToZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RedeemRoundsToZero(uint256 shares)
func (iCollectionsVault *ICollectionsVault) UnpackRedeemRoundsToZeroError(raw []byte) (*ICollectionsVaultRedeemRoundsToZero, error) {
	out := new(ICollectionsVaultRedeemRoundsToZero)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "RedeemRoundsToZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultRepayBorrowFailed represents a RepayBorrowFailed error raised by the ICollectionsVault contract.
type ICollectionsVaultRepayBorrowFailed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RepayBorrowFailed()
func ICollectionsVaultRepayBorrowFailedErrorID() common.Hash {
	return common.HexToHash("0xb69b1dedc876e8deac3f6fe824ee78595fda6d07fed9a66cf5b063b31ea21188")
}

// UnpackRepayBorrowFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RepayBorrowFailed()
func (iCollectionsVault *ICollectionsVault) UnpackRepayBorrowFailedError(raw []byte) (*ICollectionsVaultRepayBorrowFailed, error) {
	out := new(ICollectionsVaultRepayBorrowFailed)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "RepayBorrowFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultReqAmountExceeds represents a ReqAmountExceeds error raised by the ICollectionsVault contract.
type ICollectionsVaultReqAmountExceeds struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReqAmountExceeds()
func ICollectionsVaultReqAmountExceedsErrorID() common.Hash {
	return common.HexToHash("0x5975e7080a0031ac951021fa0d0c0f84e1e333ab9ab6eb03f178153d5971eb67")
}

// UnpackReqAmountExceedsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReqAmountExceeds()
func (iCollectionsVault *ICollectionsVault) UnpackReqAmountExceedsError(raw []byte) (*ICollectionsVaultReqAmountExceeds, error) {
	out := new(ICollectionsVaultReqAmountExceeds)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "ReqAmountExceeds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultShareBalanceUnderflow represents a ShareBalanceUnderflow error raised by the ICollectionsVault contract.
type ICollectionsVaultShareBalanceUnderflow struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ShareBalanceUnderflow()
func ICollectionsVaultShareBalanceUnderflowErrorID() common.Hash {
	return common.HexToHash("0xcf26439bbc0b19c6aab21d2c03d00c69ae5bfb2deac10a5e3da6855ca4d66caf")
}

// UnpackShareBalanceUnderflowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ShareBalanceUnderflow()
func (iCollectionsVault *ICollectionsVault) UnpackShareBalanceUnderflowError(raw []byte) (*ICollectionsVaultShareBalanceUnderflow, error) {
	out := new(ICollectionsVaultShareBalanceUnderflow)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "ShareBalanceUnderflow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultTotalAllocExceeds represents a TotalAllocExceeds error raised by the ICollectionsVault contract.
type ICollectionsVaultTotalAllocExceeds struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TotalAllocExceeds()
func ICollectionsVaultTotalAllocExceedsErrorID() common.Hash {
	return common.HexToHash("0xf632a5a1027eda4d660e0c02710621b999558108f51774d9634b9e8d01e64cb2")
}

// UnpackTotalAllocExceedsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TotalAllocExceeds()
func (iCollectionsVault *ICollectionsVault) UnpackTotalAllocExceedsError(raw []byte) (*ICollectionsVaultTotalAllocExceeds, error) {
	out := new(ICollectionsVaultTotalAllocExceeds)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "TotalAllocExceeds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultUnauthorizedCollectionAccess represents a UnauthorizedCollectionAccess error raised by the ICollectionsVault contract.
type ICollectionsVaultUnauthorizedCollectionAccess struct {
	CollectionAddress common.Address
	Operator          common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnauthorizedCollectionAccess(address collectionAddress, address operator)
func ICollectionsVaultUnauthorizedCollectionAccessErrorID() common.Hash {
	return common.HexToHash("0x533d2ab59dda97bc1adf51b17eff135dd3d895ee7473f783081ad597c37b1476")
}

// UnpackUnauthorizedCollectionAccessError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnauthorizedCollectionAccess(address collectionAddress, address operator)
func (iCollectionsVault *ICollectionsVault) UnpackUnauthorizedCollectionAccessError(raw []byte) (*ICollectionsVaultUnauthorizedCollectionAccess, error) {
	out := new(ICollectionsVaultUnauthorizedCollectionAccess)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "UnauthorizedCollectionAccess", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultVaultInsufficientBalancePostLMWithdraw represents a Vault_InsufficientBalancePostLMWithdraw error raised by the ICollectionsVault contract.
type ICollectionsVaultVaultInsufficientBalancePostLMWithdraw struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Vault_InsufficientBalancePostLMWithdraw()
func ICollectionsVaultVaultInsufficientBalancePostLMWithdrawErrorID() common.Hash {
	return common.HexToHash("0x724da007ac24ad369783e82a6944e1d508cc4812a904fb28d056eae41d03f06b")
}

// UnpackVaultInsufficientBalancePostLMWithdrawError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Vault_InsufficientBalancePostLMWithdraw()
func (iCollectionsVault *ICollectionsVault) UnpackVaultInsufficientBalancePostLMWithdrawError(raw []byte) (*ICollectionsVaultVaultInsufficientBalancePostLMWithdraw, error) {
	out := new(ICollectionsVaultVaultInsufficientBalancePostLMWithdraw)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "VaultInsufficientBalancePostLMWithdraw", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ICollectionsVaultYieldAlreadyApplied represents a YieldAlreadyApplied error raised by the ICollectionsVault contract.
type ICollectionsVaultYieldAlreadyApplied struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error YieldAlreadyApplied()
func ICollectionsVaultYieldAlreadyAppliedErrorID() common.Hash {
	return common.HexToHash("0xdfd22544f24cb0d0157208da7e936a56a3481caafcd09c9c1511143f7e6de1e8")
}

// UnpackYieldAlreadyAppliedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error YieldAlreadyApplied()
func (iCollectionsVault *ICollectionsVault) UnpackYieldAlreadyAppliedError(raw []byte) (*ICollectionsVaultYieldAlreadyApplied, error) {
	out := new(ICollectionsVaultYieldAlreadyApplied)
	if err := iCollectionsVault.abi.UnpackIntoInterface(out, "YieldAlreadyApplied", raw); err != nil {
		return nil, err
	}
	return out, nil
}
