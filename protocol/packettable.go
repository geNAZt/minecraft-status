/*
   Copyright 2013 Matthew Collins (purggames@gmail.com)
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package protocol

import (
	"reflect"
)

type State int

const (
	Handshaking State = iota
	Play
	Login
	Status
)

type Direction int
type Packet interface{}

const (
	Serverbound Direction = iota
	Clientbound
)

var (
	packets = [4][2][]reflect.Type{
		Handshaking: [2][]reflect.Type{
			Clientbound: []reflect.Type{},
			Serverbound: []reflect.Type{
				reflect.TypeOf((*Handshake)(nil)).Elem(),
			},
		},
		Status: [2][]reflect.Type{
			Clientbound: []reflect.Type{
				reflect.TypeOf((*StatusResponse)(nil)).Elem(),
				reflect.TypeOf((*StatusPing)(nil)).Elem(),
			},
			Serverbound: []reflect.Type{
				reflect.TypeOf((*StatusGet)(nil)).Elem(),
				reflect.TypeOf((*ClientStatusPing)(nil)).Elem(),
			},
		},
	}

	packetsToID = [2]map[reflect.Type]int{
		Clientbound: map[reflect.Type]int{},
		Serverbound: map[reflect.Type]int{},
	}
)

func init() {
	for _, st := range packets {
		for d, dir := range st {
			for i, p := range dir {
				if _, ok := packetsToID[d][p]; ok {
					panic("Duplicate packet " + p.Name())
				}

				packetsToID[d][p] = i
			}
		}
	}
}
