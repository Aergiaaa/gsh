package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type BUILTIN int

const (
	ECHO BUILTIN = iota
	EXIT
	TYPE
	HISTORY
	PWD
	CD
)

func getBuiltin(cmd string) BUILTIN {
	switch cmd {
	case "echo":
		return ECHO
	case "exit":
		return EXIT
	case "type":
		return TYPE
	case "history":
		return HISTORY
	case "pwd":
		return PWD
	case "cd":
		return CD
	default:
		return -1
	}
}

func pwd(r io.Reader, w io.Writer) {
	defer drainStdin(r)()

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%s\n", dir)
}

func cd(args []string, r io.Reader, w io.Writer) {
	defer drainStdin(r)()

	var err error
	dir := args[0]
	if dir == "~" {
		dir, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(w, "%v\n", err)
		}
	}

	err = os.Chdir(dir)
	if err != nil {
		fmt.Fprintf(w, "cd: no such file or directory: %s\n", args[0])
	}
}

func _type(cmds []string, r io.Reader, w io.Writer) {
	defer drainStdin(r)()

	for _, cmd := range cmds {
		builtin := getBuiltin(cmd)
		if builtin != -1 {
			fmt.Fprintf(w, "%s is a shell builtin\n", cmd)
			return
		}

		if path, err := exec.LookPath(cmd); err == nil {
			fmt.Fprintf(w, "%s is %s\n", cmd, path)
			return
		}

		fmt.Fprintln(w, cmd+": not found")
	}
}

func echo(args []string, r io.Reader, w io.Writer) {
	defer drainStdin(r)()
	fmt.Fprintln(w, strings.Join(args, " "))
}

func exit(r io.Reader) {
	defer drainStdin(r)()
	fmt.Println("exit")
	os.Exit(0)
}
