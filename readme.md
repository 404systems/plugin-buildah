# plugin-buildah
[![status-badge](https://woodpecker.404systems.net/api/badges/3/status.svg?events=push%2Ctag%2Cmanual)](https://woodpecker.404systems.net/repos/3)

Run buildah as a woodpecker plugin.

## Development
You can test the image by creating an env file with your settings.
```
PLUGIN_REGISTRY=ghcr.io
PLUGIN_REPO=404systems/plugin-buildah
PLUGIN_USERNAME=foo
PLUGIN_PASSWORD=bar
PLUGIN_TAGS=dev
PLUGIN_MIRROR=mirror.woodpecker.svc.usmn1.internal:80
PLUGIN_CACHE_REPO=404systems/plugin-buildah-cache
PLUGIN_SKIP_PUSH=true
PLUGIN_SKIP_BUILD=false
BUILDAH_ISOLATION=chroot
```
nu
```nu
podman build . -t plugin-buildah ; podman run --rm -it -v ./:/workspace:ro -v /dev/fuse:/dev/fuse --env-file env plugin-buildah
```
bash
```sh
podman build . -t plugin-buildah && podman run --rm -it -v ./:/workspace:ro -v /dev/fuse:/dev/fuse --env-file env plugin-buildah
```

My goal is to dogfood the build process, using plugin-buildah to build plugin-builah. As a result, the build process happens on my private server via Woodpecker and not Github Actions.
This acts as a good test for whether or not features work correctly.  

## Contributing
Contributions are welcome as long as they are made in good faith and avoid using AI/LLMs/Agents to _directly_ write code.

Some features that I'm looking to implement next are:
- Configuring multiple registry logins
- Pushing to multiple registries/destinations
- Outputting the built image to a file instead of pushing it
