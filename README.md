# drone-convert-starlark

A conversion plugin that provides optional support for generating pipeline configuration files via Starlark scripting. _Please note this project requires Drone server version 1.4 or higher._

Questions and Support:<br/>
https://discourse.drone.io

Bug Tracker:<br/>
https://discourse.drone.io/c/bugs

Features and Roadmap:<br/>
https://trello.com/b/ttae5E5o/drone

## Installation

Create a shared secret:

```text
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```text
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --restart=always \
  --name=starlark drone/drone-convert-starlark
```

Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_CONVERT_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
```

## External loads/imports (optional)

Starlark/Bazel has support for [loading](https://docs.bazel.build/versions/master/build-ref.html#load) extensions (modules). This is useful for cases where you'd like to share re-usable logic. drone-convert-starlark supports the ability to define extension repositories that pipelines can load extensions from.

For example:

```python
# Relative load.
load("@test//:steps.star", "example_step")

# Absolute load.
load("@test//subpackage:pipelines.star", "example_pipeline")

def main(ctx):
    return example_pipeline("sample", steps = example_step())
``` 

In this case, `test` is the repo name. The first `load` imports an extension named `steps_extension.star` and extracts the `example_step` symbol for use in our Drone pipeline. The second example drills into the `subpackage` directory within the `test` repo to load an extension called `pipelines.star`, then extracts a symbol named `example_pipeline`.

To make this work for your own pipelines, you'll need to set an `DRONE_STARLARK_REPO_PATHS` environment variable, with the value being a comma-separated list of `repoName:repoPath` pairs. For example:

```text
DRONE_STARLARK_REPO_PATHS=test:/var/lib/drone/starlark/test,othertest:/var/lib/drone/starlark/othertest
```

After setting the above, the converter will output the defined repos during startup.

_Note: This will only allow pipelines to load from the remote repos you define. A build will not be able to `load()` files in the same repo as the pipeline (`.drone.star`). 

## Testing

You can test the extension using the command line utility. Provide the command line utility with the conversion extension endpoint and secret.

```text
$ export DRONE_CONVERT_ENDPOINT=http://1.2.3.4:3000
$ export DRONE_CONVERT_SECRET=bea26a2221fd8090ea38720fc445eca6
```

Use the command line utility to convert the Starlark script:

```
$ drone plugins convert path/to/.drone.star
```
