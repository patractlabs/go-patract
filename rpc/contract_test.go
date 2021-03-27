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
	origin := utils.MustDecodeAccountIDFromSS58("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")

	id := rpc.GetContractAccountID(origin, types.NewHash(codeHash), instantiateSalt)

	idOk := utils.MustDecodeAccountIDFromSS58(contractAddress)

	require.Equal(t, id, idOk)
}
