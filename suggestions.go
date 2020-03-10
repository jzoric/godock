package main

import (
	"context"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (g godock) imagesSuggestion() []prompt.Suggest {
	images, _ := g.client.ImageList(context.Background(), types.ImageListOptions{All: true})
	var suggestions []prompt.Suggest

	for _, image := range images {
		ins, _, _ := g.client.ImageInspectWithRaw(context.Background(), image.ID)
		if len(ins.RepoTags) == 0 {

		} else {
			for _, tag := range ins.RepoTags {
				suggestions = append(suggestions, prompt.Suggest{Text: image.ID[7:], Description: tag + " (" + getHumanReadableSize(image.Size) + ")"})
			}

		}
	}
	return suggestions
}

func getHumanReadableSize(size int64) string {
	const unit = 1000
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "kMGTPE"[exp])
}

func (g godock) mainSuggestions() []prompt.Suggest {
	var suggestions []prompt.Suggest
	suggestions = append(suggestions, prompt.Suggest{Text: "rm", Description: "Remove containers"})
	return append(suggestions, prompt.Suggest{Text: "rmi", Description: "Remove docker image"})
}

func (g godock) containersSuggestion() []prompt.Suggest {
	containers, _ := g.client.ContainerList(context.Background(), types.ContainerListOptions{
		Quiet:   false,
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.Args{},
	})

	var suggestions []prompt.Suggest

	for _, container := range containers {
		suggestions = append(suggestions, prompt.Suggest{Text: container.ID, Description: container.Names[0]})
	}

	return suggestions
}
