package main

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

const (
	testVolumeName    = "docker-volume-iso-test"
	testContainerName = "docker-volume-iso-test"
)

func clean(ctx context.Context, c *docker.Client) {
	log.WithField("method", "clean").Infof("Cleaning test artifacts")
	cro := types.ContainerRemoveOptions{Force: true}
	c.ContainerKill(ctx, testContainerName, "KILL")
	c.ContainerRemove(ctx, testContainerName, cro)
	c.VolumeRemove(ctx, testVolumeName, true)
}

func TestDriver(test *testing.T) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		test.Fatalf("Could not create docker client %v", err.Error())
	}
	ctx := context.Background()
	defer clean(ctx, cli)

	test.Run("Create", func(t *testing.T) {
		vcb := volume.VolumesCreateBody{
			Driver: "iso",
			DriverOpts: map[string]string{
				"iso": "dsl-4.11.rc2.iso",
			},
			Name: testVolumeName,
		}
		_, err := cli.VolumeCreate(ctx, vcb)
		if err != nil {
			t.Fatalf("Failed to create %v: %v", testVolumeName, err.Error())
		}
	})

	test.Run("Inspect", func(t *testing.T) {
		v, err := cli.VolumeInspect(ctx, testVolumeName)
		if err != nil {
			t.Fatalf("Failed to inspect %v: %v", v.Name, err.Error())
		}
		if v.Name != testVolumeName {
			t.Fatalf("Incorrect volume name. Want %v; got %v", v.Name, testVolumeName)
		}
	})

	test.Run("RunTestContainer", func(t *testing.T) {
		testVol, err := cli.VolumeInspect(ctx, testVolumeName)
		if err != nil {
			t.Fatalf("Failed to inspect for runnning, dis bad :( %v: %v", testVolumeName, err.Error())
		}
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "busybox",
			Cmd:   []string{"/bin/sh", "-c", "[[ -n \"$(ls /mnt)\" ]] && exit 0 || exit 1"},
			Tty:   true,
		}, &container.HostConfig{
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: testVol.Name,
					Target: "/mnt",
				},
			},
		},
			nil, testContainerName)
		if err != nil {
			t.Fatalf("Failed to create %v: %v", testContainerName, err.Error())
		}
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			t.Fatalf("Failed to start %v: %v", testContainerName, err.Error())
		}
		exitCode, err := cli.ContainerWait(ctx, resp.ID)
		if err != nil {
			t.Fatalf("Failed waiting for %v: %v", testContainerName, err.Error())
		}
		if exitCode != 0 {
			t.Fatalf("/mnt in %v was emtpy, expected non-empty listing", testContainerName)
		}
	})
}
