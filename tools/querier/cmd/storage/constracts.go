package storagecmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initContractsCmds(root *cobra.Command) {
	storageContractsCmd := &cobra.Command{
		Use:   "contracts",
		Short: "Query storage state for contracts module",
	}

	storageContractsCmd.AddCommand(&cobra.Command{
		Use:   "schedule",
		Short: "Query current schedule for contracts module",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.NewLogger()

			cli, err := api.NewClient(logger, viper.GetString("url"))
			if err != nil {
				return errors.Wrap(err, "new client failed")
			}

			var schedule types.Schedule
			if err := cli.GetStorageLatest(&schedule, "Contracts", "CurrentSchedule", nil, nil); err != nil {
				return errors.Wrap(err, "get storage latest failed")
			}

			bz, _ := json.Marshal(schedule)
			fmt.Println(string(bz))

			return nil
		},
	})

	codeByHashCmd(storageContractsCmd)

	root.AddCommand(storageContractsCmd)
}

const (
	flagToFilePath = "file"
)

func codeByHashCmd(contractsCmd *cobra.Command) {
	res := &cobra.Command{
		Use:   "code",
		Short: "Query code by hash for contracts module",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.NewLogger()

			logger.Debug("args", "a", args)

			codeHash := types.MustHexDecodeString(args[0])

			cli, err := api.NewClient(logger, viper.GetString("url"))
			if err != nil {
				return errors.Wrap(err, "new client failed")
			}

			var codeBz []byte
			if err := cli.GetStorageLatest(&codeBz, "Contracts", "PristineCode", codeHash, nil); err != nil {
				return errors.Wrap(err, "get storage latest failed")
			}

			toFilePath := cmd.Flag(flagToFilePath).Value.String()
			if toFilePath == "" {
				bz, _ := json.Marshal(codeBz)
				fmt.Println(string(bz))
			} else if err := ioutil.WriteFile(toFilePath, codeBz, 0600); err != nil {
				return errors.Wrap(err, "write storage latest failed")
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	res.Flags().StringP(flagToFilePath, "f", "", "dump code to file")

	contractsCmd.AddCommand(res)
}
