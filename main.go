package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/docker/docker/client"
	"log"
	"os"
	"os/exec"
	"strings"
)

type godock struct {
	client *client.Client
}

func NewGodock() (*godock, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &godock{client: dockerClient}, nil
}

func (g godock) executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "exit" {
		fmt.Println("Cu!")
		os.Exit(0)
		return
	}

	cmd := exec.Command("/bin/sh", "-c", "docker "+s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	}
	return
}

func (g godock) completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return prompt.FilterHasPrefix(g.mainSuggestions(), d.GetWordBeforeCursor(), true)
	}
	args := parseArgs(d.TextBeforeCursor())

	switch args[0] {
	case "rmi":
		return g.imagesSuggestion()
	case "rm":
		return g.containersSuggestion()
	}
	return []prompt.Suggest{}
}

func parseArgs(t string) []string {
	splits := strings.Split(t, " ")
	args := make([]string, 0, len(splits))

	for i := range splits {
		if i != len(splits)-1 && splits[i] == "" {
			continue
		}
		args = append(args, splits[i])
	}
	return args
}

func main() {
	goDock, err := NewGodock()
	if err != nil {
		log.Fatalf("can't start godock: %+v", err)
	}
	fmt.Printf("godock version: %s\n", "0.0.1")
	fmt.Println("Please use `exit` to exit!")
	defer fmt.Println("Bye!")

	for {
		dockerCommand := prompt.New(goDock.executor, goDock.completer)
		dockerCommand.Run()
	}
}
