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

		switch command {
		case "exit 0":
			os.Exit(0)
		default:
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
