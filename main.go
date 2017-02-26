package main

import (
	"fmt"
	"os"
	"os/exec"
)

func createCmd() *exec.Cmd {
	// Use bash to parse given executable path and args, so we don't have
	// to care about removing quotes, ...
	var bashCmdArgs []string
	bashCmdArgs = append(bashCmdArgs, "-c")
	bashCmdArgs = append(bashCmdArgs, os.Args[1:]...)
	return exec.Command("bash", bashCmdArgs...)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid args")
		os.Exit(1)
	}

	for {
		cmd := createCmd()
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Failed to start program: %v. Exiting.\n", err)
			break
		}
		fmt.Printf("Successfully started program. Monitoring...\n")
		err = cmd.Wait()
		fmt.Printf("Program exited with error %v.\n", err)
	}
}
