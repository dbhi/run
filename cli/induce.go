package main

import (
	"github.com/dbhi/run/lib"
	v "github.com/spf13/viper"
	"github.com/umarcor/cobra"
)

// induceCmd represents the induce command
var induceCmd = &cobra.Command{
	Use:   "induce",
	Short: "Induce subgraphs",
	Long:  `Induce subgraph for the given nodes.`,
	//Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		lib.Induce(v.GetString("graph"), v.GetString("output"), args)
	},
}

func init() {
	rootCmd.AddCommand(induceCmd)
	/*
		f := induceCmd.PersistentFlags()
		// Bind the full flag set to the configuration
		err := v.BindPFlags(f)
		checkErr(err)
	*/
}
