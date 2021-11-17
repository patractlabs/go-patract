package dns_test

import (
	"context"
	"fmt"
	"github.com/patractlabs/go-patract/types"
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
	test.ByExternEnv(t, func(logger log.Logger, env test.Env) {
		require := require.New(t)
		//contractAccountID := initDNS(t, logger, env, signature.TestKeyringPairAlice)
		//fmt.Println("============================ success create contract")
		//fmt.Println(contractAccountID)
		//contractAccountID := initERC20(t, logger, env, signature.TestKeyringPairAlice)
		//b := "5CcZdeQEH7Q6qy1PqE6uaPns5u1rJtjdi6yHjmoMX3gEpkMW" // 合约地址
		//contractAccountID := initERC20(t, logger, env, signature.TestKeyringPairAlice)
		ss58byte := types.NewSS58Codec([]byte{})
		b := "5Ce5P6WV7Z4p75DGPgmnsFm73HKvFkifivp7cEJirsgsat2q" // 合约地址
		contractAccountID, _ := ss58byte.DecodeAccountID(b)
		fmt.Println(contractAccountID)
		fmt.Println(contractAccountID)
		//strcode, _ := ss58byte.EncodeAccountID(contractAccountID)
		fmt.Println("====================================================")
		fmt.Println("====================================================")
		fmt.Println(contractAccountID)
		rpcAPI, err := rpc.NewContractAPI(env.URL())
		require.Nil(err)

		metaBz, err := ioutil.ReadFile(dnsMetaPath)
		require.Nil(err)
		rpcAPI.WithMetaData(metaBz)

		rpcAPI.WithLogger(logger)

		dnsAPI := dns.New(rpcAPI, contractAccountID)

		aliceCtx := rpc.NewCtx(context.Background()).WithFrom(signature.TestKeyringPairAlice)
		//fmt.Println("register ============================================")
		//_, err = dnsAPI.Register(aliceCtx, initName)
		//fmt.Println(err)
		//require.Nil(err)

		//fmt.Println("transfer ============================================")
		//dnsAPI.Transfer(aliceCtx, initName, bob)
		//fmt.Println("setaddress ============================================")

		//_, err = dnsAPI.SetAddress(aliceCtx, initName, charlie)
		//require.Nil(err)

		resName, err := dnsAPI.GetAddress(aliceCtx, initName)
		fmt.Println("==========================================")
		fmt.Println("==========================================")
		fmt.Println("==========================================")
		fmt.Println(resName)
		fmt.Println(err)
		require.Nil(err)
		require.Equalf(resName, nil, "Alice's authority has been transferred to Bob.")

		//// Switch to Bob's identity to send the transaction
		//bobCtx := rpc.NewCtx(context.Background()).WithFrom(TestKeyringPairBob)
		//res, err := dnsAPI.SetAddress(bobCtx, initName, charlie)
		//fmt.Println("=============================================")
		//fmt.Println("=============================================")
		//fmt.Println(res)
		//fmt.Println(err)
		//require.Nil(err)
		//
		//resName, err = dnsAPI.GetAddress(bobCtx, initName)
		//fmt.Println("-____________________________________________")
		//fmt.Println("-____________________________________________")
		//fmt.Println("-____________________________________________")
		//fmt.Println(resName)
		//require.Nil(err)
		//require.Equalf(resName, charlie, "Bob successfully set the address.")
	})
}
