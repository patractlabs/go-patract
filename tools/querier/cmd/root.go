package cmd

import (
	"fmt"
	"os"

	storagecmd "github.com/patractlabs/go-patract/tools/querier/cmd/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "querier",
		Short: "A tool to query datas from node",
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("url", "u", "ws://localhost:9944", "RPC url for node")
	if err := viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url")); err != nil {
		panic(err)
	}

	storagecmd.AppendCmds(rootCmd)
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
