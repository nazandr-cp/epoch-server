package service

import (
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/merkle"
)

func TestLazyDistributor_BuildMerkleRoot(t *testing.T) {
	ld := &LazyDistributor{
		proofGenerator: merkle.NewProofGenerator(),
	}

	tests := []struct {
		name    string
		entries []storage.MerkleEntry
	}{
		{
			name:    "empty entries",
			entries: []storage.MerkleEntry{},
		},
		{
			name: "single entry",
			entries: []storage.MerkleEntry{
				{Address: "0x1", TotalEarned: big.NewInt(1000000)},
			},
		},
		{
			name: "multiple entries",
			entries: []storage.MerkleEntry{
				{Address: "0x1", TotalEarned: big.NewInt(1000000)},
				{Address: "0x2", TotalEarned: big.NewInt(2000000)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to merkle.Entry format for testing
			entries := make([]merkle.Entry, len(tt.entries))
			for i, entry := range tt.entries {
				entries[i] = merkle.Entry{
					Address:     entry.Address,
					TotalEarned: entry.TotalEarned,
				}
			}
			root := ld.proofGenerator.BuildMerkleRoot(entries)
			if len(tt.entries) == 0 && root != [32]byte{} {
				t.Error("Empty entries should produce zero root")
			}
			if len(tt.entries) > 0 && root == [32]byte{} {
				t.Error("Non-empty entries should produce non-zero root")
			}
		})
	}
}
