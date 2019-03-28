package main

import (
	"fmt"
	"log"
	"os"
	"path"

	v "github.com/spf13/viper"
	"github.com/umarcor/cobra"
	"github.com/umarcor/run/lib"
)

// induceCmd represents the induce command
var induceCmd = &cobra.Command{
	Use:   "induce",
	Short: "Induce subgraphs",
	Long:  `Induce subgraph for the given nodes.`,
	//	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("INDUCE")
		// TODO if 'graph' is empty, check if './graph.dot' exists
		s, err := lib.GetSubgraphs(v.GetString("graph"), false)
		checkErr(err)
		if len(s) == 0 {
			log.Fatal("Something went wrong. Empty subgraph map.")
		}
		o := v.GetString("output")
		if o != "" {
			checkErr(os.MkdirAll(o, 0644))
		}
		if len(args) != 0 {
			for _, a := range args {
				d, ok := s[a]
				if !ok {
					// TODO Check if it is a mid node
					fmt.Printf("subgraph for node '%s' not found\n", a)
					continue
				}
				if e := lib.WriteFile(path.Join(o, a+".dot"), d); e != nil {
					log.Fatal(err)
				}
			}
			return
		}
		for a, d := range s {
			if e := lib.WriteFile(path.Join(o, a+".dot"), d); e != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(induceCmd)

	f := induceCmd.PersistentFlags()

	// Bind the full flag set to the configuration
	err := v.BindPFlags(f)
	checkErr(err)
}
