package cmd

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func logAccountID(codec *types.SS58Codec, accountID types.AccountID) {
	ss58, _ := codec.EncodeAccountID(accountID)

	fmt.Printf("ss58\t:\t%s\n", ss58)
	fmt.Printf("hex\t:\t%s\n", types.HexEncodeToString(accountID[:]))
}

func parseAccountIDCmd(rootCmd *cobra.Command) {
	res := &cobra.Command{
		Use:   "id",
		Short: "Parse AccountID to ss58 or hex",
		RunE: func(cmd *cobra.Command, args []string) error {
			ss58Codec := types.NewSS58Codec([]byte(viper.GetString("ss58Prefix")))

			// if is SS58
			ss58AccountID, err := ss58Codec.DecodeAccountID(args[0])
			if err == nil {
				logAccountID(ss58Codec, ss58AccountID)
			}

			// if is byte hex
			hexBytes, err := types.HexDecodeString(args[0])
			if err == nil && len(hexBytes) == 32 {
				logAccountID(ss58Codec, types.NewAccountID(hexBytes))
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	rootCmd.AddCommand(res)
}
