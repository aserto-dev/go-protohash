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
	ti "github.com/aserto-dev/go-protohash/tests/internal"
	"google.golang.org/protobuf/proto"
)

// TestOtherTypes performs tests on types that do not have their own test file.
func TestOtherTypes(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []ti.TestCase{
		///////////
		//  Nil. //
		///////////
		// {
		// 	Protos: []proto.Message{
		// 		nil,
		// 	},
		// 	EquivalentJSONString: "null",
		// 	EquivalentObject:     nil,
		// 	ExpectedHashString:   "0",
		// },

		/////////////////////
		// Boolean fields. //
		/////////////////////
		{
			Protos: []proto.Message{
				&api.Simple{BoolField: true},
			},
			EquivalentJSONString: "{\"bool_field\": true}",
			EquivalentObject:     map[string]bool{"bool_field": true},
			ExpectedHashString:   "f7a206297de86dbe",
		},

		// {
		// 	Protos: []proto.Message{
		// 		&pb2_latest.Simple{BoolField: proto.Bool(false)},
		// 		// proto3 scalar fields set to their default value are considered empty.
		// 	},
		// 	EquivalentJSONString: "{\"bool_field\": false}",
		// 	EquivalentObject:     map[string]bool{"bool_field": false},
		// 	ExpectedHashString:   "1ab5ecdbe4176473024f7efd080593b740d22d076d06ea6edd8762992b484a12",
		// },

		///////////////////
		// Bytes fields. //
		///////////////////
		{
			Protos: []proto.Message{
				&api.Simple{BytesField: []byte{0, 0, 0}},
			},
			// No equivalent JSON: JSON does not have a "bytes" type.
			EquivalentObject:   map[string][]byte{"bytes_field": []byte("\000\000\000")},
			ExpectedHashString: "cb59b0693719a410",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
