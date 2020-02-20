# Docker volume plugin for ISO's

This plugin will allow you to mount an iso and present it as a volume to
a docker container. Since root is required for mounting best to let docker do
it for you?! Also it's giving me a chance to make something with
[Go](https://golang.org).

## Caveats

The ISO and volume are read-only.
