package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Query storage state from node",
}

func init() {
	initContractsCmds(storageCmd)
	rootCmd.AddCommand(storageCmd)
}

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

			meta, err := cli.API().RPC.State.GetMetadataLatest()
			if err != nil {
				return errors.Wrap(err, "get metadata failed")
			}

			key, err := types.CreateStorageKey(meta, "Contracts", "CurrentSchedule", nil, nil)
			if err != nil {
				return errors.Wrap(err, "create storage key failed")
			}

			var schedule types.Schedule
			ok, err := cli.API().RPC.State.GetStorageLatest(key, &schedule)
			if err != nil {
				return errors.Wrap(err, "get storage failed")
			}

			if !ok {
				return errors.Errorf("no found storages")
			}

			bz, _ := json.Marshal(schedule)
			fmt.Println(string(bz))

			return nil
		},
	})

	root.AddCommand(storageContractsCmd)
}
