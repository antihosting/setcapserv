/*
 *
 * Copyright 2022-present Zander Schwid & Co. LLC. All rights reserved.
 *
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func startBackground(filename string) error {

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

	for _, command := range Commands {
		args = append(args, "-c", command)
	}

	args = append(args, filename)

	cmd := exec.Command(executable, args...)
	fmt.Printf("Run cmd: %v\n", formatCommand(executable, args))

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println("Daemon started in background. PID is : ", cmd.Process.Pid)
	return nil
}

func formatCommand(exe string, args []string) string {
	var out strings.Builder
	out.WriteString(exe)
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			out.WriteByte(' ')
			out.WriteByte('"')
			out.WriteString(arg)
			out.WriteByte('"')
		} else {
			out.WriteByte(' ')
			out.WriteString(arg)
		}
	}
	return out.String()
}