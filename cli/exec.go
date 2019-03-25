package main

import (
	"fmt"

	v "github.com/spf13/viper"
	"github.com/umarcor/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec list of tasks",
	Long:  `Exec list of tasks for the given nodes.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range args {
			fmt.Println("EXEC", a)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	f := execCmd.PersistentFlags()

	// Bind the full flag set to the configuration
	err := v.BindPFlags(f)
	checkErr(err)
}
