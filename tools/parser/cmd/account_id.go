package cmd

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/spf13/cobra"
)

func logAccountID(accountID types.AccountID) {
	ss58, _ := utils.EncodeAccountIDToSS58(accountID)

	fmt.Printf("ss58\t:\t%s\n", ss58)
	fmt.Printf("hex\t:\t%s\n", types.HexEncodeToString(accountID[:]))
}

func parseAccountIDCmd(rootCmd *cobra.Command) {
	res := &cobra.Command{
		Use:   "id",
		Short: "Parse AccountID to ss58 or hex",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.NewLogger()

			logger.Debug("args", "a", args)

			// if is SS58
			ss58AccountID, err := utils.NewAccountIDFromSS58(args[0])
			if err == nil {
				logAccountID(ss58AccountID)
			}

			// if is byte hex
			hexBytes, err := types.HexDecodeString(args[0])
			if err == nil && len(hexBytes) == 32 {
				logAccountID(types.NewAccountID(hexBytes))
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	rootCmd.AddCommand(res)
}
