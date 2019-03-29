package main

import (
	"github.com/umarcor/cobra"
	"github.com/umarcor/run/lib"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec list of tasks",
	Long:  `Exec list of tasks for the given nodes.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lib.Exec(args)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	/*
		f := execCmd.PersistentFlags()
		// Bind the full flag set to the configuration
		err := v.BindPFlags(f)
		checkErr(err)
	*/
}
