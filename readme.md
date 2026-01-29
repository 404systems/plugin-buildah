# plugin-buildah
Run buildah as a woodpecker plugin.  
Supports privileged and non-privileged execution.  

## Privileged
This is the default mode. By default buildah will use rootless isolation, so you should run with a non-root user.  


## Non-privileged
Add these environment variables to enable non-privileged execution.
```bash
BUILDAH_ISOLATION=chroot
STORAGE_DRIVER=vfs
```