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

// TestMaps performs tests on how maps are handled.
func TestMaps(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []tc.TestCase{
		////////////////////
		//  Boolean maps. //
		////////////////////
		{
			Protos: []proto.Message{
				&api.BoolMaps{BoolToString: map[bool]string{true: "NOT FALSE", false: "NOT TRUE"}},
			},
			// No equivalent JSON object because JSON map keys must be strings.
			EquivalentObject:   map[string]map[bool]string{"bool_to_string": {true: "NOT FALSE", false: "NOT TRUE"}},
			ExpectedHashString: "6f6b5869cdd9333",
		},

		////////////////////
		//  Integer maps. //
		////////////////////
		{
			Protos: []proto.Message{
				&api.IntMaps{IntToString: map[int64]string{0: "ZERO"}},
			},
			// No equivalent JSON object because JSON map keys must be strings.
			EquivalentObject:   map[string]map[int64]string{"int_to_string": {0: "ZERO"}},
			ExpectedHashString: "cb97c968692e8b24",
		},

		///////////////////
		//  String maps. //
		///////////////////
		{
			Protos: []proto.Message{
				&api.StringMaps{StringToString: map[string]string{"foo": "bar"}},
			},
			EquivalentJSONString: "{\"string_to_string\": {\"foo\": \"bar\"}}",
			EquivalentObject:     map[string]map[string]string{"string_to_string": {"foo": "bar"}},
			ExpectedHashString:   "c5f3d4ac79aa224b",
		},

		{
			Protos: []proto.Message{
				&api.StringMaps{StringToString: map[string]string{
					"": "你好", "你好": "\u03d3", "\u03d3": "\u03d2\u0301"}},
			},
			EquivalentJSONString: "{\"string_to_string\": {\"\": \"你好\", \"你好\": \"\u03d3\", \"\u03d3\": \"\u03d2\u0301\"}}",
			EquivalentObject:     map[string]map[string]string{"string_to_string": {"": "你好", "你好": "\u03d3", "\u03d3": "\u03d2\u0301"}},
			ExpectedHashString:   "fd2644e21e9d8a32",
		},

		//////////////////////////////
		//  Maps of proto messages. //
		//////////////////////////////
		{
			Protos: []proto.Message{
				&api.StringMaps{StringToSimple: map[string]*api.Simple{"foo": {}}},
			},
			EquivalentJSONString: "{\"string_to_simple\": {\"foo\": {}}}",
			EquivalentObject:     map[string]map[string]map[string]string{"string_to_simple": {"foo": {}}},
			ExpectedHashString:   "c76d7fcd4b54fd92",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
