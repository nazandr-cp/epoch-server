package storage

import (
	"testing"

	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseWrapper(t *testing.T) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	t.Run("BadgerDatabase", func(t *testing.T) {
		config := StorageConfig{
			Type: "badger",
			Path: t.TempDir(),
		}

		wrapper, err := NewDatabaseWrapper(config, logger)
		require.NoError(t, err)
		require.NotNil(t, wrapper)

		db := wrapper.GetDB()
		assert.NotNil(t, db)

		err = wrapper.Close()
		assert.NoError(t, err)
	})

	t.Run("UnsupportedType", func(t *testing.T) {
		config := StorageConfig{
			Type: "unsupported",
			Path: t.TempDir(),
		}

		_, err := NewDatabaseWrapper(config, logger)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported storage type")
	})
}
