package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command = strings.TrimSpace(command)
		commands := strings.Split(command, " ")
		switch commands[0] {
		case "echo":
			fmt.Println(strings.Join(commands[1:], " "))
		case "exit":
			os.Exit(0)
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
