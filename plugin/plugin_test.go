// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/drone/drone-convert-starlark/plugin/starlark/repo"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// empty context
var noContext = context.Background()

// Helper that creates a repo registry with the test repo added.
func newTestRegistry(t *testing.T) *repo.Registry {
	registry := repo.NewRegistry()
	if err := registry.Add("test", "testdata"); err != nil {
		t.Error("unable to add test registry", err)
	}
	return registry
}

// Simplest case return of a single pipeline with no extension loading.
func TestPlugin_Single_Pipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.yml",
		},
	}

	plugin := New(newTestRegistry(t))

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if config != nil {
		t.Error("Want nil config when configuration is not starlark file")
		return
	}

	before, err := ioutil.ReadFile("testdata/single_pipeline.star")
	if err != nil {
		t.Error(err)
		return
	}
	after, err := ioutil.ReadFile("testdata/single_pipeline.star.golden")
	if err != nil {
		t.Error(err)
		return
	}

	req.Repo.Config = "single_pipeline.star"
	req.Config.Data = string(before)
	config, err = plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if config == nil {
		t.Error("Want non-nil configuration")
		return
	}

	if want, got := config.Data, string(after); want != got {
		t.Errorf("Want %q got %q", want, got)
	}
}

// Return a list of pipelines instead of a singleton.
func TestPlugin_Multi_Pipeline(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/multi_pipeline.star")
	if err != nil {
		t.Error(err)
		return
	}
	after, err := ioutil.ReadFile("testdata/multi_pipeline.star.golden")
	if err != nil {
		t.Error(err)
		return
	}

	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.star",
		},
		Config: drone.Config{
			Data: string(before),
		},
	}

	plugin := New(newTestRegistry(t))
	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	config, err = plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if config == nil {
		t.Error("Want non-nil configuration")
		return
	}

	if want, got := config.Data, string(after); want != got {
		t.Errorf("Want %q got %q", want, got)
	}
}

// Load an extension while extension loading is disabled.
func TestPlugin_Loading_Disabled(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/load_valid.star")
	if err != nil {
		t.Error(err)
		return
	}

	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.star",
		},
		Config: drone.Config{
			Data: string(before),
		},
	}

	// Empty registry with no repos added should disable extension loading.
	registry := repo.NewRegistry()
	plugin := New(registry)
	_, err = plugin.Convert(noContext, req)
	if err == nil {
		t.Errorf("Want ErrLoadingDisabled got nil")
	}
}

// Load a valid and known extension. The happy path case.
func TestPlugin_Load_Valid_Extension(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/load_valid.star")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/load_valid.star.golden")
	if err != nil {
		t.Error(err)
		return
	}

	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.star",
		},
		Config: drone.Config{
			Data: string(before),
		},
	}

	plugin := New(newTestRegistry(t))
	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
	}

	if want, got := string(after), config.Data; want != got {
		t.Errorf("Want %q got %q", want, got)
	}
}

// Load an extension that does not exist.
func TestPlugin_Load_Unknown_Extension(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/load_unknown_import.star")
	if err != nil {
		t.Error(err)
		return
	}

	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.star",
		},
		Config: drone.Config{
			Data: string(before),
		},
	}

	plugin := New(newTestRegistry(t))
	_, err = plugin.Convert(noContext, req)
	if err == nil {
		t.Error("expected eval error")
	}
}

// Load an extension from a repo name that hasn't been registered.
func TestPlugin_Load_Unknown_Repo(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/load_unknown_repo.star")
	if err != nil {
		t.Error(err)
		return
	}

	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.star",
		},
		Config: drone.Config{
			Data: string(before),
		},
	}

	plugin := New(newTestRegistry(t))
	_, err = plugin.Convert(noContext, req)
	if err == repo.ErrNoSuchRepo {
		t.Error("expected ErrNoSuchRepo")
	}
}
