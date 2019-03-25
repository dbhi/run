package main

import (
	"fmt"

	v "github.com/spf13/viper"
	"github.com/umarcor/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the topological order",
	Long:  `List the topological order of the subgraph(s) for the given node(s).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range args {
			fmt.Println("LIST", a)

			/*
				err := ioutil.WriteFile("testdata/hello", message, 0644)
				if err != nil {
					log.Fatal(err)
				}*/
			/*
				d, ok := s["happ"]
					if ok {
						l := lib.GetTaskList(d, "ghdl -a [UUT]")
						fmt.Println(l)
					}
					//t := lib.GetTaskListAll(s)
			*/
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	f := listCmd.PersistentFlags()

	// Bind the full flag set to the configuration
	err := v.BindPFlags(f)
	checkErr(err)
}
