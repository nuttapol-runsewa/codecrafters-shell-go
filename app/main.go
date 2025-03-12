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

		args := strings.Split(command, " ")
		if handler, ok := builtinCommands[args[0]]; ok {
			handler(args)
		} else {
			executeExternalCommand(args)
		}
	}
}

func handleType(args []string) {
	if len(args) < 2 {
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

	if builtinCommands[args[1]] {
		fmt.Println(args[1] + " is a shell builtin")
	} else if path, err := exec.LookPath(args[1]); err == nil {
		fmt.Println(args[1], "is", path)
	} else {
		fmt.Println(args[1] + ": not found")
	}
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args[1:], " "))
}

func handleExit(args []string) {
	if len(args) > 1 {
		if code, err := strconv.Atoi(args[1]); err == nil {
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
	if len(args) < 2 {
		fmt.Println("cd: missing argument")
		return
	}

	target := args[1]
	if target == "~" {
		target = os.Getenv("HOME")
	}

	err := os.Chdir(target)
	if err != nil {
		fmt.Println("cd: " + target + ": No such file or directory")
	}
}

func executeExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(args[0] + ": command not found")
		return
	}
	fmt.Print(string(stdout))
}
