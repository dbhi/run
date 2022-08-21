package main

import (
	"github.com/spf13/pflag"
	v "github.com/spf13/viper"

	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	au "github.com/logrusorgru/aurora"
)

func FlagFuncs(f *pflag.FlagSet) (flag func(k string, i interface{}, u string), flagP func(k, p string, i interface{}, u string)) {
	flag = func(k string, i interface{}, u string) {
		switch y := i.(type) {
		case bool:
			f.Bool(k, y, u)
		case int:
			f.Int(k, y, u)
		case string:
			f.String(k, y, u)
		}
		v.SetDefault(k, i)
	}
	flagP = func(k, p string, i interface{}, u string) {
		switch y := i.(type) {
		case bool:
			f.BoolP(k, p, y, u)
		case int:
			f.IntP(k, p, y, u)
		case string:
			f.StringP(k, p, y, u)
		}
		v.SetDefault(k, i)
	}
	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(au.Red(err))
		os.Exit(1)
	}
}

/*
func ExecTimedCmd(bin string, args []string, cmdOut, cmdErr *bytes.Buffer) {
	time_path, err := exec.LookPath("time")
	if err != nil {
		checkErr(fmt.Errorf("Please, install 'time': %s", err))
	}
	ExecCmd(time_path, append([]string{"-v", bin}, args...), cmdOut, cmdErr)
}
*/

func ExecCmd(bin string, args []string, cmdOut, cmdErr *bytes.Buffer) {
	fmt.Println(append([]string{bin}, args...))
	cmd := exec.Command(bin, args...)
	done := make(chan bool)
	/*                                          |-----> cmdOut.Write()
	     |-> .StdoutPipe() -> reader -> scan() -|-> |
	cmd -|                                          |-> fmt.Printf()
	     |-> .StderrPipe() -> reader -> scan() -|-> |
	                                            |-----> cmdErr.Write()
	*/
	for _, v := range []struct {
		rc  func() (io.ReadCloser, error)
		buf *bytes.Buffer
	}{
		{cmd.StdoutPipe, cmdOut},
		{cmd.StderrPipe, cmdErr},
	} {
		reader, err := v.rc()
		if err != nil {
			checkErr(fmt.Errorf("Error creating Std*Pipe for Cmd: %s", err))
		}
		go func(s *bufio.Scanner, b *bytes.Buffer) {
			for s.Scan() {
				fmt.Printf("%s\n", s.Text())
				if b != nil {
					b.Write([]byte(s.Text() + "\n"))
				}
			}
			done <- true
		}(bufio.NewScanner(reader), v.buf)
	}
	err := cmd.Run()
	if err != nil {
		checkErr(fmt.Errorf("Error running Cmd: %s", err))
	}
	for i := 0; i < 2; i++ {
		<-done
	}
}
