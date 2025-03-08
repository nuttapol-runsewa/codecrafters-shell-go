package main

import (
	"bufio"
	"fmt"
	"os"
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
			} else {
				fmt.Println(args[1] + ": not found")
			}
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "exit":
			code, _ := strconv.Atoi(args[1])
			os.Exit(code)
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
