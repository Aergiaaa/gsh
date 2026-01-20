package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Line struct {
	cmd  string
	args []string
}

type builtinCompleter struct {
	lastPref  string
	lastMatch []string
	tabCount  int
}

var builtins = []string{"echo", "exit", "type", "history", "pwd", "cd"}

func (bc *builtinCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	text := string(line[:pos])

	parts := strings.Fields(text)
	if len(parts) == 0 {
		return nil, 0
	}

	prefix := parts[len(parts)-1]

	// Only autocomplete if we're at the first word (command position)
	// Check if there's a space after the last field - if so, don't complete
	if len(parts) > 1 || (len(text) > 0 && text[len(text)-1] == ' ') {
		bc.tabCount = 0
		return nil, 0
	}

	// List of builtin commands
	execPath := getExecPath()

	commands := append(builtins, execPath...)
	// Find matching commands
	var matches []string
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, cmd)
		}
	}

	sort.Strings(matches)

	if prefix != bc.lastPref {
		bc.tabCount = 0
		bc.lastPref = prefix
		bc.lastMatch = matches
	}

	bc.tabCount++

	if len(matches) == 0 {
		fmt.Print("\a")
		bc.tabCount = 0
		return nil, len(prefix)
	}

	if len(matches) == 1 {
		suffix := matches[0][len(prefix):] + " "
		bc.tabCount = 0
		return [][]rune{[]rune(suffix)}, len(prefix)
	}

	lcp := longestCommonPrefix(matches)
	if len(lcp) > len(prefix) {
		suffix := lcp[len(prefix):]
		bc.tabCount = 0
		bc.lastPref = lcp
		return [][]rune{[]rune(suffix)}, len(prefix)
	}

	if bc.tabCount == 1 {
		fmt.Print("\a")
		return nil, len(prefix)
	}

	fmt.Println()
	fmt.Print(strings.Join(matches, "  "))
	fmt.Println()
	bc.tabCount = 0
	return [][]rune{[]rune("")}, len(prefix)
}

func longestCommonPrefix(matches []string) string {
	if len(matches) == 0 {
		return ""
	}
	if len(matches) == 1 {
		return matches[0]
	}

	prefix := matches[0]
	for _, match := range matches[1:] {
		i := 0
		for i < len(prefix) && i < len(match) && prefix[i] == match[i] {
			i++
		}

		prefix = prefix[:i]
		if prefix == "" {
			return ""
		}
	}

	return prefix
}

func getExecPath() []string {
	envPath := os.Getenv("PATH")

	paths := strings.Split(envPath, string(os.PathListSeparator))

	var execAble []string
	seen := make(map[string]bool)
	for _, b := range builtins {
		seen[b] = true
	}

	for _, dir := range paths {
		entrys, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entrys {
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			name := entry.Name()

			// check if executable
			if info.Mode()&0111 != 0 && !seen[name] {
				execAble = append(execAble, name)
				seen[name] = true
			}
		}
	}
	return execAble
}
