// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// empty context
var noContext = context.Background()

// mock github token
const mockToken = "d7c559e677ebc489d4e0193c8b97a12e"

func TestPlugin(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.yml",
		},
	}

	plugin := New()

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if config != nil {
		t.Error("Want nil config when configuration is not startlark file")
		return
	}

	before, err := ioutil.ReadFile("testdata/single.star")
	if err != nil {
		t.Error(err)
		return
	}
	after, err := ioutil.ReadFile("testdata/single.star.golden")
	if err != nil {
		t.Error(err)
		return
	}

	req.Repo.Config = "single.star"
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

func TestPlugin_Multi(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/multi.star")
	if err != nil {
		t.Error(err)
		return
	}
	after, err := ioutil.ReadFile("testdata/multi.star.golden")
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

	plugin := New()
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
