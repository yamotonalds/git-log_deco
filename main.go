package main

import (
	"fmt"
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

type CommandDecorator struct {
	command string
}

func (d *CommandDecorator) Decorate(message string) string {
	decorated, err := exec.Command(d.command, message).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(decorated)
}

func main() {
	cmd := exec.Command("git", "log")
	rawOutput, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	output := bufio.NewReader(rawOutput)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var decorator Decorator = &CommandDecorator{command: "echo-sd"}
	//var decorator Decorator = &SandwichDecorator{separator: "-----------------------"}
	messageLines := []string{}
	isMessageLine := false
	for {
                rawLine, _, err := output.ReadLine()
                line := string(rawLine)
		if line == "" {
			if len(messageLines) == 0 {
				isMessageLine = true
				fmt.Println("")
			} else {
				isMessageLine = false
				fmt.Println(decorator.Decorate(strings.Join(messageLines, "\n")))
				messageLines = []string{}
			}		
		} else {
			if isMessageLine {
				messageLines = append(messageLines, strings.TrimSpace(line))
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

