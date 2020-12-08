package utils_test

import (
	"testing"

	"github.com/patractlabs/go-patract/utils"
	"github.com/stretchr/testify/require"
)

var (
	address = "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
)

func TestSS58(t *testing.T) {
	id, err := utils.NewAccountIDFromSS58(address)
	require.Nil(t, err)

	add, err := utils.EncodeAccountIDToSS58(id)
	require.Nil(t, err)

	require.Equal(t, address, add)
	t.Logf("Address %x %v", id, id)

}
