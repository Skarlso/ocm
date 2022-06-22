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

package core

import (
	"encoding/json"

	"github.com/modern-go/reflect2"

	"github.com/open-component-model/ocm/pkg/runtime"
)

type AccessSpecRef struct {
	generic   *GenericAccessSpec
	evaluated AccessSpec
}

var _ AccessSpec = (*AccessSpecRef)(nil)

func NewAccessSpecRef(spec AccessSpec) *AccessSpecRef {
	if reflect2.IsNil(spec) {
		return nil
	}
	if r, ok := spec.(*AccessSpecRef); ok {
		return &AccessSpecRef{generic: r.generic, evaluated: r.evaluated}
	}
	r := &AccessSpecRef{}
	r.Set(spec)
	return r
}

func NewRawAccessSpecRef(data []byte, unmarshaler runtime.Unmarshaler) (*AccessSpecRef, error) {
	var spec GenericAccessSpec

	if unmarshaler == nil {
		unmarshaler = runtime.DefaultYAMLEncoding
	}
	err := unmarshaler.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}
	return &AccessSpecRef{generic: &spec}, nil
}

func (a *AccessSpecRef) Set(spec AccessSpec) {
	if g, ok := spec.(*GenericAccessSpec); ok {
		a.evaluated = nil
		a.generic = g
	} else {
		a.evaluated = spec
		a.generic = nil
	}
}

func (a *AccessSpecRef) GetType() string {
	if a.evaluated != nil {
		return a.evaluated.GetType()
	}
	return a.generic.GetType()
}

func (a *AccessSpecRef) GetKind() string {
	if a.evaluated != nil {
		return a.evaluated.GetKind()
	}
	return a.generic.GetKind()
}

func (a *AccessSpecRef) GetVersion() string {
	if a.evaluated != nil {
		return a.evaluated.GetVersion()
	}
	return a.generic.GetVersion()
}

func (a *AccessSpecRef) IsLocal(ctx Context) bool {
	a.assure(ctx)
	if a.evaluated != nil {
		return a.evaluated.IsLocal(ctx)
	}
	return false
}

func (a *AccessSpecRef) AccessMethod(access ComponentVersionAccess) (AccessMethod, error) {
	if err := a.assure(access.GetContext()); err != nil {
		return nil, err
	}
	return a.evaluated.AccessMethod(access)
}

func (a *AccessSpecRef) Evaluate(ctx Context) (AccessSpec, error) {
	err := a.assure(ctx)
	if err != nil {
		return nil, err
	}
	return a.evaluated, nil
}

func (a *AccessSpecRef) assure(ctx Context) error {
	var err error
	if a.evaluated == nil {
		a.evaluated, err = a.generic.Evaluate(ctx)
	}
	return err
}

// UnmarshalJSON implements a custom json unmarshal method for an access spec ref
func (a *AccessSpecRef) UnmarshalJSON(data []byte) error {
	a.evaluated = nil
	a.generic = &GenericAccessSpec{}
	return json.Unmarshal(data, a.generic)
}

// MarshalJSON implements a custom json unmarshal method for a unstructured type.
func (a *AccessSpecRef) MarshalJSON() ([]byte, error) {
	if a.evaluated == nil {
		return json.Marshal(a.generic)
	}
	return json.Marshal(a.evaluated)
}
