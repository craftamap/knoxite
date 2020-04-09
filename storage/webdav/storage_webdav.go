package webdav

/*
 * knoxite
 *     Copyright (c) 2020, Fabian Siegel <fabians1999@gmail.com>
 *
 *   For license see LICENSE
 */

import (
	"errors"
	"net/url"

	"github.com/studio-b12/gowebdav"

	knoxite "github.com/knoxite/knoxite/lib"
)

//StorageWebDav stores data on a WebDav Server
type StorageWebDav struct {
	URL    url.URL
	Client *gowebdav.Client
	knoxite.StorageFilesystem
}

// Error declaration
var (
	ErrInvalidAuthentication = errors.New("Wrong Username or Password")
)

func init() {
	knoxite.RegisterBackendFactory(&StorageWebDav{})
}

// NewBackend returns a StorageWebDav backend
func (*StorageWebDav) NewBackend(u url.URL) (knoxite.Backend, error) {
	u0, _ := url.Parse(u.String())
	if u0.Scheme == "webdav" {
		u0.Scheme = "http"
	} else if u0.Scheme == "webdavs" {
		u0.Scheme = "https"
	}

	userinfo := u.User
	username := userinfo.Username()
	passwd, _ := userinfo.Password()

	webdavClient := gowebdav.NewClient(u0.String(), username, passwd)
	storage := StorageWebDav{
		URL:    u,
		Client: webdavClient,
	}

	storagedav, err := knoxite.NewStorageFilesystem("", &storage)
	storage.StorageFilesystem = storagedav
	if err != nil {
		return &StorageWebDav{}, err
	}

	return &storage, nil

}

// Location returns the type and location of the repository
func (backend *StorageWebDav) Location() string {
	return backend.URL.String()
}

// Close - We do not need to Close this backend
func (backend *StorageWebDav) Close() error {
	return nil
}

// Protocols returns the Protocol Schemes supported by this backend
func (backend *StorageWebDav) Protocols() []string {
	// Those protocols are not offical protocols, but because webdav uses http, and the
	// http backend already exists, we have to use webdav(s)
	// This protocol scheme is also used by file explorers like dolphin
	return []string{"webdav", "webdavs"}
}

// Description returns a user-friendly description for this backend
func (backend *StorageWebDav) Description() string {
	return "WebDav Storage (Supports {Own/Next}Cloud)"
}

// AvailableSpace is not available (yet?)
func (backend *StorageWebDav) AvailableSpace() (uint64, error) {
	// TODO: This is actually possible, but im leaving it out for now
	return uint64(0), knoxite.ErrAvailableSpaceUnknown
}

// CreatePath creates a path on the remote
func (backend *StorageWebDav) CreatePath(path string) error {
	return backend.Client.MkdirAll(path, 0755)
}

// DeleteFile deletes a remote file
func (backend *StorageWebDav) DeleteFile(path string) error {
	return backend.Client.Remove(path)
}

// DeletePath deletes a directory and its contents
func (backend *StorageWebDav) DeletePath(path string) error {
	return backend.Client.Remove(path)
}

// ReadFile reads the file
func (backend *StorageWebDav) ReadFile(path string) ([]byte, error) {
	return backend.Client.Read(path)
}

// WriteFile writes a file
func (backend *StorageWebDav) WriteFile(path string, data []byte) (size uint64, err error) {
	err = backend.Client.Write(path, data, 0644)
	return uint64(len(data)), err
}

// Stat returns the file size by using the backends Stat function
func (backend *StorageWebDav) Stat(path string) (uint64, error) {
	stat, err := backend.Client.Stat(path)
	if err != nil {
		return 0, err
	}
	return uint64(stat.Size()), nil

}
