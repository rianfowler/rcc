// A generated module for Rcc functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/rcc/internal/dagger"
)

type Rcc struct{}

func (m *Rcc) BuildNodeApp(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	// First, prepare the source code in a container
	container := dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/workspace", source).
		WithWorkdir("/workspace")

	// Export the prepared source to a temporary directory
	_, err := container.Directory("/workspace").Export(ctx, "prepared-source")
	if err != nil {
		return "", err
	}

	// Now run pack build on the host machine using a shell command
	_, err = dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/source", dag.Directory()).
		WithExec([]string{
			"sh", "-c",
			"cd /source && pack build demo-node-app --path . --builder paketobuildpacks/builder-jammy-base --env BP_NODE_VERSION=18 --env BP_DISABLE_SBOM=true",
		}).Stdout(ctx)

	if err != nil {
		return "", err
	}

	return "Image built successfully as demo-node-app", nil
}

// func (m *Rcc) BuildNodeApp(
// 	ctx context.Context,
// 	source *dagger.Directory,
// 	// docker *dagger.Socket,
// ) (string, error) {
// 	// apply a 2‑minute timeout
// 	// ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
// 	// defer cancel()

// 	// reference your local source

// 	// // buildpacks lifecycle command
// 	// creatorArgs := []string{
// 	// 	"/cnb/lifecycle/creator",
// 	// 	// "-app=/tmp/src1",
// 	// 	//	"-buildpacks=paketo-buildpacks/nodejs",
// 	// 	"ttl.sh/demo-node-app:30m",
// 	// }

// 	out, err := dag.Container().
// 		From("paketobuildpacks/builder-jammy-full:latest").

// 		// 1) run as root so we can write everywhere
// 		WithUser("root").

// 		// 2) give ourselves a clean, container‑owned /tmp
// 		WithMountedTemp("/tmp").

// 		// 3) mount your source at /src (outside of /tmp)
// 		WithDirectory("/src", source).

// 		// 4) copy into /tmp/src1 (inside your tmpfs)
// 		WithExec([]string{"mkdir", "-p", "/tmp/src1"}).
// 		WithExec([]string{"cp", "-r", "/src/.", "/tmp/src1"}).

// 		// 5) switch into that staging area
// 		WithWorkdir("/tmp/src1").

// 		// 6) bump CNB API so creator respects creation-time, labels, etc.
// 		WithEnvVariable("CNB_PLATFORM_API", "0.14").

// 		// 7) invoke the lifecycle creator directly
// 		WithExec([]string{
// 			"/cnb/lifecycle/creator",
// 			"-app=/tmp/src1",
// 			//	"-buildpack=paketo-buildpacks/nodejs",
// 			"ttl.sh/demo-node-app:30m",
// 		}).

// 		// run it and grab all its logs
// 		Stdout(ctx)

// 	if err != nil {
// 		return "", err
// 	}
// 	return out, nil
// }

// func (m *Rcc) BuildNodeApp(
// 	ctx context.Context,
// 	source *dagger.Directory,
// 	// docker *dagger.Socket,
// ) (string, error) {

// 	// ────────────────────────────────────────────────────────────────────────────────
// 	// 1) Spin up a real Docker daemon as a long‑running service
// 	dind := dag.Container().
// 		From("docker:dind").
// 		// Docker‐in‐Docker needs privileged mode
// 		// WithPrivileged(true).
// 		// Expose the Docker API port
// 		.with
// 		WithExposedPort(2375).
// 		AsService(dagger.ContainerAsServiceOpts{})

// 	dind.Start(ctx)

// 	// 2) Spin up your local registry as another service
// 	reg := dag.Container().
// 		From("registry:2").
// 		WithExposedPort(5000).
// 		AsService(dagger.ContainerAsServiceOpts{})

// 	reg.Start(ctx)

// 	// ────────────────────────────────────────────────────────────────────────────────
// 	// 3) Now run pack, wiring up BOTH services over TCP:
// 	ctr := dag.Container().
// 		From("buildpacksio/pack:latest").
// 		// bind hostname "docker" → our dind
// 		WithServiceBinding("docker", dind).
// 		// bind hostname "registry" → our registry
// 		WithServiceBinding("registry", reg).
// 		// tell Pack to use TCP instead of the socket
// 		WithEnvVariable("DOCKER_HOST", "tcp://docker:2375").
// 		WithDirectory("/app", source).
// 		WithExec([]string{
// 			"pack", "build", "demo-node-app",
// 			"--path", "/app",
// 			"--builder", "paketobuildpacks/builder-jammy-base",
// 			"--buildpack", "paketo-buildpacks/nodejs",
// 			"--env", "BP_DISABLE_SBOM=true",
// 			"--creation-time", "now",
// 			"--publish",
// 			"--tag", "registry:5000/demo-node-app:latest",
// 		})

// 	// // 1) Start an ephemeral registry service on port 5000
// 	// registry := dag.Container().
// 	// 	From("registry:2").
// 	// 	// Expose the registry's HTTP port
// 	// 	WithExposedPort(5000).
// 	// 	// AsService makes this a long‑running service that Dagger will start & health‑check
// 	// 	AsService(dagger.ContainerAsServiceOpts{
// 	// 		Args: []string{"registry", "serve", "/etc/docker/registry/config.yml"},
// 	// 	})

// 	// registry.Start(ctx)
// 	// ctr := dag.Container().
// 	// 	From("buildpacksio/pack:latest").
// 	// 	// WithUnixSocket("/var/run/docker.sock", docker).
// 	// 	WithServiceBinding("registry", registry). // "registry" → our service
// 	// 	WithDirectory("/app", source).
// 	// 	WithEnvVariable("CNB_PLATFORM_API", "0.9").
// 	// 	WithExec([]string{
// 	// 		// call the creator binary directly
// 	// 		"/cnb/lifecycle/creator",
// 	// 		"-app=/app",
// 	// 		"registry:5000/demo-node-app:latest",
// 	// 	})
// 	// WithExec([]string{
// 	// 	"pack", "build", "demo-node-app",
// 	// 	"--path", "/app",
// 	// 	"--builder", "paketobuildpacks/builder-jammy-base",
// 	// 	"--buildpack", "paketo-buildpacks/nodejs",
// 	// 	"--env", "BP_DISABLE_SBOM=true",
// 	// 	"--creation-time", "now",
// 	// 	"--publish",
// 	// 	"--tag", "registry:5000/demo-node-app:latest",
// 	// 	// any other --env flags here too
// 	// })

// 	out, err := ctr.Stdout(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	return out, nil

// }
