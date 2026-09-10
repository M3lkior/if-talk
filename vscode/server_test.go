package vscode_test

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dgageot/demoit/vscode"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestBindSourceResolvesARelativeFolder(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to read the working directory: %v", err)
	}

	got, err := vscode.BindSource(".")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	if got != cwd {
		t.Errorf("got %q, want %q", got, cwd)
	}
}

// An absolute presentation folder has to be used as it is. Joining it to the
// working directory concatenated the two, so `demoit /path/to/talk` bound
// `<cwd>/path/to/talk` at /app — a directory that does not exist. Docker then
// created an empty one, the container exited, and the slide showed VS Code
// with no file tree and no error anywhere.
func TestBindSourceKeepsAnAbsoluteFolder(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()

	got, err := vscode.BindSource(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	if got != folder {
		t.Errorf("got %q, want %q — an absolute folder must not be joined to the working directory", got, folder)
	}
}

// A bind source that does not exist makes Docker invent an empty directory and
// the container die, which is the same silent failure by another route.
func TestBindSourceRejectsAMissingFolder(t *testing.T) {
	t.Parallel()

	if _, err := vscode.BindSource(filepath.Join(t.TempDir(), "pas-la")); err == nil {
		t.Fatal("got no error for a folder that does not exist, want one")
	}
}

func TestBindSourceRejectsAFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "demoit.md")
	if err := os.WriteFile(path, []byte("# une slide\n"), 0o600); err != nil {
		t.Fatalf("unable to write the file: %v", err)
	}

	if _, err := vscode.BindSource(path); err == nil {
		t.Fatal("got no error for a file, want one — /app must be a directory")
	}
}

func TestPlanCreatesWhenNoContainerExists(t *testing.T) {
	t.Parallel()

	got, err := vscode.Plan(false, false, "", "/talks/mine")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	if got != vscode.Create {
		t.Errorf("got %v, want Create", got)
	}
}

func TestPlanReusesAContainerServingTheSameFolder(t *testing.T) {
	t.Parallel()

	got, err := vscode.Plan(true, true, "/talks/mine", "/talks/mine")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	if got != vscode.Reuse {
		t.Errorf("got %v, want Reuse — recreating it would restart code-server for nothing", got)
	}
}

// Two code-server instances cannot both hold the shared port, so a container
// serving another folder must not be taken over: doing so silently swapped the
// other talk's files into this deck's VS Code slides, and the instance that
// lost the race never recovered.
func TestPlanRefusesToStealAContainerServingAnotherFolder(t *testing.T) {
	t.Parallel()

	_, err := vscode.Plan(true, true, "/talks/autre", "/talks/mine")
	if err == nil {
		t.Fatal("got no error, want one — another demoit is serving that container")
	}
	for _, want := range []string{"/talks/autre", "/talks/mine"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got %q, want it to name %q so the conflict is obvious", err, want)
		}
	}
}

func TestPlanRecreatesAStoppedContainer(t *testing.T) {
	t.Parallel()

	for _, bind := range []string{"/talks/mine", "/talks/autre"} {
		bind := bind
		t.Run(bind, func(t *testing.T) {
			t.Parallel()

			got, err := vscode.Plan(true, false, bind, "/talks/mine")
			if err != nil {
				t.Fatalf("got error %v, want none — a stopped container serves nobody", err)
			}
			if got != vscode.Create {
				t.Errorf("got %v, want Create", got)
			}
		})
	}
}

// readArchive returns every entry of a tar archive, keyed by name.
func readArchive(t *testing.T, archive []byte) map[string]*tar.Header {
	t.Helper()

	entries := map[string]*tar.Header{}
	reader := tar.NewReader(bytes.NewReader(archive))
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unable to read the archive: %v", err)
		}
		entries[header.Name] = header
	}

	return entries
}

// readArchivedFile returns the content of one file of a tar archive.
func readArchivedFile(t *testing.T, archive []byte, name string) []byte {
	t.Helper()

	reader := tar.NewReader(bytes.NewReader(archive))
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unable to read the archive: %v", err)
		}
		if header.Name != name {
			continue
		}

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("unable to read %s: %v", name, err)
		}

		return content
	}

	t.Fatalf("the archive holds no %s, only %v", name, readArchive(t, archive))

	return nil
}

