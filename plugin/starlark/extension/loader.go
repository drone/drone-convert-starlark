// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package extension provides an API for loading Starlark extensions.
package extension

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/drone/drone-convert-starlark/plugin/starlark/repo"

	"github.com/bazelbuild/bazel-gazelle/label"
	"go.starlark.net/starlark"
)

var (
	// ErrNoRepoSpecified indicates that an import was attempted
	// without a repo specified.
	ErrNoRepoSpecified = errors.New("starlark: imports must include a valid repo name")

	// ErrInvalidExtension indicates a load attempt against something
	// that is not a starlark extension.
	ErrInvalidExtension = errors.New("starlark: can't load files that are not starlark extensions")
)

// parseLabelStr converts a bazel label string into a Label instance.
func parseLabelStr(labelStr string) (*label.Label, error) {
	// Breaks a label string into a struct that separates out the
	// repo name, package path, and extension name.
	parsed, err := label.Parse(labelStr)
	if err != nil {
		return nil, err
	}
	// We don't (yet) support loading extensions from within the repo
	// that is being built.
	if parsed.Repo == "" {
		return nil, ErrNoRepoSpecified
	}
	if !IsValidFilename(parsed.Name) {
		return nil, ErrInvalidExtension
	}
	return &parsed, nil
}

// pathFromLabel converts a parsed label into an absolute path.
func pathFromLabel(repoPath string, pLabel *label.Label) string {
	if pLabel.Relative {
		// No package path provided, so assume root level of the repo.
		return path.Join(repoPath, pLabel.Name)
	}
	// Package path provided. Get more specific.
	return path.Join(repoPath, pLabel.Pkg, pLabel.Name)
}

// loadExtension loads a bazel extension from the given path.
func loadFromFile(t *starlark.Thread, extensionPath string) (starlark.StringDict, error) {
	extensionContents, err := ioutil.ReadFile(extensionPath)
	if err != nil {
		return nil, err
	}

	globals, err := starlark.ExecFile(t, extensionPath, extensionContents, nil)
	if err != nil {
		return nil, err
	}
	return globals, nil
}

// Load loads an extension from the repository registry.
func Load(t *starlark.Thread, registry *repo.Registry, labelStr string) (starlark.StringDict, error) {
	pLabel, err := parseLabelStr(labelStr)
	if err != nil {
		logrus.Debugln("error while parsing extension label:", err)
		return nil, err
	}

	repoPath, err := registry.Path(pLabel.Repo)
	if err != nil {
		logrus.Debugln("invalid registry name during extension load:", err)
		return nil, err
	}

	extensionPath := pathFromLabel(repoPath, pLabel)
	logrus.Debugln("extension path resolved to", extensionPath)
	return loadFromFile(t, extensionPath)
}

// IsValidFilename returns true if the filename appears to be
// a Starlark/Bazel language file.
func IsValidFilename(name string) bool {
	switch {
	case strings.HasSuffix(name, ".star"):
		return true
	case strings.HasSuffix(name, ".starlark"):
		return true
	case strings.HasSuffix(name, ".bzl"):
		return true
	default:
		return false
	}
}
