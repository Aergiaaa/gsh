package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	if env, exist := os.LookupEnv("HISTFILE"); exist {
		HISTFILE = env
	}
	history := loadHist(HISTFILE)

	defer saveHist(history, HISTFILE)

	buffer, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		InterruptPrompt: "vvv",
		AutoComplete: &builtinCompleter{},
	})
	if err != nil {
		fmt.Printf("error reading line: %v\n", err)
	}
	defer buffer.Close()

	for {
		line, err := buffer.Readline()
		if err != nil {
			continue
		}

		if strings.TrimSpace(line) != "" {
			history = append(history, line)
			appendHist(line, HISTFILE)
		}

		cmds := splitByPipes(line)

		if len(cmds) == 0 {
			continue
		}

		if len(cmds) == 1 {
			handleSingleCommand(cmds[0], &history)
			continue
		}

		handlePipeline(cmds, &history)
	}
}
