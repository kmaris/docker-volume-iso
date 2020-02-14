package main

import (
	"path/filepath"
	"strconv"
	"os/user"

	"github.com/docker/go-plugins-helpers/volume"
)

func main() {
	d := newIsoDriver(filepath.Join(volume.DefaultDockerRootDirectory, "iso"))
	h := volume.NewHandler(d)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
	h.ServeUnix("iso", gid)
}
