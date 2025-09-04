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
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
	}

	for {
		fmt.Print("$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			return
		}

		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}
		cmd, arg, _ := strings.Cut(command, " ")
		args := parseArgs(arg)
		if handler, ok := builtinCommands[cmd]; ok {
			handler(args)
		} else {
			executeExternalCommand(cmd, args)
		}
	}
}

func parseArgs(inputString string) []string {
	var words []string
	var insideSingleQuotes bool
	var insideDoubleQuotes bool
	var escaped bool
	var currentWordBuilder strings.Builder

	for _, char := range inputString {
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
		case char == '"':
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

func handleType(args []string) {
	if len(args) < 1 {
		fmt.Println("type: missing argument")
		return
	}

	builtinCommands := map[string]bool{
		"echo": true,
		"exit": true,
		"type": true,
		"pwd":  true,
		"cd":   true,
	}

	if builtinCommands[args[0]] {
		fmt.Println(args[0] + " is a shell builtin")
	} else if path, err := exec.LookPath(args[0]); err == nil {
		fmt.Println(args[0], "is", path)
	} else {
		fmt.Println(args[0] + ": not found")
	}
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleExit(args []string) {
	if len(args) > 0 {
		if code, err := strconv.Atoi(args[0]); err == nil {
			os.Exit(code)
		} else {
			fmt.Println("exit: numeric argument required")
			os.Exit(1)
		}
	}
	os.Exit(0)
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

	err := os.Chdir(target)
	if err != nil {
		fmt.Println("cd: " + target + ": No such file or directory")
	}
}

func executeExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(command + ": command not found")
		return
	}
	fmt.Print(string(stdout))
}
