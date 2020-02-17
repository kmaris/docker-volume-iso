package main

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
)

func TestDriver(test *testing.T) {
	testVolumeName := "docker-volume-iso-test"
	dc, err := docker.NewEnvClient()
	if err != nil {
		test.Fatalf("Could not create docker client %v", err.Error())
	}

	test.Run("Create", func(*testing.T) {
		vcb := volume.VolumeCreateBody{
			Driver:     "docker-volume-iso",
			DriverOpts: map[string]string{"iso": "dsl-4.4.10-initrd.iso"},
			Name:       testVolumeName,
		}
		v, err := dc.VolumeCreate(context.Background(), vcb)

	})

}
