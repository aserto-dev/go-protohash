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

// TestEmptyFields checks that empty proto fields are handled properly.
func TestEmptyFields(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []ti.TestCase{
		{
			Protos: []proto.Message{
				&api.Empty{},

				// Empty repeated fields are ignored.
				&api.Repetitive{StringField: []string{}},

				// Empty map fields are ignored.
				&api.StringMaps{StringToString: map[string]string{}},

				// Proto3 scalar fields set to their default values are considered empty.
				&api.Simple{BoolField: false},
				&api.Simple{BytesField: []byte{}},
				&api.Simple{DoubleField: 0},
				&api.Simple{DoubleField: 0.0},
				&api.Simple{Fixed32Field: 0},
				&api.Simple{Fixed64Field: 0},
				&api.Simple{FloatField: 0},
				&api.Simple{FloatField: 0.0},
				&api.Simple{Int32Field: 0},
				&api.Simple{Int64Field: 0},
				&api.Simple{Sfixed32Field: 0},
				&api.Simple{Sfixed64Field: 0},
				&api.Simple{Sint32Field: 0},
				&api.Simple{Sint64Field: 0},
				&api.Simple{StringField: ""},
				&api.Simple{Uint32Field: 0},
				&api.Simple{Uint64Field: 0},
			},
			EquivalentJSONString: "{}",
			EquivalentObject:     map[string]interface{}{},
			ExpectedHashString:   "0",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
