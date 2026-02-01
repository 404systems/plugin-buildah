# Overview
A Woodpecker CI plugin to build OCI containers.  

## Features
- Run with privileged enabled or disabled (more details below)
- Automatically generate image tags based on the git tag
- Use an OCI registry to store layer caches

## Privileged and non-privileged execution
### Privileged (default)
By default plugin-buildah runs with `rootless` isolation for RUN commands and uses `fuse-overlayfs` for layer storage.  
Fuse-overlayfs requires access to `/dev/fuse`, but root is not needed. The container runs with user:group `1000:1000`, so no extra options are needed to keep it secure.

```
steps:
- name: build
  ...
  privileged: true
```
Lastly, you need to add `ghcr.io/404systems/plugin-buildah` to the `WOODPECKER_PLUGINS_PRIVILEGED` env variable on your Woodpecker server.

This is the recommended mode for running on kubernetes because in kubernetes the only easy way to allow a container to use `/dev/fuse` is by setting `privileged: true`.

### Non-privileged
To run the container without any system privileges, we need to change a few environment variables: BUILDAH_ISOLATION and STORAGE_DRIVER.

#### BUILDAH_ISOLATION
With buildah there are two isolation modes we care about, `rootless` and `chroot`. By default the image uses `rootless`, but that still needs some system privileges.  
If that is an issue we can use `chroot` instead and omit `privileged: true` in our config.  
However, our container still needs `/dev/fuse`. Under docker or podman that is easy to include.
```
steps:
- name: build
  ...
  environment:
    BUILDAH_ISOLATION: chroot
  volumes:
  - /dev/fuse:/dev/fuse
```
From what I understand, `chroot` is both less secure and less capable than the `rootless` (default for this image), or `oci` options. But since the build is running in a rootless container which isn't allowed any host system privileges, the security impact is negligable. Depending on your hosts OCI runtime, you may be able to run with `rootless` since some runtimes allow non-root users to use some namespacing features.  

#### STORAGE_DRIVER
If are running in a system where `/dev/fuse` isn't available, we can fall back to `vfs`.

```
steps:
- name: build
  ...
  environment:
    BUILDAH_ISOLATION: chroot
    STORAGE_DRIVER: vfs
```
VFS creates a complete copy on each layer operation. With a pretty basic Containerfiles you could end up with many gigabytes of layer data during runtime.  
When running on a host which doesn't have a CoW filesystem this can be a major bottleneck.  
On my dev laptop (BTRFS on NVME using podman) I wasn't able to reasonably tell a difference between VFS and fuse-overlay in terms of build time.  

## Settings

| Setting Name      | Default     | Description                                                                                                                                        |
| ----------------- | ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `repo`            | _none_      | sets the repository to push to                                                                                                                     |
| `registry`        | `docker.io` | sets the registry to push to                                                                                                                       |
| `username`        | _none_      | sets the username to authenticate with                                                                                                             |
| `password`        | _none_      | sets the password to authenticate with                                                                                                             |
| `containerfile`   | _none_      | sets the path to the containerfile, if the file is in the root of the git repo, buildah will find it automatically                                 |
| `context`         | `.`         | sets the build context                                                                                                                             |
| `target`          | _none_      | sets the build target                                                                                                                              |
| `tags`            | _none_      | sets the tags that will be added to the image                                                                                                      |
| `auto_tag`        | _none_      | automatically add the git tag as a tag on the image. if the git tag is a semver then it will also add tags for the major version and minor version |
| `search_registry` | `docker.io` | sets the registry that will be searched when no registry is supplied in an image name                                                              |
| `mirror`          | _none_      | sets the mirror that will be used for the search registry                                                                                          |
| `cache_repo`      | _none_      | sets the repo in the`registry` that will be used for layer caching                                                                                 |
| `skip_push`       | _none_      | enabling this will skip pushing layer cache or result images                                                                                       |
| `skip_build`      | _none_      | enabling this will skip building, but computed buildah arguments will still be shown                                                               |
| `retries`         | `1`         | sets the number of additional attempts to make when pulling an image or layer                                                                      |

## Additional Environment Variables

These environment variables are read by the wrapper, but are not expected to be used by the end user.  
The binary doesn't set any defaults because when running without container isolation we want to avoid accidental file overwrites.  
However, in the Containerfile we set these to be in the /tmp directory.

| Environment Variable   | Default                           | Description                                                                |
| ---------------------- | --------------------------------- | -------------------------------------------------------------------------- |
| `AUTHS_FILE`           | _none_ and `/tmp/auths.json`      | sets the file that will be created/used for storing credentials            |
| `REGISTRIES_FILE`      | _none_ and `/tmp/registries.conf` | sets the file that will be create/used for storing registry configurations |
