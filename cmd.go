package main

import (
	"io"
	"os/exec"
)

func run(
	cmd BUILTIN, ln Line,
	stdin io.Reader, stdout io.Writer, stderr io.Writer,
	hist *[]string,
) {
	switch cmd {
	case ECHO:
		echo(ln.args, stdin, stdout)
	case EXIT:
		exit(stdin)
	case TYPE:
		_type(ln.args, stdin, stdout)
	case HISTORY:
		history(ln.args, stdin, stdout, hist)
	case PWD:
		pwd(stdin, stdout)
	case CD:
		cd(ln.args, stdin, stdout)
	default:
		runBin(ln.cmd, ln.args, stdin, stdout, stderr)
	}
}

func runBin(cmd string, args []string, stdin io.Reader, stdout, stderr io.Writer) {
	if _, err := exec.LookPath(cmd); err != nil {
		notExist(cmd)
		return
	}

	extCmd := exec.Command(cmd, args...)
	extCmd.Stdin = stdin
	extCmd.Stderr = stderr
	extCmd.Stdout = stdout

	_ = extCmd.Run()
}
