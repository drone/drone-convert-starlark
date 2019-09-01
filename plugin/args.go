// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"github.com/drone/drone-go/drone"
	"go.starlark.net/starlark"
)

// TODO(bradrydzewski) add repository id
// TODO(bradrydzewski) add repository timeout
// TODO(bradrydzewski) add repository counter
// TODO(bradrydzewski) add repository created
// TODO(bradrydzewski) add repository updated
// TODO(bradrydzewski) add repository synced
// TODO(bradrydzewski) add repository version

// TODO(bradrydzewski) add build id
// TODO(bradrydzewski) add build number
// TODO(bradrydzewski) add build parent
// TODO(bradrydzewski) add build timestamp
// TODO(bradrydzewski) add build params
// TODO(bradrydzewski) add build started
// TODO(bradrydzewski) add build finished
// TODO(bradrydzewski) add build created
// TODO(bradrydzewski) add build updated
// TODO(bradrydzewski) add build version
// TODO(bradrydzewski) add build created

func createArgs(repo drone.Repo, build drone.Build) []starlark.Value {
	dict := new(starlark.Dict)
	dict.SetKey(starlark.String("build"), fromBuild(build))
	dict.SetKey(starlark.String("repo"), fromRepo(repo))
	return starlark.Tuple([]starlark.Value{dict})
}

func fromBuild(v drone.Build) *starlark.Dict {
	dict := new(starlark.Dict)
	dict.SetKey(starlark.String("event"), starlark.String(v.Event))
	dict.SetKey(starlark.String("action"), starlark.String(v.Action))
	dict.SetKey(starlark.String("cron"), starlark.String(v.Cron))
	dict.SetKey(starlark.String("environent"), starlark.String(v.Deploy))
	dict.SetKey(starlark.String("link"), starlark.String(v.Link))
	dict.SetKey(starlark.String("branch"), starlark.String(v.Target))
	dict.SetKey(starlark.String("source"), starlark.String(v.Source))
	dict.SetKey(starlark.String("before"), starlark.String(v.Before))
	dict.SetKey(starlark.String("after"), starlark.String(v.After))
	dict.SetKey(starlark.String("target"), starlark.String(v.Target))
	dict.SetKey(starlark.String("ref"), starlark.String(v.Ref))
	dict.SetKey(starlark.String("commit"), starlark.String(v.After))
	dict.SetKey(starlark.String("title"), starlark.String(v.Title))
	dict.SetKey(starlark.String("message"), starlark.String(v.Message))
	dict.SetKey(starlark.String("source_repo"), starlark.String(v.Fork))
	dict.SetKey(starlark.String("author_login"), starlark.String(v.Author))
	dict.SetKey(starlark.String("author_name"), starlark.String(v.AuthorName))
	dict.SetKey(starlark.String("author_email"), starlark.String(v.AuthorEmail))
	dict.SetKey(starlark.String("author_avatar"), starlark.String(v.AuthorAvatar))
	dict.SetKey(starlark.String("sender"), starlark.String(v.Sender))
	return dict
}

func fromRepo(v drone.Repo) *starlark.Dict {
	dict := new(starlark.Dict)
	dict.SetKey(starlark.String("uid"), starlark.String(v.UID))
	dict.SetKey(starlark.String("name"), starlark.String(v.Name))
	dict.SetKey(starlark.String("namespace"), starlark.String(v.Namespace))
	dict.SetKey(starlark.String("slug"), starlark.String(v.Slug))
	dict.SetKey(starlark.String("git_http_url"), starlark.String(v.HTTPURL))
	dict.SetKey(starlark.String("git_ssh_url"), starlark.String(v.SSHURL))
	dict.SetKey(starlark.String("link"), starlark.String(v.Link))
	dict.SetKey(starlark.String("branch"), starlark.String(v.Branch))
	dict.SetKey(starlark.String("config"), starlark.String(v.Config))
	dict.SetKey(starlark.String("private"), starlark.Bool(v.Private))
	dict.SetKey(starlark.String("visibility"), starlark.String(v.Visibility))
	dict.SetKey(starlark.String("active"), starlark.Bool(v.Active))
	dict.SetKey(starlark.String("trusted"), starlark.Bool(v.Trusted))
	dict.SetKey(starlark.String("protected"), starlark.Bool(v.Protected))
	dict.SetKey(starlark.String("ignore_forks"), starlark.Bool(v.IgnoreForks))
	dict.SetKey(starlark.String("ignore_pull_requests"), starlark.Bool(v.IgnorePulls))
	return dict
}
