package dns_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/test/contracts"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

const (
	dnsWasmPath = "../../test/contracts/ink/dns.wasm"
	dnsMetaPath = "../../test/contracts/ink/dns.json"
)

var (
	bob     = utils.MustAccountIDFromSS58("5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty")
	charlie = utils.MustAccountIDFromSS58("5FLSigC9HGRKVhB9FiEo4Y3koPsNmBmLJbpXg2mp1hXcS59Y")

	TestKeyringPairBob, _ = signature.KeyringPairFromSecret("//Bob", 42)
	instantiateSalt       = []byte("ysncz3nbjjzoc7s07of3malp9d")

	initName = types.NewHash(types.MustHexDecodeString("0x43a3d97a2bf7293835cdecdd25df56bb18c9d9c0e856246622a6c5164c2f2e0e"))
)

func initDNS(t *testing.T, logger log.Logger, env test.Env, authKey signature.KeyringPair) types.AccountID {
	require := require.New(t)

	codeBytes, err := ioutil.ReadFile(dnsWasmPath)
	require.Nil(err)

	cApi, err := rpc.NewContractAPI(env.URL())
	require.Nil(err)

	metaBz, err := ioutil.ReadFile(dnsMetaPath)
	fmt.Println("--------------------------------")
	//fmt.Println(metaBz)
	require.Nil(err)
	cApi.WithMetaData(metaBz)

	ctx := api.NewCtx(context.Background()).WithFrom(authKey)

	var endowment uint64 = 1000000000000

	// Instantiate
	_, contractAccount, err := cApi.InstantiateWithCode(ctx, logger,
		types.NewCompactBalance(endowment),
		types.NewCompactGas(test.DefaultGas),
		contracts.CodeHashDNS,
		codeBytes,
		instantiateSalt,
		[]string{"new"},
		//nil,
	)
	require.Nil(err)

	// check code
	var codeBz []byte
	if err := cApi.Native().Cli.GetStorageLatest(&codeBz,
		"Contracts", "PristineCode",
		contracts.CodeHashDNS[:], nil); err != nil {
		require.Nil(err)
	}

	t.Logf("constract %s", types.HexEncodeToString(contractAccount[:]))
	return contractAccount
}
