package main

import (
	"github.com/docker/go-plugins-helpers/volume"
)

const socketName = "iso"

func main() {
	d := newIsoDriver("/var/lib/docker/volumes")
	h := volume.NewHandler(d)
	h.ServeUnix(socketName, 0)
}
