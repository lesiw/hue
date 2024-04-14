package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"lesiw.io/flag"
)

var flags = flag.NewSet(os.Stderr, "hue COMMAND")
var errParse = errors.New("parse error")
var outprefix, errprefix string

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, errParse) {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}

func run() (err error) {
	if outprefix = os.Getenv("HUEOUT"); outprefix == "" {
		outprefix = "\033[39m"
	} else if outprefix, err = colorToAnsi(outprefix); err != nil {
		return err
	}
	if errprefix = os.Getenv("HUEERR"); errprefix == "" {
		errprefix = "\033[31m"
	} else if errprefix, err = colorToAnsi(errprefix); err != nil {
		return err
	}
	if len(os.Args) < 2 {
		flags.PrintError("no command given")
		return errParse
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdout := make(chanWriter)
	stderr := make(chanWriter)
	done := make(chan bool)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	go func() {
		var ee *exec.ExitError
		if err := cmd.Run(); err != nil && !errors.As(err, &ee) {
			stderr <- err.Error() + "\n"
		}
		done <- true
	}()
	var last chanWriter
	for {
		select {
		case m := <-stdout:
			if last != stdout {
				fmt.Fprint(os.Stdout, "\033[0m", outprefix)
			}
			fmt.Fprint(os.Stdout, m)
			last = stdout
		case m := <-stderr:
			if last != stderr {
				fmt.Fprint(os.Stderr, "\033[0m", errprefix)
			}
			fmt.Fprint(os.Stderr, m)
			last = stderr
		case <-done:
			return nil
		}
	}
}

type chanWriter chan string

func (cw chanWriter) Write(b []byte) (int, error) {
	cw <- string(b)
	return len(b), nil
}
