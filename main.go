package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
	"log"
	"os"
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

func main() {
	goDock, err := NewGodock()
	if err != nil {
		log.Fatalf("can't start godock: %+v", err)
	}

	command, err := commands()
	if err != nil {
		fmt.Println(err)
	}

	switch command {
	case 0:
		itemList := goDock.GetAllImages()
		details := `
--------- Details ----------
{{ "ID:" | faint }}	{{ .ID }}
{{ "Tag:" | faint }}	{{ .Tag }}
{{ "Size:" | faint }}	{{ .Size }}`
		dockerImageId, err := getSelectedItem(itemList, details)
		if err != nil {
			os.Exit(2)
		}
		fmt.Println(dockerImageId)
	case 1:
		allImages := goDock.GetAllContainers()
		details := `
--------- Details ----------
{{ "ID:" | faint }}	{{ .ID }}
{{ "Tag:" | faint }}	{{ .Tag }}`
		if len(allImages) == 0 {
			fmt.Println("No containers...")
			return
		}
		dockerContainerId, err := getSelectedItem(allImages,details)
		if err != nil {
			os.Exit(2)
		}
		fmt.Println(dockerContainerId)
	case 3:
		os.Exit(0)
	}

}

func getSelectedItem(itemList []dockerImage, details string) (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   ">> {{ .ID | cyan }} ({{ .Tag | red }})",
		Inactive: "  {{ .ID | cyan }} ({{ .Tag | red }})",
		Selected: ">> {{ .ID | red | cyan }}",
		Details: details,
	}

	searcher := func(input string, index int) bool {
		image := itemList[index]
		name := strings.Replace(strings.ToLower(image.Tag), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "",
		Items:     itemList,
		Templates: templates,
		Size:      5,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	return itemList[i].ID, err
}

func commands() (int, error) {
	prompt := promptui.Select{
		Label: "",
		Items: []string{"rmi", "rm", "exit"},
		Size:  5,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1, err
	}
	return i, nil
}
