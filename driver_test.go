package main

import (
	"context"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
)

const (
	testDir = "testdata"
)

func TestMain(m *testing.M) {
	//TODO: Set the log level to match the docker daemon?
	log.SetLevel(log.DebugLevel)
	if err := os.MkdirAll(testDir, 0777); err != nil {
		log.Fatalf("Could not create testdata dir %v", testDir)
	}
	m.Run()
}

func TestDriver(test *testing.T) {
	testVolumeName := "wat"
	dc, err := docker.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create docker client %v", err.Error())
	}

	test.Run("Create", func(*testing.T) {
		vcb := volume.VolumeCreateBody{
			Driver:     "docker-volume-iso",
			DriverOpts: map[string]string{"iso": "dsl-4.4.10-initrd.iso"},
			Name:       testVolumeName,
		}
		v, err := dc.VolumeCreate(context.Background(), vcb)
		// Now check things....
	})

}
