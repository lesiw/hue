package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"

	"github.com/netflix/go-iomux"
	"lesiw.io/flag"
)

var flags = flag.NewSet(os.Stderr, "hue COMMAND")
var errParse = errors.New("parse error")
var outprefix, errprefix string
var stdout, stderr *os.File
var defers deferlist

type msg int

const (
	msgout msg = iota
	msgerr
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		defers.run()
		os.Exit(1)
	}()
	if err := run(); err != nil {
		if !errors.Is(err, errParse) {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}

func run() (err error) {
	defer defers.run()
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
	cmd.Stdin = os.Stdin
	if os.Getenv("HUEASYNC") == "1" {
		return runAsync(cmd)
	}
	return runSync(cmd)
}

func runSync(cmd *exec.Cmd) (err error) {
	mux := &iomux.Mux[msg]{}
	defers.add(func() { _ = mux.Close() })
	if stdout, err = mux.Tag(msgout); err != nil {
		return fmt.Errorf("failed creating mux tag: %s", err)
	}
	if stderr, err = mux.Tag(msgerr); err != nil {
		return fmt.Errorf("failed creating mux tag: %s", err)
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	ctx, cancel := context.WithCancel(context.Background())
	defers.add(cancel)
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("could not start command: %s", err)
	}
	go func() { _ = cmd.Wait(); cancel() }()
	var last msg
	var buf []byte
	var m msg
	for err != io.EOF {
		buf, m, err = mux.Read(ctx)
		if err != nil && err != io.EOF {
			return err
		}
		switch m {
		case msgout:
			if last != msgout {
				fmt.Print("\033[0m", outprefix)
			}
			fmt.Print(string(buf))
			last = msgout
		case msgerr:
			if last != msgerr {
				fmt.Print("\033[0m", errprefix)
			}
			fmt.Print(string(buf))
			last = msgerr
		}
	}
	return nil
}

func runAsync(cmd *exec.Cmd) error {
	stdout := make(chanWriter)
	stderr := make(chanWriter)
	done := make(chan bool, 1)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start command: %s", err)
	}
	go func() { _ = cmd.Wait(); done <- true }()
	defers.add(func() { done <- true })
	var last chanWriter
	for {
		select {
		case m := <-stdout:
			if last != stdout {
				fmt.Print("\033[0m", outprefix)
			}
			fmt.Print(m)
			last = stdout
		case m := <-stderr:
			if last != stderr {
				fmt.Print("\033[0m", errprefix)
			}
			fmt.Print(m)
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
