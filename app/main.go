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

		cmd, args := parseArgs(command)
		if handler, ok := builtinCommands[cmd]; ok {
			handler(args)
		} else {
			executeExternalCommand(cmd, args)
		}
	}
}

func parseArgs(command string) (string, []string) {
	var args []string
	var currentArg strings.Builder
	inQuotes := false
	var cmd string
	firstWord := true

	for _, char := range command {
		if char == '\'' && !inQuotes {
			inQuotes = true
		} else if char == '\'' && inQuotes {
			inQuotes = false
			unquoted, err := strconv.Unquote("'" + currentArg.String() + "'")
			if err != nil {
				args = append(args, currentArg.String())
			} else {
				args = append(args, unquoted)
			}
			currentArg.Reset()
		} else if char == ' ' && !inQuotes {
			if currentArg.Len() > 0 {
				if firstWord {
					cmd = currentArg.String()
					firstWord = false
				} else {
					args = append(args, currentArg.String())
				}
				currentArg.Reset()
			}
		} else {
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		if firstWord {
			cmd = currentArg.String()
		} else {
			args = append(args, currentArg.String())
		}
	}
	if firstWord {
		cmd = currentArg.String()
	}

	return cmd, args
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
		fmt.Println(args[0] + ": command not found")
		return
	}
	fmt.Print(string(stdout))
}
