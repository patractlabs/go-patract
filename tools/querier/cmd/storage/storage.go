package storagecmd

import (
	"github.com/spf13/cobra"
)

// AppendCmds append the cmd from package
func AppendCmds(root *cobra.Command) {
	var storageCmd = &cobra.Command{
		Use:   "storage",
		Short: "Query storage state from node",
	}

	// sub cmds by contracts
	initContractsCmds(storageCmd)

	root.AddCommand(storageCmd)
}
