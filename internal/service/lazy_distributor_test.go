package service

import (
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/storage"
)

func TestLazyDistributor_BuildMerkleRoot(t *testing.T) {
	ld := &LazyDistributor{}

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
			root := ld.buildMerkleRoot(tt.entries)
			if len(tt.entries) == 0 && root != [32]byte{} {
				t.Error("Empty entries should produce zero root")
			}
			if len(tt.entries) > 0 && root == [32]byte{} {
				t.Error("Non-empty entries should produce non-zero root")
			}
		})
	}
}
