package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		builtinCommands := map[string]bool{
			"echo": true,
			"exit": true,
			"type": true,
			"pwd":  true,
			"cd":   true,
		}

		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command = strings.TrimSpace(command)
		args := strings.Split(command, " ")
		switch args[0] {
		case "type":
			if builtinCommands[args[1]] {
				fmt.Println(args[1] + " is a shell builtin")
			} else if path, err := exec.LookPath(args[1]); err == nil {
				fmt.Println(args[1], "is", path)
			} else {
				fmt.Println(args[1] + ": not found")
			}
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "exit":
			code, _ := strconv.Atoi(args[1])
			os.Exit(code)
		case "pwd":
			dir, _ := os.Getwd()
			fmt.Println(dir)
		case "cd":
			if args[1] == "~" {
				args[1] = os.Getenv("HOME")
			}
			err := os.Chdir(args[1])
			if err != nil {
				fmt.Println("cd: " + args[1] + ": No such file or directory")
			}
		default:
			cmd := exec.Command(args[0], args[1:]...)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(args[0] + ": command not found")
			} else {
				fmt.Print(string(stdout))
			}
		}
	}
}
