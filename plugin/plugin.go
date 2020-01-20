// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/drone/drone-convert-starlark/metrics"

	"github.com/drone/drone-convert-starlark/plugin/starlark/repo"

	"github.com/drone/drone-convert-starlark/plugin/starlark/extension"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
)

const (
	separator = "---"
	newline   = "\n"
)

// limits generated configuration file size.
const limit = 1000000

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

	// ErrMaximumSize indicates the starlark script generated a
	// file that exceeds the maximum allowed file size.
	ErrMaximumSize = errors.New("starlark: maximum file size exceeded")

	// ErrLoadingDisabled indicates the starlark script is attempting to
	// load an external file while extension loading is disabled.
	ErrLoadingDisabled = errors.New("starlark: extension loading is disabled")
)

// NewRegistry returns a new converter plugin.
func New(repoRegistry *repo.Registry) converter.Plugin {
	return &plugin{repoRegistry: repoRegistry}
}

type plugin struct {
	repoRegistry *repo.Registry
}

func (p *plugin) recordSuccess(convertStart time.Time) {
	convertEnd := time.Now()
	convertDuration := convertEnd.Sub(convertStart)
	metrics.ConversionLatency.Observe(convertDuration.Seconds())
	metrics.ConversionSuccessCount.Inc()
}

func (p *plugin) recordFail() {
	metrics.ConversionFailCount.Inc()
}

func (p *plugin) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// if the file is not a Starlark script return no-content.
	// this will instruct the caller to use the configuration
	// file as-is.
	if !extension.IsValidFilename(req.Repo.Config) {
		return nil, nil
	}

	convertStart := time.Now()
	metrics.ConversionCount.Inc()
	droneCfg, err := p.doConversion(ctx, req)
	if err != nil {
		p.recordFail()
		return nil, err
	}

	// Only report timings on success, since timings for failures may be
	// much different (ex: syntax error stopping execution before a load().
	p.recordSuccess(convertStart)
	return droneCfg, nil
}

func (p *plugin) doConversion(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	thread := &starlark.Thread{
		Name: "drone",
		Load: p.loadExtension,
		Print: func(_ *starlark.Thread, msg string) {
			logrus.WithFields(logrus.Fields{
				"namespace": req.Repo.Namespace,
				"name":      req.Repo.Name,
			}).Traceln(msg)
		},
	}

	// This starts the parsing (and possibly loads/imports).
	globals, err := starlark.ExecFile(thread, req.Repo.Config, []byte(req.Config.Data), nil)
	if err != nil {
		return nil, err
	}

	// find the main method in the starlark script and  cast to a callable type.
	// If not callable the script is invalid.
	mainVal, ok := globals["main"]
	if !ok {
		return nil, ErrMainMissing
	}
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, ErrMainInvalid
	}

	// create the input args and invoke the main method using the input args.
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

	// this is a temporary workaround until we
	// implement a LimitWriter.
	if b := buf.Bytes(); len(b) > limit {
		return nil, ErrMaximumSize
	}

	return &drone.Config{
		Data: buf.String(),
	}, nil
}

func (p *plugin) loadExtension(t *starlark.Thread, labelStr string) (starlark.StringDict, error) {
	if p.repoRegistry.Len() == 0 {
		return nil, ErrLoadingDisabled
	}

	logrus.Debugln("attempting to load extension", labelStr)
	loaded, err := extension.Load(t, p.repoRegistry, labelStr)
	if err != nil {
		logrus.Debugln("error while loading extension:", err)
		return nil, err
	}

	logrus.Debugln("successfully loaded extension", labelStr)
	return loaded, err
}
