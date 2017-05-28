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

type Decorator interface {
	Decorate(message string) string
}

type SandwichDecorator struct {
	separator string
}

func (d *SandwichDecorator) Decorate(message string) string {
	return strings.Join([]string{d.separator, message, d.separator}, "\n")
}

type SdDecorator struct {
}

func (d *SdDecorator) Decorate(message string) string {
	decorated, err := exec.Command("echo-sd", message).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(decorated)
}

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

	var decorator Decorator = &SdDecorator{}
	//var decorator Decorator = &SandwichDecorator{separator: "-----------------------"}
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
				fmt.Println(decorator.Decorate(strings.Join(message, "\n")))
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

