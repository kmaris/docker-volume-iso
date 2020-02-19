package main

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
)

const (
	testVolumeName := "docker-volume-iso-test"
)

func clean(dc *Client) {
	dc.VolumeRemote(context.Background(), testVolumeName, true)
}

func TestDriver(test *testing.T) {
	testVolumeName := "docker-volume-iso-test"
	dc, err := docker.NewEnvClient()
	if err != nil {
		test.Fatalf("Could not create docker client %v", err.Error())
	}

	defer clean(dc)
	// Create the volume then do a docker volume ls to see if it exists.
	// Doesn't matter if it's mounted atm.
	test.Run("Create", func(*testing.T) {
		vcb := volume.VolumeCreateBody{
			Driver:     "test",
			DriverOpts: map[string]string{"iso": "dsl-4.4.10-initrd.iso"},
			Name:       testVolumeName,
		}
		v, err := dc.VolumeCreate(context.Background(), vcb)

	})

}
