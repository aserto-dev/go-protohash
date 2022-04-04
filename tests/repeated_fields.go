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

// TestRepeatedFields performs tests on how repeated fields are handled.
func TestRepeatedFields(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []tc.TestCase{
		///////////////////
		//  Empty lists. //
		///////////////////

		// Empty repeated fields are ignored when taking a protobuf's objecthash.
		// This is the case for both Proto2 and Proto3.
		{
			Protos: []proto.Message{
				&api.Repetitive{
					BoolField:       []bool{},
					BytesField:      [][]byte{},
					DoubleField:     []float64{},
					Fixed32Field:    []uint32{},
					Fixed64Field:    []uint64{},
					FloatField:      []float32{},
					Int32Field:      []int32{},
					Int64Field:      []int64{},
					Sfixed32Field:   []int32{},
					Sfixed64Field:   []int64{},
					Sint32Field:     []int32{},
					Sint64Field:     []int64{},
					StringField:     []string{},
					Uint32Field:     []uint32{},
					Uint64Field:     []uint64{},
					SimpleField:     []*api.Simple{},
					RepetitiveField: []*api.Repetitive{},
					SingletonField:  []*api.Singleton{},
				},
			},
			EquivalentJSONString: "{}",
			EquivalentObject:     map[string]interface{}{},
			ExpectedHashString:   "0",
		},

		//////////////////////////
		//  Lists with strings. //
		//////////////////////////
		{
			Protos: []proto.Message{
				&api.Repetitive{StringField: []string{""}},
			},
			EquivalentJSONString: "{\"string_field\": [\"\"]}",
			EquivalentObject:     map[string][]string{"string_field": {""}},
			ExpectedHashString:   "bab48eecfa8cd51a",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{StringField: []string{"foo"}},
			},
			EquivalentJSONString: "{\"string_field\": [\"foo\"]}",
			EquivalentObject:     map[string][]string{"string_field": {"foo"}},
			ExpectedHashString:   "e781d93648f4e29b",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{StringField: []string{"foo", "bar"}},
			},
			EquivalentJSONString: "{\"string_field\": [\"foo\", \"bar\"]}",
			EquivalentObject:     map[string][]string{"string_field": {"foo", "bar"}},
			ExpectedHashString:   "5e398a810a1e8af7",
		},

		///////////////////////
		//  Lists with ints. //
		///////////////////////

		// JSON treats all numbers as floats, so it is not possible to have an equivalent JSON string.

		{
			Protos: []proto.Message{
				&api.Repetitive{Int64Field: []int64{0}},
			},
			EquivalentObject:   map[string][]int64{"int64_field": {0}},
			ExpectedHashString: "88abed3eda001f87",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{Int64Field: []int64{-2, -1, 0, 1, 2}},
			},
			EquivalentObject:   map[string][]int64{"int64_field": {-2, -1, 0, 1, 2}},
			ExpectedHashString: "8e40d97221a0dba3",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{Int64Field: []int64{123456789012345, 678901234567890}},
			},
			EquivalentObject:   map[string][]int64{"int64_field": {123456789012345, 678901234567890}},
			ExpectedHashString: "e7c4423fe65d2f08",
		},

		/////////////////////////
		//  Lists with floats. //
		/////////////////////////
		{
			Protos: []proto.Message{
				&api.Repetitive{FloatField: []float32{0}},
			},
			EquivalentJSONString: "{\"float_field\": [0]}",
			EquivalentObject:     map[string][]float32{"float_field": {0}},
			ExpectedHashString:   "88abed3eda001f87",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{FloatField: []float32{0.0}},
			},
			EquivalentJSONString: "{\"float_field\": [0.0]}",
			EquivalentObject:     map[string][]float32{"float_field": {0.0}},
			ExpectedHashString:   "88abed3eda001f87",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{FloatField: []float32{-2, -1, 0, 1, 2}},
			},
			EquivalentJSONString: "{\"float_field\": [-2, -1, 0, 1, 2]}",
			EquivalentObject:     map[string][]float32{"float_field": {-2, -1, 0, 1, 2}},
			ExpectedHashString:   "3df0da89a348c288",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{FloatField: []float32{1, 2, 3}},
			},
			EquivalentJSONString: "{\"float_field\": [1, 2, 3]}",
			EquivalentObject:     map[string][]float32{"float_field": {1, 2, 3}},
			ExpectedHashString:   "96c4f986cedc148",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{DoubleField: []float64{1.2345, -10.1234}},
			},
			EquivalentJSONString: "{\"double_field\": [1.2345, -10.1234]}",
			EquivalentObject:     map[string][]float64{"double_field": {1.2345, -10.1234}},
			ExpectedHashString:   "d317c8afdac508cc",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{DoubleField: []float64{1.0, 1.5, 0.0001, 1000.9999999, 2.0, -23.1234, 2.32542}},
			},
			EquivalentJSONString: "{\"double_field\": [1.0, 1.5, 0.0001, 1000.9999999, 2.0, -23.1234, 2.32542]}",
			EquivalentObject:     map[string][]float64{"double_field": {1.0, 1.5, 0.0001, 1000.9999999, 2.0, -23.1234, 2.32542}},
			ExpectedHashString:   "ca0d702cfcb510b9",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{DoubleField: []float64{123456789012345, 678901234567890}},
			},
			EquivalentJSONString: "{\"double_field\": [123456789012345, 678901234567890]}",
			EquivalentObject:     map[string][]float64{"double_field": {123456789012345, 678901234567890}},
			ExpectedHashString:   "66ed2a9a6f6b8684",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
