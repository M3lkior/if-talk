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
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// UserSettingsPath is where code-server reads its user settings, relative to
// the container's home directory.
const UserSettingsPath = ".local/share/code-server/User/settings.json"

// homeDir is the home directory of the user code-server runs as in the image.
const homeDir = "/home/coder"

// userSettings is the settings file copied into every container demoit
// creates.
//
// zenMode.restore is the one that matters. VS Code Web keeps its layout in the
// browser, not in the container, so Zen Mode — one Cmd+K Z away, and easy to
// hit by accident on stage — came back on every later load: no file tree, no
// activity bar, no status bar, and nothing in the deck to say why. The setting
// is read from the container, which is the only side demoit controls.
//
// workbench.startupEditor keeps the "Get Started" walkthrough from opening over
// the demo and hiding the files the slide exists to show.
const userSettings = `{
    "zenMode.restore": false,
    "workbench.startupEditor": "none"
}
`

// UserSettings returns a tar archive of the code-server settings, to be
// extracted in the container's home directory.
//
// It carries a directory entry for every parent of the settings file: Docker
// creates a missing parent as root, and code-server, which runs as uid:gid,
// writes the rest of its state next to settings.json.
func UserSettings(uid, gid int) ([]byte, error) {
	var archive bytes.Buffer
	writer := tar.NewWriter(&archive)

	dir := ""
	for _, segment := range strings.Split(path.Dir(UserSettingsPath), "/") {
		dir = path.Join(dir, segment)

		if err := writer.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     dir + "/",
			Mode:     0o755,
			Uid:      uid,
			Gid:      gid,
		}); err != nil {
			return nil, fmt.Errorf("unable to archive the %s directory: %w", dir, err)
		}
	}

	if err := writer.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     UserSettingsPath,
		Mode:     0o644,
		Size:     int64(len(userSettings)),
		Uid:      uid,
		Gid:      gid,
	}); err != nil {
		return nil, fmt.Errorf("unable to archive %s: %w", UserSettingsPath, err)
	}
	if _, err := writer.Write([]byte(userSettings)); err != nil {
		return nil, fmt.Errorf("unable to archive %s: %w", UserSettingsPath, err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("unable to close the settings archive: %w", err)
	}

	return archive.Bytes(), nil
}

// CopyUserSettings extracts the code-server settings in the home directory of
// a container.
func CopyUserSettings(ctx context.Context, docker client.CommonAPIClient, containerID string, uid, gid int) error {
	archive, err := UserSettings(uid, gid)
	if err != nil {
		return err
	}

	if err := docker.CopyToContainer(ctx, containerID, homeDir, bytes.NewReader(archive), types.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("unable to write the vscode settings in the container: %w", err)
	}

	return nil
}
