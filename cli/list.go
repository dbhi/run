package main

import (
	v "github.com/spf13/viper"
	"github.com/umarcor/cobra"
	"github.com/umarcor/run/lib"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the topological order",
	Long:  `List the topological order of the subgraph(s) for the given node(s).`,
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lib.List(v.GetString("graph"), args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	/*
		f := listCmd.PersistentFlags()
		// Bind the full flag set to the configuration
		err := v.BindPFlags(f)
		checkErr(err)
	*/
}
