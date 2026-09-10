/*
Copyright 2018 Google LLC
Copyright 2022 David Gageot

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vscode

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/dgageot/demoit/files"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"golang.org/x/term"
)

const (
	Port          = 18080
	dockerImage   = "codercom/code-server:4.8.3@sha256:c0e99db852b2c4e3602c912b658d1fd509d02643a1c437d1f7bb80197535cd76"
	containerName = "demoit-vscode"
	appDir        = "/app"
)

var defaultFlags = []string{"--auth=none", "--disable-telemetry", "--disable-update-check", "--force"}

// Start makes sure a code-server container is serving this presentation's
// folder, and reports why it could not.
//
// It runs on every VS Code slide rather than once per process. Running once
// meant a single failure — Docker not up yet, or another demoit having taken
// the container — disabled VS Code for the rest of the run with no way back.
func Start() error {
	return startVsCodeServer(context.Background())
}

// existing reports whether a demoit-vscode container exists, whether it is
// running, and which host directory it currently has bound at /app.
func existing(ctx context.Context, docker client.CommonAPIClient) (bool, bool, string) {
	info, err := docker.ContainerInspect(ctx, containerName)
	if err != nil {
		return false, false, ""
	}

	bind := ""
	for _, mount := range info.Mounts {
		if mount.Destination == appDir {
			bind = mount.Source
			break
		}
	}

	return true, info.State != nil && info.State.Running, bind
}

// BindSource returns the host directory to bind at /app inside the code-server
// container: the presentation folder, resolved against the working directory
// when it is relative.
//
// filepath.Join(cwd, root) is not a substitute. For an absolute root it
// concatenates the two, so `demoit /path/to/talk` bound `<cwd>/path/to/talk` —
// a directory that does not exist. Docker then created an empty one, the
// container exited, and the slide rendered VS Code with no file tree and no
// error anywhere. The folder is checked here so that failure is reported
// rather than discovered on stage.
func BindSource(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("can't resolve the presentation folder %q: %w", root, err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("can't read the presentation folder %q: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("the presentation folder %q is not a directory", abs)
	}

	return abs, nil
}

// Action is what Start has to do about the container that may already exist.
type Action int

const (
	// Create means no usable container exists and one must be created.
	Create Action = iota
	// Reuse means a running container already serves the right folder.
	Reuse
)

// Plan decides what to do about an existing demoit-vscode container.
//
// It returns an error when another demoit process is already serving a
// different presentation folder: the container name and the host port are
// shared, and two code-servers cannot both hold that port. Taking the
// container over used to be silent, which meant the other talk's files
// appeared in this deck's VS Code slides while the instance that lost the
// race never recovered — Start ran once per process and never tried again.
// Refusing keeps the failure legible and lets whichever instance owns the
// container keep working.
func Plan(found, running bool, existingBind, wantBind string) (Action, error) {
	if !found || !running {
		return Create, nil
	}

	if existingBind != wantBind {
		// The container outlives the demoit that created it, so stopping the
		// other instance is not enough on its own — say how to release it.
		return Create, fmt.Errorf(
			"the %s container is already serving VS Code for %q on port %d, so this presentation (%q) cannot: stop the other demoit, then run `docker rm -f %s`",
			containerName, existingBind, Port, wantBind, containerName)
	}

	return Reuse, nil
}

func startVsCodeServer(ctx context.Context) error {
	bind, err := BindSource(files.Root)
	if err != nil {
		return err
	}

	client, err := newDockerlient(ctx)
	if err != nil {
		return fmt.Errorf("docker is not available: %w", err)
	}

	found, running, currentBind := existing(ctx, client)

	action, err := Plan(found, running, currentBind, bind)
	if err != nil {
		return err
	}
	if action == Reuse {
		return nil
	}

	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("can't get current user: %w", err)
	}

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return fmt.Errorf("can't read the current user id %q: %w", user.Uid, err)
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return fmt.Errorf("can't read the current group id %q: %w", user.Gid, err)
	}

	// Ignore error
	_ = client.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true})

	rc, err := client.ImagePull(ctx, dockerImage, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("unable to initiate vscode image pull: %w", err)
	}
	defer rc.Close()
	if err := streamDockerMessages(rc); err != nil {
		return fmt.Errorf("unable to pull vscode image: %w", err)
	}

	body, err := client.ContainerCreate(ctx, &containertypes.Config{
		Image: dockerImage,
		User:  fmt.Sprintf("%d:%d", uid, gid),
		Cmd:   defaultFlags,
		ExposedPorts: nat.PortSet{
			nat.Port("8080/tcp"): struct{}{},
		},
	}, &containertypes.HostConfig{
		Binds: []string{bind + ":" + appDir},
		PortBindings: nat.PortMap{
			nat.Port("8080/tcp"): []nat.PortBinding{{HostPort: strconv.Itoa(Port)}},
		},
	}, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("unable to create a docker container for vscode: %w", err)
	}

	if err := CopyUserSettings(ctx, client, body.ID, uid, gid); err != nil {
		return err
	}

	if err := client.ContainerStart(ctx, body.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start a docker container for vscode: %w", err)
	}

	return nil
}

func newDockerlient(ctx context.Context) (client.CommonAPIClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error getting docker client: %w", err)
	}
	cli.NegotiateAPIVersion(ctx)

	return cli, nil
}

func streamDockerMessages(src io.Reader) error {
	dst := os.Stdout
	termFd, isTerm := isTerminal(dst)
	return jsonmessage.DisplayJSONMessagesStream(src, dst, termFd, isTerm, nil)
}

func isTerminal(w io.Writer) (uintptr, bool) {
	type descriptor interface {
		Fd() uintptr
	}

	if f, ok := w.(descriptor); ok {
		termFd := f.Fd()
		isTerm := term.IsTerminal(int(termFd))
		return termFd, isTerm
	}

	return 0, false
}
