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

package common

import (
	"reflect"
	"strings"

	"github.com/gardener/ocm/pkg/ocm/core"
	"github.com/gardener/ocm/pkg/ocm/runtime"
)

type AccessType struct {
	runtime.JSONTypedObjectCodecBase
	spectype reflect.Type
	name     string
}

func NewAccessType(name string, proto core.AccessSpec) core.AccessType {
	t := reflect.TypeOf(proto)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return &AccessType{
		spectype: t,
		name:     name,
	}
}

func (t *AccessType) GetSpecType() reflect.Type {
	return t.spectype
}

func (t *AccessType) GetName() string {
	i := strings.LastIndex(t.name, "/")
	if i < 0 {
		return t.name
	}
	return t.name[:i]
}

func (t *AccessType) GetVersion() string {
	i := strings.LastIndex(t.name, "/")
	if i < 0 {
		return "v1"
	}
	return t.name[i+1:]
}

func (t *AccessType) Decode(data []byte) (runtime.TypedObject, error) {
	obj := reflect.New(t.spectype)
	return runtime.UnmarshalInto(data, obj.Interface().(runtime.TypedObject))
}