// VS Code Web keeps its layout in the browser, not in the container, so a
// speaker who once hit Zen Mode got it restored on every later load: no file
// tree, no activity bar, no status bar, and nothing in the deck to explain it.
// zenMode.restore is read from the container, which is the only side demoit
// controls.
func TestUserSettingsTurnOffZenModeRestore(t *testing.T) {
	t.Parallel()

	archive, err := vscode.UserSettings(501, 20)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(readArchivedFile(t, archive, vscode.UserSettingsPath), &settings); err != nil {
		t.Fatalf("the settings are not valid JSON: %v", err)
	}

	if got := settings["zenMode.restore"]; got != false {
		t.Errorf(`got zenMode.restore = %v, want false — Zen Mode must not survive a reload`, got)
	}
}

// The "Get Started" walkthrough opens over the demo on a fresh container and
// hides the very files the slide is there to show.
func TestUserSettingsOpenNoStartupEditor(t *testing.T) {
	t.Parallel()

	archive, err := vscode.UserSettings(501, 20)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(readArchivedFile(t, archive, vscode.UserSettingsPath), &settings); err != nil {
		t.Fatalf("the settings are not valid JSON: %v", err)
	}

	if got := settings["workbench.startupEditor"]; got != "none" {
		t.Errorf(`got workbench.startupEditor = %v, want "none"`, got)
	}
}

// Docker creates a missing parent of a copied file as root. code-server runs
// as the host user and writes its own state next to settings.json, so every
// directory the archive creates has to belong to that user.
func TestUserSettingsOwnEveryDirectoryTheyCreate(t *testing.T) {
	t.Parallel()

	archive, err := vscode.UserSettings(501, 20)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	entries := readArchive(t, archive)

	for _, name := range []string{
		".local/",
		".local/share/",
		".local/share/code-server/",
		".local/share/code-server/User/",
	} {
		header, ok := entries[name]
		if !ok {
			t.Fatalf("the archive holds no %s, so Docker creates it as root", name)
		}
		if header.Typeflag != tar.TypeDir {
			t.Errorf("%s is not a directory entry", name)
		}
		if header.Mode&0o700 != 0o700 {
			t.Errorf("got mode %#o for %s, want it writable by its owner", header.Mode, name)
		}
	}

	for name, header := range entries {
		if header.Uid != 501 || header.Gid != 20 {
			t.Errorf("got %s owned by %d:%d, want 501:20", name, header.Uid, header.Gid)
		}
	}
}

// fakeDocker records the one call CopyUserSettings makes. Embedding the
// interface keeps the fake to the method under test; any other call panics,
// which is what a surprise call deserves.
type fakeDocker struct {
	client.CommonAPIClient
	container string
	dstPath   string
	content   []byte
	err       error
}

func (f *fakeDocker) CopyToContainer(_ context.Context, container, dstPath string, content io.Reader, _ types.CopyToContainerOptions) error {
	f.container = container
	f.dstPath = dstPath

	read, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	f.content = read

	return f.err
}

// The archive holds paths relative to the home directory, so it has to be
// extracted there: anywhere else and code-server never reads the settings.
func TestCopyUserSettingsExtractsTheArchiveInTheContainerHome(t *testing.T) {
	t.Parallel()

	docker := &fakeDocker{}

	if err := vscode.CopyUserSettings(context.Background(), docker, "abc123", 501, 20); err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if docker.container != "abc123" {
		t.Errorf("got container %q, want %q", docker.container, "abc123")
	}
	if docker.dstPath != "/home/coder" {
		t.Errorf("got destination %q, want %q", docker.dstPath, "/home/coder")
	}

	var settings map[string]any
	if err := json.Unmarshal(readArchivedFile(t, docker.content, vscode.UserSettingsPath), &settings); err != nil {
		t.Fatalf("the copied archive holds no readable settings: %v", err)
	}
	if got := settings["zenMode.restore"]; got != false {
		t.Errorf("got zenMode.restore = %v, want false", got)
	}
}

// A container that starts without the settings is the bug all over again, so
// the failure has to reach the slide rather than be logged and forgotten.
func TestCopyUserSettingsReportsADockerFailure(t *testing.T) {
	t.Parallel()

	docker := &fakeDocker{err: errors.New("no such container")}

	err := vscode.CopyUserSettings(context.Background(), docker, "abc123", 501, 20)
	if err == nil {
		t.Fatal("got no error, want one")
	}
	if !strings.Contains(err.Error(), "no such container") {
		t.Errorf("got %q, want it to carry the docker error", err)
	}
}
