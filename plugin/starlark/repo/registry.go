// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package repo provides a registry of repo names, mapped to paths.
package repo

import (
	"errors"
	"os"
	"regexp"
)

var (
	// Used for validating repo names.
	labelRepoRegexp = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

	// ErrMalformedRepoName is returned when an invalid repo name is passed.
	ErrMalformedRepoName = errors.New("starlark: malformed repo name")

	// ErrNoSuchRepo is returned when a lookup happens with a syntactically
	// valid repo that can't be found in the registry.
	ErrNoSuchRepo = errors.New("starlark: no such repo")

	// ErrInvalidRepoPath is returned when a repo definition contains
	// an invalid path.
	ErrInvalidRepoPath = errors.New("starlark: invalid repo path")

	// ErrDuplicateAdd is returned when attempting to add a repo that
	// is already in the registry.
	ErrDuplicateAdd = errors.New("starlark: attempted to add a repo that is already registered")
)

// Registry provides an API for registering and resolving repo paths from names.
type Registry struct {
	nameToPath map[string]string
}

// All returns all registered repo names and their respective paths.
func (r *Registry) All() map[string]string {
	return r.nameToPath
}

// Len returns the number of registered repositories.
func (r *Registry) Len() int {
	return len(r.nameToPath)
}

// Add registers a new repository in the registry.
func (r *Registry) Add(repoName string, repoPath string) error {
	if err := validateRepoName(repoName); err != nil {
		return err
	}
	if err := validateRepoPath(repoPath); err != nil {
		return err
	}
	if _, lookupErr := r.Path(repoName); lookupErr == nil {
		return ErrDuplicateAdd
	}
	r.nameToPath[repoName] = repoPath
	return nil
}

// Path returns the filesystem path where the given registry is located.
func (r *Registry) Path(repoName string) (string, error) {
	path, found := r.nameToPath[repoName]
	if !found {
		return "", ErrNoSuchRepo
	}
	return path, nil
}

// NewRegistry instantiates a new registry.
func NewRegistry() *Registry {
	return &Registry{
		nameToPath: make(map[string]string),
	}
}

// validateRepoName ensures that a repo name conforms to the Starlark spec.
func validateRepoName(repoName string) error {
	if !labelRepoRegexp.MatchString(repoName) {
		return ErrMalformedRepoName
	}
	return nil
}

// validateRepoPath ensures that the provided path is a valid repository.
func validateRepoPath(repoPath string) error {
	if repoPath == "" {
		return nil
	}
	fi, err := os.Stat(repoPath)
	if err != nil {
		return err
	}
	if mode := fi.Mode(); mode.IsDir() {
		return nil
	}
	return ErrInvalidRepoPath
}
