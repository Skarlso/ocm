// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"github.com/gardener/ocm/pkg/common"
	"github.com/gardener/ocm/pkg/oci/artdesc"
)

type Repository interface {
	ExistsArtefact(name string, version string) (bool, error)
	LookupArtefact(name string, version string) (ArtefactAccess, error)
	ComposeArtefact(name string, version string) (ArtefactComposer, error)
	WriteArtefact(ArtefactAccess) (ArtefactAccess, error)
}

type BlobAccess = common.BlobAccess
type DataAccess = common.DataAccess

type ArtefactAccess interface {
	GetRepository() Repository

	GetArtefactDescriptor() artdesc.ArtefactDescriptor
	GetManifest(digest string) ManifestAccess
	GetBlob(digest string) BlobAccess
}

////////////////////////////////////////////////////////////////////////////////
// technical abstraction

type ManifestAccess interface {
	GetManifest() artdesc.Manifest
	GetBlob(digest string) BlobAccess
}

type ArtefactComposer interface {
	ArtefactAccess

	AddManifest(artdesc.Manifest) (digest string, err error)
	AddBlob(BlobAccess) (digest string, err error)
	Update() error
}

////////////////////////////////////////////////////////////////////////////////
// logical abstraction

/*
type ManifestAccess interface {
	GetManifestMeta() *ManifestMeta  // for index members
	GetManifest() artdesc.Manifest
	GetBlob(digest string) BlobAccess
}

type ArtefactComposer interface {
	ArtefactAccess

	AsManifest(bool) ManifestComposer
	AsIndex(keepManifest bool) ManifestComposer
}

type LayerMeta struct {
}

type ManifestMeta struct {
}

type IndexComposer interface {
	GetIndex() artdesc.Index
	GetBlob(digest string) BlobAccess
	GetManifest(digest string) artdesc.Manifest
	ComposeManifest(*ManifestMeta) (ManifestComposer, error)
	AddManifest(ArtefactAccess) error
	Update() error
}

type ManifestComposer interface {
	ManifestAccess
	AddLayer(*LayerMeta,BlobAccess) error
	Update() error
}
*/
