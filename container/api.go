package container

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

/* source examples : https://docs.docker.com/develop/sdk/examples/
https://github.com/moby/moby/tree/master/client

API references : https://godoc.org/github.com/docker/docker/client
*/

// DockerClient is the docker API client
func DockerClient() client.APIClient {
	// create client with docker verision 1.39
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	return cli
}

// ListContainers list all containers of the current docker host
func ListContainers() {
	cli := DockerClient()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

// ResponseCreatedBody is the response of the container creation request
type ResponseCreatedBody struct {

	// The ID of the created container
	// Required: true
	ID string `json:"Id"`

	// Warnings encountered when creating the container
	// Required: true
	Warnings []string `json:"Warnings"`
}

// RunContainer run a container in background
func RunContainer(name string) ResponseCreatedBody {
	ctx := context.Background()
	cli := DockerClient()

	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	createdBody := ResponseCreatedBody{}
	createdBody.ID = resp.ID
	createdBody.Warnings = resp.Warnings

	return createdBody
}

// PrintLogsOfContainer show logs of a container
func PrintLogsOfContainer(containerID string) io.ReadCloser {
	ctx := context.Background()
	cli := DockerClient()

	options := types.ContainerLogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)

	return out
}

// StopAllRunningContainers stop all containers of current docker host
func StopAllRunningContainers() {
	ctx := context.Background()
	cli := DockerClient()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}

// ListAllImages list all images of current docker host
func ListAllImages() []types.ImageSummary {
	cli := DockerClient()

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	return images
}

// PullImage pull an image from docker hub registry
func PullImage(imageName string) io.ReadCloser {
	ctx := context.Background()
	cli := DockerClient()

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	return out
}

// PullImageWithAuth pull an image from docker hub registry
func PullImageWithAuth(imageName string, username string, password string) io.ReadCloser {
	ctx := context.Background()
	cli := DockerClient()

	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	return out
}

// CommitContainer commit a container to create an image from its contents
func CommitContainer(imageName string, cmd []string, newImageName string) types.IDResponse {
	ctx := context.Background()
	cli := DockerClient()

	createResp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   cmd,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	commitResp, err := cli.ContainerCommit(ctx, createResp.ID, types.ContainerCommitOptions{Reference: newImageName})
	if err != nil {
		panic(err)
	}

	return commitResp
}
