package main

import (
	"fmt"
	"io"
	"os"
)

func notExist(cmd string) {
	fmt.Printf("%s: command not found\n", cmd)
}

func parseQuote(str string) (arr []string) {
	var res []byte
	var (
		quote    = false
		dquote   = false
		backlash = false
	)

	var (
		quotePrior  = 0
		dquotePrior = 0
	)

	for i := range str {
		if backlash {
			res = append(res, str[i])
			backlash = false
			continue
		}
		if str[i] == '\\' {
			if !(len(str) > i+1) {
				res = append(res, str[i])
				continue
			}

			next := str[i+1]
			if dquote {
				if isSpecialChar(next) {
					backlash = true
					continue
				}

				res = append(res, str[i])
				continue
			}

			if quote {
				res = append(res, str[i])
			}

			backlash = true
			continue
		}

		if str[i] == '"' {
			if quote && quotePrior > dquotePrior {
				res = append(res, str[i])
			}
			if dquote {
				dquote = false
				if quote && dquotePrior > quotePrior {
					quotePrior = 0
					quote = false
				}
				dquotePrior = 0
				continue
			}
			dquote = true
			dquotePrior++
			if !quote {
				dquotePrior++
			}
			continue
		}

		if str[i] == '\'' {
			if dquote && dquotePrior > quotePrior {
				res = append(res, str[i])
			}
			if quote {
				quote = false
				if dquote && quotePrior > dquotePrior {
					dquotePrior--
					dquote = false
				}
				quotePrior--
				continue
			}
			quote = true
			quotePrior++
			if !dquote {
				quotePrior++
			}
			continue
		}
		if quote {
			res = append(res, str[i])
			continue
		}

		if dquote {
			res = append(res, str[i])
			continue
		}

		if str[i] == ' ' {
			if len(res) > 0 {

				arr = append(arr, string(res))
				res = nil
			}
			continue
		}

		res = append(res, str[i])
	}

	// edge cases for if the end of the line
	// is straight newline
	if len(res) > 0 {
		arr = append(arr, string(res))
	}

	return arr
}

func isSpecialChar(b byte) bool {
	return b == '"' ||
		b == '\\' ||
		b == '$' ||
		b == '`'
}

func isRedirect(str string) bool {
	return str == ">" ||
		str == "1>" ||
		str == "2>" ||
		str == ">>" ||
		str == "1>>" ||
		str == "2>>"
}

func drainStdin(stdin io.Reader) func() {
	if stdin == os.Stdin {
		return func() {} // no-op if it's terminal stdin
	}

	done := make(chan bool)
	go func() {
		io.Copy(io.Discard, stdin)
		done <- true
	}()

	return func() { <-done }
}
