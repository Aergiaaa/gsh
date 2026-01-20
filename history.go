package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const MAX_HIST_SIZE = 1000

var (
	HISTFILE      = ".gsh_history"
	lastAppendIdx = make(map[string]int)
)

func history(args []string, r io.Reader, w io.Writer, hist *[]string) {
	defer drainStdin(r)()

	if hist == nil {
		return
	}

	// Handle -r flag for reading external history
	if len(args) >= 2 && args[1] != "" {
		switch args[0] {
		case "-r":
			extHist := loadHist(args[1])
			*hist = append(*hist, extHist...)
			return
		case "-w":
			saveHist(*hist, args[1])
			return
		case "-a":
			fPath := args[1]
			lastIdx := lastAppendIdx[fPath]
			newCmds := (*hist)[lastIdx:]
			appendHistFile(newCmds, fPath)
			lastAppendIdx[fPath] = len(*hist)
			return
		}
	}

	start := 0

	// If an argument is provided, show last N commands
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			start = max(len(*hist)-n, 0)
		}
	}

	for i := start; i < len(*hist); i++ {
		fmt.Fprintf(w, "%5d  %s\n", i+1, (*hist)[i])
	}
}

// loadHistory reads history from file on startup
func loadHist(path string) []string {
	var history []string

	histPath := getHistoryPath(path)
	file, err := os.Open(histPath)
	if err != nil {
		// File doesn't exist yet, return empty history
		return history
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			history = append(history, line)
		}
	}

	return history
}

func saveHist(history []string, path string) {
	histPath := getHistoryPath(path)

	// Keep only last N commands
	start := 0
	if len(history) > MAX_HIST_SIZE {
		start = len(history) - MAX_HIST_SIZE
	}

	file, err := os.Create(histPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := start; i < len(history); i++ {
		fmt.Fprintln(writer, history[i])
	}
	writer.Flush()
}

// appendHistory appends a single command to history file immediately
func appendHist(cmd, path string) {
	histPath := getHistoryPath(path)

	file, err := os.OpenFile(histPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintln(file, cmd)
}

func appendHistFile(cmd []string, path string) {
	histPath := getHistoryPath(path)

	file, err := os.OpenFile(histPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	for _, v := range cmd {
		fmt.Fprintln(file, v)
	}

}

// getHistoryPath returns the full path to history file
func getHistoryPath(path string) string {
	// If path is absolute, use it as-is
	if filepath.IsAbs(path) {
		return path
	}

	// If file exists at the given path, use it
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// Check if file exists in current directory
	currDir, err := os.Getwd()
	if err == nil {
		currPath := filepath.Join(currDir, path)
		if _, err := os.Stat(currPath); err == nil {
			return currPath
		}
	}

	// For HISTFILE constant, use home directory
	if path == HISTFILE {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path)
		}
	}

	// Default: return path as-is (will be created if needed)
	return path
}
