```markdown
# gsh

Get some SHell
or
Go Shell

A simple shell implementation written in Go.

## Features

- Interactive command-line interface
- Command history support
- Pipeline support (command chaining with pipes)
- Built-in command completion
- Persistent history across sessions

## Installation

```bash
go install github.com/Aergiaaa/gsh@latest
```

Or build from source:

```bash
git clone https://github.com/Aergiaaa/gsh.git
cd gsh
make build
```

## Usage

Run the shell:

```bash
gsh
```

### Environment Variables

- `HISTFILE`: Set custom history file location (default: `.gsh_history`)

### Built-in Commands

echo
type
history
pwd
cd
exit

## Examples

```bash
# Run a simple command
> ls -la

# Use pipes
> cat file.txt | grep "pattern"
```

## Development

### Prerequisites

- Go 1.25.5 or higher

### Building

```bash
make install
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 

## Author

Aergiaaa
```
