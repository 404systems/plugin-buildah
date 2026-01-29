# plugin-buildah

Run buildah as a woodpecker plugin.
Supports privileged and non-privileged execution.

## Settings


| Setting Name      | Default     | Description                                                                                                                                        |
| ------------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
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

These variables are not expected to be used in Woodpecker. They exist for local development and debugging purposes.


| Environment Variable | Default                | Description                                                                |
| ---------------------- | ------------------------ | ---------------------------------------------------------------------------- |
| `AUTHS_FILE`         | `/tmp/auths.json`      | sets the file that will be created/used for storing credentials            |
| `REGISTRIES_FILE`    | `/tmp/registries.conf` | sets the file that will be create/used for storing registry configurations |

## Privileged

This is the default mode. By default buildah will use rootless isolation, so you should run with a non-root user.

Testing under podman

```env
PLUGIN_REGISTRY=harbor.404systems.net
PLUGIN_REPO=library/plugin-buildah
PLUGIN_USERNAME=david
PLUGIN_PASSWORD=divad
PLUGIN_TAGS=latest
PLUGIN_MIRROR=mirror.woodpecker.svc.usmn1.internal:80
PLUGIN_CACHE_REPO=library/plugin-buildah
PLUGIN_SKIP_PUSH=true
PLUGIN_SKIP_BUILD=false
```

```bash
podman run --rm -it -v ./:/workspace:ro --env-file env --privileged --user 1000:1000 plugin-buildah
```

## Non-privileged

Add these environment variables to enable non-privileged execution.

```bash
BUILDAH_ISOLATION=chroot
STORAGE_DRIVER=vfs
```
