// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package accessobj

import (
	"io"
	"os"
	"sync"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/opencontainers/go-digest"
)

type FileSystemBlobAccess struct {
	sync.RWMutex
	base *AccessObject
}

func NewFileSystemBlobAccess(access *AccessObject) *FileSystemBlobAccess {
	return &FileSystemBlobAccess{
		base: access,
	}
}

func (a *FileSystemBlobAccess) Access() *AccessObject {
	return a.base
}

func (a *FileSystemBlobAccess) IsReadOnly() bool {
	return a.base.IsReadOnly()
}

func (a *FileSystemBlobAccess) IsClosed() bool {
	return a.base.IsClosed()
}

func (a *FileSystemBlobAccess) Write(path string, mode vfs.FileMode, opts ...accessio.Option) error {
	return a.base.Write(path, mode, opts...)
}

func (a *FileSystemBlobAccess) Update() error {
	return a.base.Update()
}

func (a *FileSystemBlobAccess) Close() error {
	return a.base.Close()
}

func (a *FileSystemBlobAccess) GetState() State {
	return a.base.GetState()
}

// DigestPath returns the path to the blob for a given name.
func (a *FileSystemBlobAccess) DigestPath(digest digest.Digest) string {
	return a.BlobPath(common.DigestToFileName(digest))
}

// BlobPath returns the path to the blob for a given name.
func (a *FileSystemBlobAccess) BlobPath(name string) string {
	return a.base.GetInfo().SubPath(name)
}

func (a *FileSystemBlobAccess) GetBlobData(digest digest.Digest) (int64, accessio.DataAccess, error) {
	if a.IsClosed() {
		return accessio.BLOB_UNKNOWN_SIZE, nil, accessio.ErrClosed
	}
	path := a.DigestPath(digest)
	if ok, err := vfs.FileExists(a.base.GetFileSystem(), path); ok {
		return accessio.BLOB_UNKNOWN_SIZE, accessio.DataAccessForFile(a.base.GetFileSystem(), path), nil
	} else {
		if err != nil {
			return accessio.BLOB_UNKNOWN_SIZE, nil, err
		}
		return accessio.BLOB_UNKNOWN_SIZE, nil, accessio.ErrBlobNotFound(digest)
	}
}

func (a *FileSystemBlobAccess) GetBlobDataByName(name string) (accessio.DataAccess, error) {
	if a.IsClosed() {
		return nil, accessio.ErrClosed
	}
	path := a.BlobPath(name)
	if ok, err := vfs.FileExists(a.base.GetFileSystem(), path); ok {
		return accessio.DataAccessForFile(a.base.GetFileSystem(), path), nil
	} else {
		if err != nil {
			return nil, err
		}
		return nil, accessio.ErrBlobNotFound(digest.Digest(name))
	}
}

func (a *FileSystemBlobAccess) AddBlob(blob accessio.BlobAccess) error {
	if a.base.IsClosed() {
		return accessio.ErrClosed
	}
	if a.base.IsReadOnly() {
		return accessio.ErrReadOnly
	}

	path := a.DigestPath(blob.Digest())
	if ok, err := vfs.FileExists(a.base.GetFileSystem(), path); ok {
		return nil
	} else {
		if err != nil {
			return err
		}
	}
	r, err := blob.Reader()
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := a.base.GetFileSystem().OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, a.base.GetMode()&0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	if err != nil {
		w.Close()
		return err
	}
	return w.Close()
}
