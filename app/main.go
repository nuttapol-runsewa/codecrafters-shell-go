package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	builtinCommands := map[string]func([]string){
		"echo": handleEcho,
		"exit": handleExit,
		"type": nil, // placeholder
		"pwd":  handlePwd,
		"cd":   handleCd,
	}
	builtinCommands["type"] = func(args []string) { handleType(args, builtinCommands) }
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		cmd, arg, _ := strings.Cut(line, " ")
		args := parseArgs(arg)
		if handler, ok := builtinCommands[cmd]; ok {
			handler(args)
		} else {
			executeExternalCommand(cmd, args)
		}
	}
}

func parseArgs(input string) []string {
	var (
		words              []string
		insideSingleQuotes bool
		insideDoubleQuotes bool
		escaped            bool
		currentWordBuilder strings.Builder
	)

	for _, char := range input {
		switch {
		case char == '\\' && !insideSingleQuotes && !insideDoubleQuotes:
			escaped = true
		case char == '\\' && insideSingleQuotes:
			currentWordBuilder.WriteRune(char)
			escaped = true
		case escaped:
			currentWordBuilder.WriteRune(char)
			escaped = false
		case char == '\'' && !insideDoubleQuotes:
			insideSingleQuotes = !insideSingleQuotes
		case char == '"' && !insideSingleQuotes:
			insideDoubleQuotes = !insideDoubleQuotes
		case char == ' ' && !insideSingleQuotes && !insideDoubleQuotes:
			if currentWordBuilder.Len() > 0 {
				words = append(words, currentWordBuilder.String())
				currentWordBuilder.Reset()
			}
		default:
			currentWordBuilder.WriteRune(char)
		}
	}
	if currentWordBuilder.Len() > 0 {
		words = append(words, currentWordBuilder.String())
	}
	return words
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleExit(args []string) {
	code := 0
	if len(args) > 0 {
		if val, err := strconv.Atoi(args[0]); err == nil {
			code = val
		} else {
			fmt.Println("exit: numeric argument required")
			code = 1
		}
	}
	os.Exit(code)
}

func handleType(args []string, builtinCommands map[string]func([]string)) {
	if len(args) < 1 {
		fmt.Println("type: missing argument")
		return
	}
	if _, ok := builtinCommands[args[0]]; ok {
		fmt.Printf("%s is a shell builtin\n", args[0])
		return
	}
	if path, err := exec.LookPath(args[0]); err == nil {
		fmt.Printf("%s is %s\n", args[0], path)
	} else {
		fmt.Printf("%s: not found\n", args[0])
	}
}

func handlePwd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		return
	}
	fmt.Println(dir)
}

func handleCd(args []string) {
	if len(args) < 1 {
		fmt.Println("cd: missing argument")
		return
	}
	target := args[0]
	if target == "~" {
		target = os.Getenv("HOME")
	}
	if err := os.Chdir(target); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", target)
	}
}

func executeExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s: command not found\n", command)
		return
	}
	fmt.Print(string(output))
}
