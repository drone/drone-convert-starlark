// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"

	"go.starlark.net/starlark"
)

const (
	separator = "---"
	newline   = "\n"
)

var (
	// ErrMainMissing indicates the starlark script is missing
	// the main method.
	ErrMainMissing = errors.New("starlark: missing main function")

	// ErrMainInvalid indicates the starlark script defines a
	// global variable named main, however, it is not callable.
	ErrMainInvalid = errors.New("starlark: main must be a function")

	// ErrMainReturn indicates the starlark script's main method
	// returns an invalid or unexpected type.
	ErrMainReturn = errors.New("starlark: main returns an invalid type")
)

// New returns a new converter plugin.
func New() converter.Plugin {
	return new(plugin)
}

type plugin struct {
}

func (p *plugin) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// if the file is not a Starlark script return no-content.
	// this will instruct the caller to use the configuration
	// file as-is.
	if !isStarlark(req.Repo.Config) {
		return nil, nil
	}

	thread := &starlark.Thread{Name: "drone"}
	globals, err := starlark.ExecFile(thread, req.Repo.Config, []byte(req.Config.Data), nil)
	if err != nil {
		return nil, err
	}

	// find the main method in the starlark script and
	// cast to a callable type. If not callable the script
	// is invalid.
	mainVal, ok := globals["main"]
	if !ok {
		return nil, ErrMainMissing
	}
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, ErrMainInvalid
	}

	// create the input args and invoke the main method
	// using the input args.
	args := createArgs(req.Repo, req.Build)
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch v := mainVal.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			buf.WriteString(separator)
			buf.WriteString(newline)
			if err := write(buf, item); err != nil {
				return nil, err
			}
			buf.WriteString(newline)
		}
	case *starlark.Dict:
		if err := write(buf, v); err != nil {
			return nil, err
		}
	default:
		return nil, ErrMainReturn
	}

	return &drone.Config{
		Data: buf.String(),
	}, nil
}

// helper function returns return if the configuration file
// is in Starlark format.
func isStarlark(name string) bool {
	switch {
	case strings.HasSuffix(name, ".star"):
		return true
	case strings.HasSuffix(name, ".starlark"):
		return true
	default:
		return false
	}
}
