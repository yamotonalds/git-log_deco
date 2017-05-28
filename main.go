package main

import (
	"fmt"
	"flag"
	"os/exec"
	"log"
	"bufio"
	"io"
	"strings"
)

func main() {
	flag.Parse()

	args := []string{"log"}
	args = append(args, flag.Args()...)
	logCmd := exec.Command("git", args...)
	rawOutput, err := logCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	output := bufio.NewReader(rawOutput)

	if err := logCmd.Start(); err != nil {
		log.Fatal(err)
	}

	message := []string{}
	isMessageLine := false
	for {
                rawLine, _, err := output.ReadLine()
                line := string(rawLine)
		if line == "" {
			if len(message) == 0 {
				isMessageLine = true
				fmt.Println("")
			} else {
				isMessageLine = false
				decorated, cmdErr := exec.Command("echo-sd", strings.Join(message, "\n")).CombinedOutput()
				if cmdErr != nil {
					log.Fatal(cmdErr)
				}
				fmt.Println(string(decorated))
				message = []string{}
			}		
		} else {
			if isMessageLine {
				message = append(message, strings.TrimSpace(line))
			} else {
				fmt.Println(line)
			}
		}

                if err == io.EOF {
                        break
                } else if err != nil {
                        log.Fatal(err)
                }

        }
}

