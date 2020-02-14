package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"golang.org/x/sys/unix"
)

type isoVolume struct {
	iso string
	Mountpoint string
	containerMountPoint string
}

type isoDriver struct {
	mutex   *sync.Mutex
	root    string
	volumes map[string]*isoVolume
}

func newIsoDriver(volumesRoot string) *isoDriver {
	return &isoDriver{
		root: volumesRoot,
		volumes: map[string]*isoVolume{},
	}
}

func (d isoDriver) Create(r *volume.CreateRequest) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	v := &isoVolume{}
	if iso, ok := r.Options["iso"]; !ok {
		return fmt.Errorf("'iso' path option required")
	} else {
		v.iso = iso
	}

	v.Mountpoint = filepath.Join(d.root, r.Name)
	if err := os.Mkdir(v.Mountpoint, os.FileMode(0755)); err != nil {
		return fmt.Errorf("Could not create mount point %s", v.Mountpoint)
	}

	d.volumes[r.Name] = v
	return nil
}

func (d isoDriver) Remove(r *volume.RemoveRequest) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

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
	d.mutex.Lock()
	defer d.mutex.Unlock()

	v, present := d.volumes[r.Name]
	if !present {
		return &volume.PathResponse{}, fmt.Errorf("Path volume %s not found", r.Name)
	}

	return &volume.PathResponse{Mountpoint: v.Mountpoint}, nil
}

func (d isoDriver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	v, present := d.volumes[r.Name]
	if !present {
		return &volume.MountResponse{}, fmt.Errorf("Mount volume %s not found", r.Name)
	}

	stat, err := os.Lstat(v.Mountpoint)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(v.Mountpoint, 0755); err != nil {
				return &volume.MountResponse{}, fmt.Errorf("Could not create mount point %s", v.Mountpoint)
			}
		} else {
			return &volume.MountResponse{}, fmt.Errorf(err.Error())
		}
	}
	if stat != nil && !stat.IsDir() {
		return &volume.MountResponse{}, fmt.Errorf("%s exists and is not a directory")
	}

	if err := unix.Mount(v.iso, v.Mountpoint, "iso9660", 0, "ro"); err != nil {
		return &volume.MountResponse{}, fmt.Errorf("Could not mount %s to %s", v.iso, v.Mountpoint)
	}

	return &volume.MountResponse{Mountpoint: v.Mountpoint}, nil
}

func (d isoDriver) Unmount(r *volume.UnmountRequest) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

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

	if err := unix.Unmount(v.Mountpoint, 0); err != nil {
		return fmt.Errorf("Could not unmount %s", v.Mountpoint)
	}

	return nil
}

func (d isoDriver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	v, present := d.volumes[r.Name]
	if !present {
		return &volume.GetResponse{}, fmt.Errorf("Mount volume %s not found", r.Name)
	}

	return &volume.GetResponse{Volume: &volume.Volume{Name: r.Name, Mountpoint: v.Mountpoint}}, nil
}

func (d isoDriver) List() (*volume.ListResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var volumes []*volume.Volume
	for name, v := range d.volumes {
		volumes = append(volumes, &volume.Volume{Name: name, Mountpoint: v.Mountpoint})
	}

	return &volume.ListResponse{Volumes: volumes}, nil
}

func (d isoDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "local"}}
}
