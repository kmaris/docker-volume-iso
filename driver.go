package main

//TODO: Support mounting the same ISO in multiple locations/containers?:
//https://unix.stackexchange.com/questions/520747/mount-iso-o-loop-select-loop-device

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	log "github.com/sirupsen/logrus"
)

type isoVolume struct {
	iso                 string
	Mountpoint          string
	containerMountPoint string
}

type isoDriver struct {
	sync.Mutex

	root    string
	volumes map[string]*isoVolume
}

func newIsoDriver(volumesRoot string) *isoDriver {
	return &isoDriver{
		root:    volumesRoot,
		volumes: map[string]*isoVolume{},
	}
}

// Create the volume mount point on the host
func (d isoDriver) Create(r *volume.CreateRequest) error {
	log.WithField("method", "Create").Debugf("%#v", r)

	d.Lock()
	defer d.Unlock()

	v := &isoVolume{}
	if iso, ok := r.Options["iso"]; !ok {
		return fmt.Errorf("'iso' path option required")
	} else {
		v.iso = iso
	}

	v.Mountpoint = filepath.Join(d.root, r.Name)
	if err := os.Mkdir(v.Mountpoint, os.FileMode(0755)); err != nil {
		return fmt.Errorf("Cannot create mount point %s with permissions 0755: %v", v.Mountpoint, err)
	}

	d.volumes[r.Name] = v
	return nil
}

func (d isoDriver) List() (*volume.ListResponse, error) {
	d.Lock()
	defer d.Unlock()

	var volumes []*volume.Volume
	for name, v := range d.volumes {
		volumes = append(volumes, &volume.Volume{Name: name, Mountpoint: v.Mountpoint})
	}

	return &volume.ListResponse{Volumes: volumes}, nil
}

func (d isoDriver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	d.Lock()
	defer d.Unlock()
	log.WithField("method", "Get").Debugf("%#v", r)

	v, present := d.volumes[r.Name]
	if !present {
		return &volume.GetResponse{}, fmt.Errorf("Mount volume %s not found", r.Name)
	}

	return &volume.GetResponse{Volume: &volume.Volume{Name: r.Name, Mountpoint: v.Mountpoint}}, nil
}

func (d isoDriver) Remove(r *volume.RemoveRequest) error {
	d.Lock()
	defer d.Unlock()
	log.WithField("method", "Remove").Debugf("%#v", r)

	v, present := d.volumes[r.Name]
	if !present {
		return fmt.Errorf("Volume<%s> not found to remove", r.Name)
	}

	if err := os.Remove(v.Mountpoint); err != nil {
		return fmt.Errorf("Could not remove %s", v.Mountpoint)
	}

	delete(d.volumes, r.Name)
	return nil
}

func (d isoDriver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	d.Lock()
	defer d.Unlock()
	log.WithField("method", "Path").Debugf("%#v", r)

	v, present := d.volumes[r.Name]
	if !present {
		return &volume.PathResponse{}, fmt.Errorf("Path volume %s not found", r.Name)
	}

	return &volume.PathResponse{Mountpoint: v.Mountpoint}, nil
}

func (d isoDriver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	d.Lock()
	defer d.Unlock()
	log.WithField("method", "Mount").Debugf("%#v", r)

	v, present := d.volumes[r.Name]
	if !present {
		return &volume.MountResponse{}, fmt.Errorf("Mount volume %s not found", r.Name)
	}

	stat, err := os.Lstat(v.Mountpoint)
	if err != nil && os.IsNotExist(err) {
		return &volume.MountResponse{}, fmt.Errorf("Missing mount point %s: %v", v.Mountpoint, err.Error())
	}

	if stat != nil && !stat.IsDir() {
		log.WithField("v.Mountpoint", v.Mountpoint).Debugf("%#v", stat)
		return &volume.MountResponse{}, fmt.Errorf("Mount point %s exists and is not a directory", v.Mountpoint)
	}

	if err := exec.Command("mount", v.iso, v.Mountpoint).Run(); err != nil {
		return &volume.MountResponse{}, fmt.Errorf("Could not mount %s to %s: %v", v.iso, v.Mountpoint, err.Error())
	}

	return &volume.MountResponse{Mountpoint: v.Mountpoint}, nil
}

func (d isoDriver) Unmount(r *volume.UnmountRequest) error {
	d.Lock()
	defer d.Unlock()
	log.WithField("method", "Unmount").Debugf("%#v", r)

	v, present := d.volumes[r.Name]
	if !present {
		return fmt.Errorf("Mount volume %s not found", r.Name)
	}

	_, err := os.Lstat(v.Mountpoint)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Mountpoint %s does not exist", v.Mountpoint)
		} else {
			if err != nil {
				return fmt.Errorf(err.Error())
			}
		}
	}

	if err := exec.Command("unmount", v.Mountpoint); err != nil {
		return fmt.Errorf("Could not unmount %s", v.Mountpoint)
	}

	return nil
}

func (d isoDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "local"}}
}
