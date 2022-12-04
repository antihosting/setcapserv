/*
 *
 * Copyright 2022-present Zander Schwid & Co. LLC. All rights reserved.
 *
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	rt "runtime"
	"github.com/pkg/errors"
	"strings"
)

var (
	Version string
	Build   string
)

func main() {
	rt.GOMAXPROCS(1)
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	if err := doRun(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	} else {
		return 0
	}

}

func doRun(args []string) error {

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	args = flag.CommandLine.Args()

	if len(args) == 0 {
		flag.PrintDefaults()
		return errors.New("Usage: sudo ./autoupdate [flags] file_path")
	}

	watchFilePath := args[0]

	if !*Foreground {
		// fork the process to run in background
		return startBackground(*Verbose)
	}

	var logFile *os.File
	var logWriter io.Writer

	if *LogFile == "stdout" {
		logWriter = os.Stdout
	} else if *LogFile == "stderr" {
		logWriter = os.Stderr
	} else {
		var err error
		logFile, err = os.OpenFile(*LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return errors.Errorf("fail to open file '%s', %v", *LogFile, err)
		}
		logWriter = logFile
	}
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	log := log.New(logWriter,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	log.Printf("Setcapserv %s %s\n", Version, Build)
	log.Printf("Watch File: %s\n", watchFilePath)
	log.Printf("Verbose: %v\n", *Verbose)

	commands := splitAndTrim(Commands)
	log.Printf("Commands: %v\n", commands)

	if containsFirstString(commands, "setcap") {
		root, err := isRoot()
		if err != nil {
			return err
		}

		if !root {
			log.Printf("Warning: autoupdate with option -setcap should run under root user")
		}
	}

	daemon := NewDeamon(watchFilePath, commands, log, *Verbose)
	return daemon.Run(NewApp())
}

func isRoot() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, err
	}
	return currentUser.Username == "root", nil
}


func startBackground(verbose bool) error {

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	args := []string{
		"-f",
		"-log", executable + ".log",
	}

	if *Verbose {
		args = append(args, "-v")
	}

	cmd := exec.Command(executable, args...)
	fmt.Printf("Run cmd: %v\n", cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println("Daemon started in background. PID is : ", cmd.Process.Pid)
	return nil
}

func containsFirstString(commands [][]string, value string) bool {
	for _, args := range commands {
		if len(args) > 0 && args[0] == value {
			return true
		}
	}
	return false
}

func splitAndTrim(arr []string) [][]string {
	var out [][]string
	for _, s := range arr {
		var args []string
		for _, arg := range strings.Split(s, " ") {
			a := strings.TrimSpace(arg)
			if len(a) > 0 {
				args = append(args, a)
			}
		}
		out = append(out, args)
	}
	return out
}

