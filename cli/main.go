package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/umarcor/run/lib"

	au "github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
	v "github.com/spf13/viper"
	"github.com/umarcor/cobra"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	Execute()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "run",
	Version: "v0.0.0",
	Short:   au.Sprintf(au.Cyan("[RUN] a task execution automation package")),
	Long: `A task execution automation package for complex dependency graphs.
Currently only DOT files are supported as input. To retrieve the ordered list of
tasks for a given target, use the following syntax 'leaf[|task]'. The optional
argument 'task' allows to filter the list to include only a subset of the tasks
in the subgraphs corresponding to the leaf. It can be either of:
- '>DOTID': tasks that allow build DOTID.
- 'DOTID>': tasks that depend on DOTID.
- '>DOTID>': tasks that allow to build DOTID and those that depend on it.
`,
	Args: cobra.ArbitraryArgs,
	Run: func(_ *cobra.Command, args []string) {
		s, err := lib.GetSubgraphs(v.GetString("file"), false)
		CheckErr(err)
		t := lib.GetTaskListAll(s)
		fmt.Println(t)
		//g := dot.Load()
		//fmt.Printf("%+v\n", g)
	},
}

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Set custom version template
	rootCmd.SetVersionTemplate("RUN {{printf \"version %s\" .Version}}\n")
	fmt.Println(au.Sprintf(au.Cyan("[RUN] a task execution automation package (%s)"), rootCmd.Version))
	err := rootCmd.Execute()
	CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	f := rootCmd.PersistentFlags()
	// Helper functions to set cobra and viper at once
	_, flagP := FlagFuncs(f)

	// Define flags and defaults
	f.StringVarP(&cfgFile, "config", "c", "", "config file (defaults are './.run[ext]', '$HOME/.run[ext]' or '/etc/run/.run[ext]')")
	flagP("log", "l", "stdout", "errors logger; can use 'stdout', 'stderr' or file")
	flagP("file", "f", "", "input DOT file")
	//flagP("output", "o", "", "output ('stdout' or file)")

	// Bind the full flag set to the configuration
	err := v.BindPFlags(f)
	CheckErr(err)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		CheckErr(err)

		v.AddConfigPath(".")
		v.AddConfigPath(home)
		v.AddConfigPath("/etc/run/")
		v.SetConfigName(".run")
	}

	v.SetEnvPrefix("RUN")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		// Fail with invalid config format
		if _, ok := err.(v.ConfigParseError); ok {
			CheckErr(err)
		}
	} else {
		log.Println("Using config file:", v.ConfigFileUsed())
	}

	switch l := v.GetString("log"); l {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "":
		log.SetOutput(ioutil.Discard)
	default:
		log.SetOutput(&lumberjack.Logger{
			Filename:   l,
			MaxSize:    100,
			MaxAge:     14,
			MaxBackups: 10,
		})
	}

	if !v.IsSet("indocker") {
		v.Set("indocker", false)
	}
	cmd := exec.Command("cat", "/proc/self/cgroup")
	o, err := cmd.CombinedOutput()
	CheckErr(err)
	if strings.Contains(string(o), "docker") {
		log.Println("It seems you are running RUN CLI inside Docker a container")
		v.Set("indocker", true)
	}
}
