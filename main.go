// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/drone/drone-convert-starlark/plugin/starlark/repo"

	"github.com/drone/drone-convert-starlark/plugin"
	"github.com/drone/drone-convert-starlark/server"
	"github.com/drone/drone-go/plugin/converter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind              string            `envconfig:"DRONE_BIND"`
	Debug             bool              `envconfig:"DRONE_DEBUG"`
	Secret            string            `envconfig:"DRONE_SECRET"`
	StarlarkRepoPaths map[string]string `envconfig:"DRONE_STARLARK_REPO_PATHS"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}
	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	repoRegistry := repo.NewRegistry()
	if len(spec.StarlarkRepoPaths) == 0 {
		logrus.Infoln("starlark extension loading disabled")
	} else {
		logrus.Infoln("starlark extension loading enabled")
		logrus.Infoln("defined extension repositories:")
		for repoName, repoPath := range spec.StarlarkRepoPaths {
			if err := repoRegistry.Add(repoName, repoPath); err != nil {
				logrus.Fatalln("unable to add starlark repo:", err)
			}
			logrus.Infof("  @%s = %s\n", repoName, repoPath)
		}
	}

	handler := converter.Handler(
		plugin.New(repoRegistry),
		spec.Secret,
		logrus.StandardLogger(),
	)
	healthzHandler := server.HandleHealthz()

	logrus.Infof("server listening on address %s", spec.Bind)
	http.Handle("/", handler)
	http.Handle("/healthz", healthzHandler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
