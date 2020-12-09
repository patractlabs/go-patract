package rpc_test

import (
	"testing"

	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/stretchr/testify/require"
)

func TestGetContractAccountID(t *testing.T) {
	codeHash := types.MustHexDecodeString("0x57f26a48169e57f118f6a0bca610a09b85cac7892f73c48f0c409cc8971817e7")
	data := types.MustHexDecodeString("0xd183512b0080c6a47e8d03000000000000000000")
	origin := utils.MustDecodeAccountIDFromSS58("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")

	id := rpc.GetContractAccountID(origin, types.CodeHash(types.NewHash(codeHash)), data, []byte{})

	// 5HCA3wMDkmk9MuNBjmed3rnkMwrzXa1rW4BK4JZ3JKh4ATz1
	idOk := utils.MustDecodeAccountIDFromSS58("5HCA3wMDkmk9MuNBjmed3rnkMwrzXa1rW4BK4JZ3JKh4ATz1")

	require.Equal(t, id, idOk)
}
