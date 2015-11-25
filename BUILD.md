# build

The build is orchestrated by `make`. It provides tasks for building statically
linked binaries using the `go` toolchain, docker images, and gzipped tarballs.

The resulting artifacts can be uploaded to github releases using the github
api, and pushed to a remote docker registry.

## tasks

### generic tasks

* `%.tar.gz` - builds a gzipped tarball of `%`, a directory, called `%.tar.gz`
* `%-$(GOOS)-$(GOARCH)` - builds a statically linked binary using `go build`

* `docker-build-%` - builds a docker image from a Dockerfile in directory `%`
* `docker-build-root` - builds a docker image from the root Dockerfile

* `docker-push-%` - pushes the docker image resulting from `docker-build-%`
* `docker-push-root` - pushes the root docker image

* `create-gh-release` - create a github release associated with the current tag
* `gh-release-%` - push an artifact to github releases

* `gh-token` - ensures that the env variable `GH_TOKEN` is set
* `tag` - ensures that the checked out revision is a tag
* `clean-repo` - ensures that the repo checkout is clean (no changes)
* `info` - prints project info

### per-project tasks

* `build` - all tasks resulting in artifacts
* `docker` - all tasks resulting in a docker image (dep: `build`)
* `release` - push docker images and upload artifacts (dep: `docker` `build`)
