package main

import (
	"strconv"
	"os/user"

	"github.com/docker/go-plugins-helpers/volume"
)

const socketName = "iso"

func main() {
	d := newIsoDriver("/var/lib/docker/volumes")
	h := volume.NewHandler(d)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
	h.ServeUnix(socketName, gid)
}
