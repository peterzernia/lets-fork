package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	require := assert.New(t)

	_, err := InitRDB()
	require.NoError(err)
}

func TestGenerateRandomBytes(t *testing.T) {
	require := assert.New(t)

	b, err := GenerateRandomBytes(0)
	require.NoError(err)
	require.Equal(b, []uint8([]byte{}))

	b, err = GenerateRandomBytes(36)
	require.NoError(err)
	require.Equal(len(b), 36)
}
