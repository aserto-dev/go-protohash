// Copyright 2018 The ObjectHash-Proto Authors
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

package tests

import (
	"testing"

	ph "github.com/aserto-dev/go-protohash"
	"github.com/aserto-dev/go-protohash/tests/api/v1"
	ti "github.com/aserto-dev/go-protohash/tests/internal"
	"google.golang.org/protobuf/proto"
)

// TestOneOfFields checks that oneof fields are handled properly.
func TestOneOfFields(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []ti.TestCase{
		//////////////////////////
		//  Empty oneof fields. //
		//////////////////////////
		{
			Protos: []proto.Message{
				&api.Singleton{},

				&api.Empty{},
			},
			EquivalentJSONString: "{}",
			EquivalentObject:     map[int64]string{},
			ExpectedHashString:   "0",
		},

		/////////////////////////////////////////////
		//  One of the options selected but empty. //
		/////////////////////////////////////////////
		{
			Protos: []proto.Message{
				&api.Singleton{Singleton: &api.Singleton_TheBool{}},

				&api.Singleton{Singleton: &api.Singleton_TheBool{TheBool: false}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject:   map[int64]bool{1: false},
			ExpectedHashString: "fc42f4d44454522d",
		},

		{
			Protos: []proto.Message{
				&api.Singleton{Singleton: &api.Singleton_TheString{}},

				&api.Singleton{Singleton: &api.Singleton_TheString{TheString: ""}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject:   map[int64]string{25: ""},
			ExpectedHashString: "5fd4b748b2f5442c",
		},

		{
			Protos: []proto.Message{
				&api.Singleton{Singleton: &api.Singleton_TheInt32{}},

				&api.Singleton{Singleton: &api.Singleton_TheInt32{TheInt32: 0}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject:   map[int64]int32{13: 0},
			ExpectedHashString: "cb720bf58a2c29ec",
		},

		////////////////////////////////////////////////
		//  One of the options selected with content. //
		////////////////////////////////////////////////
		//
		// For protobufs, it is legal (and backwards-compatible) to update a message by wrapping
		// an existing field within a oneof rule. Therefore, both objects (using old schema and
		// the new schema) should result in the same objecthash.
		//
		// Example:
		//
		// # Old schema:               | # New schema:
		// message Simple {            | message Singleton {
		//   string string_field = 25; |   oneof singleton {
		// }                           |     string the_string = 25;
		//                             |   }
		//                             | }
		//
		// The following examples demonstrate this equivalence.

		{
			Protos: []proto.Message{
				&api.Simple{StringField: "TEST!"},

				&api.Singleton{Singleton: &api.Singleton_TheString{TheString: "TEST!"}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject:   map[int64]string{25: "TEST!"},
			ExpectedHashString: "6e388481a9f4259e",
		},

		{
			Protos: []proto.Message{
				&api.Simple{Int32Field: 99},

				&api.Singleton{Singleton: &api.Singleton_TheInt32{TheInt32: 99}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject:   map[int64]int32{13: 99},
			ExpectedHashString: "7ad488431568d7ef",
		},

		///////////////////////////
		//  Nested oneof fields. //
		///////////////////////////
		{
			Protos: []proto.Message{
				&api.Simple{SingletonField: &api.Singleton{}},

				&api.Singleton{Singleton: &api.Singleton_TheSingleton{TheSingleton: &api.Singleton{}}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject: map[int64]map[int64]int64{35: {}},
			// EquivalentObject:   map[int64]map[int64]map[int64]int64{35: {35: {}}},
			ExpectedHashString: "88201fb960ff6465",
		},

		{
			Protos: []proto.Message{
				&api.Simple{SingletonField: &api.Singleton{Singleton: &api.Singleton_TheSingleton{TheSingleton: &api.Singleton{}}}},

				&api.Singleton{Singleton: &api.Singleton_TheSingleton{TheSingleton: &api.Singleton{Singleton: &api.Singleton_TheSingleton{TheSingleton: &api.Singleton{}}}}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			EquivalentObject:   map[int64]map[int64]map[int64]int64{35: {35: {}}},
			ExpectedHashString: "661a6df2c7688a1b",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)

		checkAsSingletonOnTheWire(t, tc, hasher)
	}
}

// Checks the provided test case after all its proto messages have been cycled
// to their wire format and unmarshaled back as a Singleton message.
func checkAsSingletonOnTheWire(t *testing.T, tc ti.TestCase, hasher *ph.ProtoHasher) {
	t.Helper()

	testCaseAfterAWireTransfer := ti.TestCase{
		Protos:               tc.Protos,
		EquivalentJSONString: tc.EquivalentJSONString,
		EquivalentObject:     tc.EquivalentObject,
		ExpectedHashString:   tc.ExpectedHashString,
	}

	for i, pb := range tc.Protos {
		testCaseAfterAWireTransfer.Protos[i] = unmarshalAsSingletonOnTheWire(t, pb)
	}

	testCaseAfterAWireTransfer.Check(t, hasher)
}

// Marshals a proto message to its wire format and returns its
// unmarshaled Singleton message.
func unmarshalAsSingletonOnTheWire(t *testing.T, original proto.Message) proto.Message {
	t.Helper()

	binary, err := proto.Marshal(original)
	if err != nil {
		t.Error(err)
	}

	singleton := &api.Singleton{}
	err = proto.Unmarshal(binary, singleton)
	if err != nil {
		t.Error(err)
	}

	return singleton
}
