package main

import (
	"fmt"
	"io"
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
		inpipe, err := cmd.StdinPipe()
		if err != nil {
			fmt.Printf("Failed to get program stdin: %v. Exiting.\n", err)
			break
		}
		outpipe, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Failed to get program stdout: %v. Exiting.\n", err)
			break
		}
		errpipe, err := cmd.StderrPipe()
		if err != nil {
			fmt.Printf("Failed to get program stderr: %v. Exiting.\n", err)
			break
		}
		err = cmd.Start()
		if err != nil {
			fmt.Printf("Failed to start program: %v. Exiting.\n", err)
			break
		}
		//err,in and out are passed as parameter to avoid race conditions among
		//different iterations of the loop
		go func(inpipe io.WriteCloser) {
			_, _ = io.Copy(inpipe, os.Stdin)
			_ = inpipe.Close()
		}(inpipe)
		go func(outpipe io.ReadCloser) {
			_, _ = io.Copy(os.Stdout, outpipe)
			_ = outpipe.Close()
		}(outpipe)
		go func(errpipe io.ReadCloser) {
			_, _ = io.Copy(os.Stderr, errpipe)
			_ = errpipe.Close()
		}(errpipe)

		fmt.Printf("Successfully started program. Monitoring...\n")
		err = cmd.Wait()
		fmt.Printf("Program exited with error %v.\n", err)
	}
}
