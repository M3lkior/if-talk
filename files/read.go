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

package files

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

var Root = "."

// Read reads a file in .demoit folder.
func Read(path ...string) ([]byte, error) {
	return os.ReadFile(fullpath(path...))
}

// Exists tests if a file exists.
func Exists(path ...string) bool {
	_, err := os.Stat(fullpath(path...))
	return err == nil
}

// Sha256 returns the sha256 digest of a file.
func Sha256(path string) (string, error) {
	file, err := os.Open(fullpath(".demoit", path))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func fullpath(path ...string) string {
	return filepath.Join(Root, filepath.Join(path...))
}
