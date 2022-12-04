/*
 *
 * Copyright 2022-present Zander Schwid & Co. LLC. All rights reserved.
 *
 */

package main

import (
	"flag"
	"time"
)

type ArrayListFlag []string

func (t *ArrayListFlag) String() string {
	return "Repeatable ordered sequence of commands"
}

func (t *ArrayListFlag) Set(s string) error {
	*t = append(*t, s)
	return nil
}

var (
	Commands    ArrayListFlag

	Verbose    = flag.Bool("v", false, "Print logs and debug information")
	Foreground = flag.Bool("f", false, "Indicator that daemon is running in foreground")
	LogFile    = flag.String("log", "stdout", "Write log to file, stdout, stderr")

	Delay      = flag.Duration("d", 3*time.Second, "Delay on update after last event")
)

func init() {
	flag.CommandLine.Var(&Commands, "c", "Run sequence of commands")
}

