# Docker volume plugin for ISO's

This plugin will allow you to mount an iso and present it as a volume to
a docker container. Since root is required for mounting best to let docker do
it for you?! Also it's giving me a chance to make something with
[Go](https://golang.org).

#### Caveat

The ISO and volume are read-only.

## Installation

Copy docker-volume-iso to a path on the system. There are systemd related
resources in the systemd/ subdirectory.
