package main

import (
	"io"
	"os"
	"strings"
)

func handleSingleCommand(cmdStr string, hist *[]string) {
	tokens := parseQuote(cmdStr)
	if len(tokens) == 0 {
		return
	}

	cmdStr = tokens[0]
	args := tokens[1:]

	// Handle redirects
	args, redirect, file := handleRedirects(args)

	// Restore stdout/stderr if redirected
	var stdout, stderr *os.File
	if file != nil {
		defer file.Close()

		switch redirect {
		case ">", "1>", ">>", "1>>":
			stdout = os.Stdout
			os.Stdout = file
		case "2>", "2>>":
			stderr = os.Stderr
			os.Stderr = file
		}
	}

	// Execute command
	ln := Line{cmd: cmdStr, args: args}
	cmd := getBuiltin(cmdStr)

	run(cmd, ln, os.Stdin, os.Stdout, os.Stderr, hist)

	if stdout != nil {
		os.Stdout = stdout
	}
	if stderr != nil {
		os.Stderr = stderr
	}
}

type parsedCmd struct {
	cmd     string
	args    []string
	builtin BUILTIN
}

func handlePipeline(commands []string, hist *[]string) {
	var cmds []parsedCmd
	for _, command := range commands {
		tokens := parseQuote(command)
		if len(tokens) == 0 {
			continue
		}

		pc := parsedCmd{
			cmd:     tokens[0],
			args:    tokens[1:],
			builtin: getBuiltin(tokens[0]),
		}

		cmds = append(cmds, pc)
	}

	if len(cmds) == 0 {
		return
	}

	pipeRead := make([]*io.PipeReader, len(cmds)-1)
	pipeWrite := make([]*io.PipeWriter, len(cmds)-1)
	for i := range len(cmds) - 1 {
		pipeRead[i], pipeWrite[i] = io.Pipe()
	}

	done := make(chan bool, len(cmds))

	for i, cmd := range cmds {
		var stdin io.Reader = os.Stdin
		var stdout io.Writer = os.Stdout
		stderr := os.Stderr

		// set input pipe
		if i > 0 {
			stdin = pipeRead[i-1]
		}

		if i < len(cmds)-1 {
			stdout = pipeWrite[i]
		} else {
			var redirect string
			var file *os.File

			cmd.args, redirect, file = handleRedirects(cmd.args)

			if file != nil {
				defer file.Close()

				switch redirect {
				case ">", ">>", "1>", "1>>":
					stdout = file
				case "2>", "2>>":
					stderr = file
				}
			}
		}

		ln := Line{
			cmd:  cmd.cmd,
			args: cmd.args,
		}

		go func(
			i int,
			builtin BUILTIN,
			ln Line,
			stdin io.Reader,
			stdout,
			stderr io.Writer,
		) {

			defer func() {
				if i < len(pipeWrite) {
					pipeWrite[i].Close()
				}

				done <- true
			}()

			run(builtin, ln, stdin, stdout, stderr, hist)

		}(i, cmd.builtin, ln, stdin, stdout, stderr)
	}

	for range len(cmds) {
		<-done
	}
}

func splitByPipes(line string) []string {
	var cmds []string
	var curr strings.Builder
	quote := false
	dquote := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		switch char {
		case '\'':
			if !dquote {
				quote = !quote
			}
			curr.WriteByte(char)
		case '"':
			if !quote {
				dquote = !dquote
			}
			curr.WriteByte(char)
		case '|':
			if !quote && !dquote {
				cmds = append(cmds, strings.TrimSpace(curr.String()))
				curr.Reset()
			} else {
				curr.WriteByte(char)
			}
		default:
			curr.WriteByte(char)
		}
	}

	if curr.Len() > 0 {
		cmds = append(cmds, strings.TrimSpace(curr.String()))
	}

	return cmds
}

func handleRedirects(args []string) ([]string, string, *os.File) {
	for i, arg := range args {
		if isRedirect(arg) && i+1 < len(args) {
			redirect := arg
			filename := args[i+1]

			var flag int
			switch redirect {
			case ">", "1>", "2>":
				flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			case ">>", "1>>", "2>>":
				flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
			default:
				return args[:i], "", nil
			}

			file, err := os.OpenFile(filename, flag, 0644)
			if err != nil {
				return args[:i], "", nil
			}

			return args[:i], redirect, file
		}
	}
	return args, "", nil
}
