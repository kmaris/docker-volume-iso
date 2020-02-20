package main

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
)

const (
	testVolumeName    = "docker-volume-iso-test"
	testContainerName = "docker-volume-iso-test"
)

func clean(dc *Client) {
	dc.VolumeRemove(context.Background(), testVolumeName, true)
}

func TestDriver(test *testing.T) {
	testVolumeName := "docker-volume-iso-test"
	dc, err := docker.NewEnvClient()
	if err != nil {
		test.Fatalf("Could not create docker client %v", err.Error())
	}

	defer clean(dc)
	// Create the volume then do a docker volume ls to see if it exists.
	test.Run("Create", func(*testing.T) {
		vcb := volume.VolumeCreateBody{
			Driver:     "iso",
			DriverOpts: map[string]string{"iso": "dsl-4.11.rc2.iso"},
			Name:       testVolumeName,
		}
		v, err := dc.VolumeCreate(context.Background(), vcb)
		if err != nil {
			test.Fatalf("Failed to create %v: %v", testVolumeName, err.Error())
		}
	})

	test.Run("RunMount", func(*testing.T) {
		ccfg := container.Config{}
		hcfg := container.HostConfig{}
		ncfg := network.NetworkingConfig{}
		container, err := dc.ContainerCreate(context.Background(), ccfg, hcfg, ncfg, testContainerName)
		if err != nil {
			test.Fatalf("Failed to run %v: %v", testContainerName, err.Error())
		}
	})

	//test.Run("Unmount", func(*testing.T) {

}
