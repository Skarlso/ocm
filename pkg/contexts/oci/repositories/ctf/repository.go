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

package ctf

import (
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/datacontext/vfsattr"
	"github.com/open-component-model/ocm/pkg/contexts/oci/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/artefactset"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/ctf/index"
)

/*
   A common transport archive is just a folder with artefact archives.
   in tar format and an index.json file. The name of the archive
   is the digest of the artefact descriptor.

   The artefact archive is a filesystem structure with a file
   artefact-descriptor.json and a folder blobs containing
   the flat blob files with the name according to the blob digest.

   Digests used as filename will replace the ":" by a "."
*/

// Repository is a closable view on a repository implementation
type Repository struct {
	view accessio.CloserView
	*RepositoryImpl
}

func (r *Repository) IsClosed() bool {
	return r.view.IsClosed()
}

func (r *Repository) Close() error {
	return r.view.Close()
}

func (r *Repository) LookupArtefact(name string, ref string) (cpi.ArtefactAccess, error) {
	return r.RepositoryImpl.LookupArtefact(r, name, ref)
}

////////////////////////////////////////////////////////////////////////////////

// RepositoryImpl is closed, if all views are released
type RepositoryImpl struct {
	refs accessio.ReferencableCloser

	ctx  cpi.Context
	spec *RepositorySpec
	base *artefactset.FileSystemBlobAccess
}

var _ cpi.Repository = (*Repository)(nil)

// New returns a new representation based repository
func New(ctx cpi.Context, spec *RepositorySpec, setup accessobj.Setup, closer accessobj.Closer, mode vfs.FileMode) (*Repository, error) {
	if spec.PathFileSystem == nil {
		spec.PathFileSystem = vfsattr.Get(ctx)
	}
	base, err := accessobj.NewAccessObject(accessObjectInfo, spec.AccessMode, spec.Options.Representation, setup, closer, mode)
	return _Wrap(ctx, spec, base, err)
}

func _Wrap(ctx cpi.Context, spec *RepositorySpec, obj *accessobj.AccessObject, err error) (*Repository, error) {
	if err != nil {
		return nil, err
	}
	r := &RepositoryImpl{
		ctx:  ctx,
		spec: spec,
		base: artefactset.NewFileSystemBlobAccess(obj),
	}
	r.refs = accessio.NewRefCloser(r, true)
	return r.View(true)
}

func (r *RepositoryImpl) View(main ...bool) (*Repository, error) {
	v, err := r.refs.View(main...)
	if err != nil {
		return nil, err
	}
	return &Repository{view: v, RepositoryImpl: r}, nil
}

func (r *RepositoryImpl) GetSpecification() cpi.RepositorySpec {
	return r.spec
}

func (r *RepositoryImpl) NamespaceLister() cpi.NamespaceLister {
	return r
}

func (r *RepositoryImpl) NumNamespaces(prefix string) (int, error) {
	return len(cpi.FilterByNamespacePrefix(prefix, r.getIndex().RepositoryList())), nil
}

func (r *RepositoryImpl) GetNamespaces(prefix string, closure bool) ([]string, error) {
	return cpi.FilterChildren(closure, cpi.FilterByNamespacePrefix(prefix, r.getIndex().RepositoryList())), nil
}

////////////////////////////////////////////////////////////////////////////////
// forward

func (r *RepositoryImpl) IsReadOnly() bool {
	return r.base.IsReadOnly()
}

func (r *RepositoryImpl) IsClosed() bool {
	return r.base.IsClosed()
}

func (r *RepositoryImpl) Write(path string, mode vfs.FileMode, opts ...accessio.Option) error {
	return r.base.Write(path, mode, opts...)
}

func (r *RepositoryImpl) Update() error {
	return r.base.Update()
}

func (r *RepositoryImpl) Close() error {
	return r.base.Close()
}

func (a *RepositoryImpl) getIndex() *index.RepositoryIndex {
	if a.IsReadOnly() {
		return a.base.GetState().GetOriginalState().(*index.RepositoryIndex)
	}
	return a.base.GetState().GetState().(*index.RepositoryIndex)
}

////////////////////////////////////////////////////////////////////////////////
// cpi.Repository methods

func (r *RepositoryImpl) ExistsArtefact(name string, tag string) (bool, error) {
	return r.getIndex().HasArtefact(name, tag), nil
}

func (r *RepositoryImpl) LookupArtefact(repo *Repository, name string, ref string) (cpi.ArtefactAccess, error) {
	a := r.getIndex().GetArtefactInfo(name, ref)
	if a == nil {
		return nil, cpi.ErrUnknownArtefact(name, ref)
	}

	ns, err := newNamespace(repo, name) // share repo view.namespace not exposed
	if err != nil {
		return nil, err
	}
	return ns.GetArtefact(ref)
}

func (r *RepositoryImpl) LookupNamespace(name string) (cpi.NamespaceAccess, error) {
	repo, err := r.View() // create new closable view
	if err != nil {
		return nil, err
	}
	return newNamespace(repo, name)
}
