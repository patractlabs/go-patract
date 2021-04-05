package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "parser",
		Short: "A tool to parse AccountID and datas",
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("ss58Prefix", "s", "", "ss58Prefix")
	if err := viper.BindPFlag("ss58Prefix", rootCmd.PersistentFlags().Lookup("ss58Prefix")); err != nil {
		panic(err)
	}

	parseAccountIDCmd(rootCmd)
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
