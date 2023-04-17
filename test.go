package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func test() {
	cmd := exec.Command("sh", "-c", "cd /src/oddjob/ssm; ./port_forward.sh -h develop")
	// some command output will be input into stderr
	// e.g.
	// cmd := exec.Command("../../bin/master_build")
	// stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()
	fmt.Println("The command is running")
	if err != nil {
		fmt.Println(err)
	}

	// print the output of the subprocess
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}
