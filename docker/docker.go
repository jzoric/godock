package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"log"
)

type docker struct {
	client *client.Client
}

func NewDocker() (*docker, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &docker{client: dockerClient}, nil
}

type Item struct {
	ID   string
	Tag  string
	Size string
}

func (d *docker) Close() {
	if err := d.client.Close(); err != nil {
		log.Fatal(err)
	}
}

func (d docker) GetAllImages() []Item {
	images, _ := d.client.ImageList(context.Background(), types.ImageListOptions{All: true})

	var listImages []Item

	for _, image := range images {
		ins, _, _ := d.client.ImageInspectWithRaw(context.Background(), image.ID)
		if len(ins.RepoTags) == 0 {

		} else {
			for _, tag := range ins.RepoTags {
				listImages = append(listImages, Item{
					ID:   image.ID[7:20],
					Tag:  tag,
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

func (d docker) GetAllContainers() []Item {
	containers, _ := d.client.ContainerList(context.Background(), types.ContainerListOptions{
		Quiet:   false,
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.Args{},
	})

	var listImages []Item

	for _, container := range containers {
		listImages = append(listImages, Item{
			ID:  container.ID[:20],
			Tag: container.Names[0],
		})
	}

	return listImages
}
