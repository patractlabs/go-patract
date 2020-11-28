package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	metadataCmd.PersistentFlags().StringP("module", "m", "", "module data to query")
	if err := viper.BindPFlag("module", metadataCmd.PersistentFlags().Lookup("module")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(metadataCmd)
}

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Query metadata for runtime from node",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.NewLogger()

		cli, err := api.NewClient(logger, viper.GetString("url"))
		if err != nil {
			return errors.Wrap(err, "new client failed")
		}

		data, err := cli.API().RPC.State.GetMetadataLatest()
		if err != nil {
			return errors.Wrap(err, "get metadata failed")
		}

		var bz []byte
		moduleToQuery := viper.GetString("module")

		if moduleToQuery == "" {
			// now will be V12
			bz, err = json.Marshal(data.AsMetadataV12)
			if err != nil {
				return errors.Wrap(err, "marshal data failed")
			}
		} else {
			for _, m := range data.AsMetadataV12.Modules {
				if strings.EqualFold(string(m.Name), moduleToQuery) {
					bz, err = json.Marshal(m)
					if err != nil {
						return errors.Wrap(err, "marshal module data failed")
					}
					break
				}
			}

			if bz == nil {
				// no found modules
				return errors.Errorf("no found modules %s", moduleToQuery)
			}
		}

		_, err = fmt.Println(string(bz))

		return err
	},
}
