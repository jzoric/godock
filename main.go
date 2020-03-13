package main

import (
	"github.com/manifoldco/promptui"
	"godock/docker"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	d, err := docker.NewDocker()
	defer d.Close()

	if err != nil {
		log.Fatalf("can't start godock: %+v", err)
	}

	command, err := getAllCommands()
	if err != nil {
		log.Fatalf("can't get all commands: %v", err)
	}

	switch command {
	case 0:
		itemList := d.GetAllImages()
		details := imageDetails()
		dockerImageId, err := getSelectedItem(itemList, details, "Images")
		if err != nil {
			log.Printf("can't get selected docker image: %v", err)
			return
		}
		err = run(" rmi "+dockerImageId)
		if err != nil {
			log.Printf("can't run the command: %v", err)
			return
		}
	case 1:
		allImages := d.GetAllContainers()
		details := containerDetails()
		if len(allImages) == 0 {
			log.Println("No containers...")
			return
		}
		dockerContainerId, err := getSelectedItem(allImages, details, "Containers")
		if err != nil {
			log.Printf("can't get selected docker container:: %v", err)
			return
		}
		err = run(" rm "+ dockerContainerId)
		if err != nil {
			log.Fatalf("can't run the command: %v", err)
		}
	case 3:
		os.Exit(0)
	}

}

func imageDetails() string {
	return `
--------- Details ----------
{{ "ID:" | faint }}	{{ .ID }}
{{ "Tag:" | faint }}	{{ .Tag }}
{{ "Size:" | faint }}	{{ .Size }}`
}

func containerDetails() string {
	return `
--------- Details ----------
{{ "ID:" | faint }}	{{ .ID }}
{{ "Tag:" | faint }}	{{ .Tag }}`
}

func getSelectedItem(itemList []docker.Item, details string, selectLabel string) (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   ">> {{ .ID | cyan }} ({{ .Tag | red }})",
		Inactive: "  {{ .ID | cyan }} ({{ .Tag | red }})",
		Selected: ">> {{ .ID | red | cyan }}",
		Details:  details,
	}

	searcher := func(input string, index int) bool {
		image := itemList[index]
		name := strings.Replace(strings.ToLower(image.Tag), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     selectLabel,
		Items:     itemList,
		Templates: templates,
		Size:      5,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	return itemList[i].ID, err
}

func getAllCommands() (int, error) {
	prompt := promptui.Select{
		Label: "Commands",
		Items: []string{
			"rmi - remove the docker image",
			"rm - remove container",
			"exit - exit godock"},
		Size: 5,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return -1, err
	}
	return i, nil
}

func run(s string) error {
	s = strings.TrimSpace(s)
	cmd := exec.Command("/bin/sh", "-c", "docker "+s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
