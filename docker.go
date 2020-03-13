package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"os"
	"os/exec"
	"strings"
)


type dockerImage struct {
	ID  string
	Tag string
	Size string
}

func (g godock) GetAllImages() []dockerImage {
	images, _ := g.client.ImageList(context.Background(), types.ImageListOptions{All: true})

	var listImages []dockerImage

	for _, image := range images {
		ins, _, _ := g.client.ImageInspectWithRaw(context.Background(), image.ID)
		if len(ins.RepoTags) == 0 {

		} else {
			for _, tag := range ins.RepoTags {
				listImages = append(listImages, dockerImage{
					ID:  image.ID[7:20],
					Tag: tag,
					Size: getHumanReadableSize(image.Size),
				})
			}

		}
	}
	return listImages
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


func (g godock) GetAllContainers() []dockerImage {
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

	var listImages []dockerImage

	for _, container := range containers {
		listImages = append(listImages, dockerImage{
			ID:  container.ID[:20],
			Tag: container.Names[0],
		})
	}

	return listImages
}


func executor(s string) error {
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
