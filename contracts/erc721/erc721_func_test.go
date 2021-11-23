package erc721_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/contracts/erc721"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestERC721(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)

		contractAccountID := initERC721(t, logger, env, signature.TestKeyringPairAlice)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(erc721MetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		erc721API := erc721.New(rpcAPI, contractAccountID)

		ctx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)

		// mint tokenId by alice
		_, err = erc721API.Mint(ctx, tokenId)
		require.Nil(err)

		aliceTotal, err := erc721API.BalanceOf(ctx, test.AliceAccountID)
		require.Nil(err)
		require.Equalf(aliceTotal, types.NewU32(1), "alice should be 1")

		resOwn, err := erc721API.OwnerOf(ctx, tokenId)
		require.Nil(err)
		require.Equalf(resOwn.Value, test.AliceAccountID, "owner should be alice")

		// transfer alice to bob
		_, err = erc721API.Transfer(ctx, bob, tokenId)
		require.Nil(err)

		bobBalance, err := erc721API.BalanceOf(ctx, bob)
		require.Nil(err)
		require.Equalf(bobBalance, types.NewU32(1), "bob Balance should be 1")

		aliceNewTotal, err := erc721API.BalanceOf(ctx, test.AliceAccountID)
		require.Nil(err)
		require.Equalf(aliceNewTotal, types.NewU32(0), "alice new add transfer should be alice none")

		resOwn, err = erc721API.OwnerOf(ctx, tokenId)
		require.Nil(err)
		require.Equalf(resOwn.Value, bob, "owner of transfer should be bob")

		bobCtx := rpc.NewCtx(context.Background()).WithFrom(TestKeyringPairBob)

		_, err = erc721API.Approve(bobCtx, charlie, tokenId)
		require.Nil(err)

		resApprove, err := erc721API.GetApproved(bobCtx, tokenId)
		require.Nil(err)
		require.Equalf(resApprove.Value, charlie, "approved for charlie")

		_, err = erc721API.SetApprovalForAll(bobCtx, test.AliceAccountID, types.NewBool(true))
		require.Nil(err)

		resSetApproval, err := erc721API.IsApprovedForAll(bobCtx, bob, test.AliceAccountID)
		require.Nil(err)
		require.Equalf(resSetApproval, types.NewBool(true), "approval is must be true")

		_, err = erc721API.Burn(bobCtx, tokenId)
		require.Nil(err)

		bobNewTotal, err := erc721API.BalanceOf(ctx, bob)
		require.Nil(err)
		require.Equalf(bobNewTotal, types.NewU32(0), "The tokenId was burned")
	})
}
