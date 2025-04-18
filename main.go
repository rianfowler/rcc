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
	docker *dagger.Socket,
) (int, error) {
	ctr := dag.Container().
		From("buildpacksio/pack:latest").
		WithDirectory("/app", source).
		WithExec([]string{
			"pack", "build", "demo-node-app",
			"--path", "/app",
			"--builder", "paketobuildpacks/builder-jammy-base",
			"--buildpack", "paketo-buildpacks/nodejs",
			"--env", "BP_DISABLE_SBOM=true",
			// any other --env flags here too
		})

	code, err := ctr.ExitCode(ctx)
	if err != nil {
		return 0, err
	}
	return code, nil
}
