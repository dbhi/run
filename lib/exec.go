package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

/*
func ExecTimedCmd(dir, bin string, args, env []string, cmdOut, cmdErr *bytes.Buffer, verbose bool) error {
	time_path, err := exec.LookPath("time")
	if err != nil {
		checkErr(fmt.Errorf("Please, install 'time': %s", err))
	}
	return ExecCmd(dir, time_path, append([]string{"-v", bin}, args...), env, cmdOut, cmdErr, verbose)
}
*/

func ExecCmd(dir, bin string, args, env []string, cmdOut, cmdErr *bytes.Buffer, verbose bool) error {
	cmd := exec.Command(bin, args...)

	cmd.Env = os.Environ()
	for _, v := range env {
		cmd.Env = append(cmd.Env, v)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	var teeOut, teeErr io.Reader
	for _, v := range []struct {
		rc  func() (io.ReadCloser, error)
		tee *io.Reader
		buf *bytes.Buffer
	}{
		{cmd.StdoutPipe, &teeOut, cmdOut},
		{cmd.StderrPipe, &teeErr, cmdErr},
	} {
		r, err := v.rc()
		if err != nil {
			return err
		}
		*v.tee = io.Reader(r)
		if v.buf != nil {
			*v.tee = io.TeeReader(r, v.buf)
		}
	}

	done := make(chan bool)
	go func() {
		s := bufio.NewScanner(io.MultiReader(teeOut, teeErr))
		for s.Scan() {
			if verbose {
				fmt.Printf("%s\n", s.Text())
			}
		}
		done <- true
	}()

	err := cmd.Run()
	if err != nil {
		return err
	}

	<-done
	return nil
}
