// Copyright 2017 The ObjectHash-Proto Authors
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
	tc "github.com/aserto-dev/go-protohash/tests/internal"
	"google.golang.org/protobuf/proto"
)

// TestIntegerFields performs tests on how integers are handled.
func TestIntegerFields(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []tc.TestCase{
		///////////////////////////////
		//  Equivalence of Integers. //
		///////////////////////////////
		{
			Protos: []proto.Message{
				&api.Fixed32Message{Values: []uint32{0, 1, 2}},
				&api.Fixed64Message{Values: []uint64{0, 1, 2}},
				&api.Int32Message{Values: []int32{0, 1, 2}},
				&api.Int64Message{Values: []int64{0, 1, 2}},
				&api.Sfixed32Message{Values: []int32{0, 1, 2}},
				&api.Sfixed64Message{Values: []int64{0, 1, 2}},
				&api.Sint32Message{Values: []int32{0, 1, 2}},
				&api.Sint64Message{Values: []int64{0, 1, 2}},
				&api.Uint32Message{Values: []uint32{0, 1, 2}},
				&api.Uint64Message{Values: []uint64{0, 1, 2}},
			},
			EquivalentObject: map[string][]int32{"values": {0, 1, 2}},
			// No equivalent JSON: JSON does not have an "integer" type. All numbers are floats.
			ExpectedHashString: "cc2ab53c181c4329",
		},

		{
			Protos: []proto.Message{
				&api.Int32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&api.Int64Message{Values: []int64{-2, -1, 0, 1, 2}},
				&api.Sfixed32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&api.Sfixed64Message{Values: []int64{-2, -1, 0, 1, 2}},
				&api.Sint32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&api.Sint64Message{Values: []int64{-2, -1, 0, 1, 2}},
			},
			EquivalentObject: map[string][]int32{"values": {-2, -1, 0, 1, 2}},
			// No equivalent JSON: JSON does not have an "integer" type. All numbers are floats.
			ExpectedHashString: "8e40d97221a0dba3",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
