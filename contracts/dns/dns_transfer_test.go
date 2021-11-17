package dns_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/contracts/dns"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/stretchr/testify/require"
)

func TestDNSTransfer(t *testing.T) {
	test.ByNodeEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		contractAccountID := initDNS(t, logger, env, signature.TestKeyringPairAlice)
		fmt.Println("============================ success create contract")
		fmt.Println(contractAccountID)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(dnsMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		dnsAPI := dns.New(rpcAPI, contractAccountID)

		aliceCtx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)
		fmt.Println("register ============================================")
		_, err = dnsAPI.Register(aliceCtx, initName)
		fmt.Println(err)
		require.Nil(err)

		fmt.Println("transfer ============================================")
		dnsAPI.Transfer(aliceCtx, initName, bob)
		fmt.Println("setaddress ============================================")

		_, err = dnsAPI.SetAddress(aliceCtx, initName, charlie)
		require.Nil(err)

		resName, err := dnsAPI.GetAddress(aliceCtx, initName)
		require.Nil(err)
		require.Equalf(resName, nil, "Alice's authority has been transferred to Bob.")

		// Switch to Bob's identity to send the transaction
		bobCtx := rpc.NewCtx(context.Background()).WithFrom(TestKeyringPairBob)
		_, err = dnsAPI.SetAddress(bobCtx, initName, charlie)
		require.Nil(err)

		resName, err = dnsAPI.GetAddress(bobCtx, initName)
		require.Nil(err)
		require.Equalf(resName, charlie, "Bob successfully set the address.")
	})
}
