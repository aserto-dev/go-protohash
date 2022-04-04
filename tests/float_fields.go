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
	"math"
	"testing"

	ph "github.com/aserto-dev/go-protohash"
	"github.com/aserto-dev/go-protohash/tests/api/v1"
	ti "github.com/aserto-dev/go-protohash/tests/internal"
	"google.golang.org/protobuf/proto"
)

// TestFloatFields performs tests on how floating point numbers are handled.
func TestFloatFields(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []ti.TestCase{
		/////////////////////////////
		//  Equivalence of Floats. //
		/////////////////////////////
		{
			Protos: []proto.Message{
				&api.DoubleMessage{Values: []float64{-2, -1, 0, 1, 2}},

				&api.FloatMessage{Values: []float32{-2, -1, 0, 1, 2}},
			},
			EquivalentObject:     map[string][]float64{"values": {-2, -1, 0, 1, 2}},
			EquivalentJSONString: "{\"values\": [-2, -1, 0, 1, 2]}",
			ExpectedHashString:   "3df0da89a348c288",
		},

		// Note that due to how floating point numbers work, we have to carefully
		// choose the values below in order for the decimal representation of the
		// test fractions to have 32-bit and 64-bit representations that are equal.
		{
			Protos: []proto.Message{
				&api.DoubleMessage{Values: []float64{0.0078125, 7.888609052210118e-31}},

				&api.FloatMessage{Values: []float32{0.0078125, 7.888609052210118e-31}},
			},
			EquivalentObject:     map[string][]float64{"values": {0.0078125, 7.888609052210118e-31}},
			EquivalentJSONString: "{\"values\": [0.0078125, 7.888609052210118e-31]}",
			ExpectedHashString:   "ad21263b80785c33",
		},

		{
			Protos: []proto.Message{
				&api.DoubleMessage{Values: []float64{-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},

				&api.FloatMessage{Values: []float32{-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},
			},
			EquivalentObject:     map[string][]float64{"values": {-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625}},
			EquivalentJSONString: "{\"values\": [-1.0, 1.5, 1000.000244140625, 1267650600228229401496703205376, 32.0, 13.0009765625]}",
			ExpectedHashString:   "780f95ae6c09ac4f",
		},

		/////////////////////////////////////////////////////////////////
		//  Non-equivalence of Floats using different representations. //
		/////////////////////////////////////////////////////////////////
		{
			Protos: []proto.Message{
				&api.FloatMessage{Value: 0.1},

				// A float64 "0.1" is not equal to a float32 "0.1".
				// However, float32 "0.1" is equal to float64 "1.0000000149011612e-1".
				&api.DoubleMessage{Value: 1.0000000149011612e-1},
			},
			EquivalentObject:     map[string]float32{"value": 0.1},
			EquivalentJSONString: "{\"value\": 1.0000000149011612e-1}", // JSON objecthash only uses 64-bit floats.
			ExpectedHashString:   "4b12cde041073c56",
		},

		// There's no float32 number that is equivalent to a float64 "0.1".
		{
			Protos: []proto.Message{
				&api.DoubleMessage{Value: 0.1},
			},
			EquivalentObject:     map[string]float64{"value": 0.1},
			EquivalentJSONString: "{\"value\": 0.1}",
			ExpectedHashString:   "3703e1c494c5c8e9",
		},

		{
			Protos: []proto.Message{
				&api.FloatMessage{Value: 1.2163543e+25},

				// The decimal representation of the equivalent 64-bit float is different.
				&api.DoubleMessage{Value: 1.2163543234531120e+25},
			},
			EquivalentObject:     map[string]float32{"value": 1.2163543e+25},
			EquivalentJSONString: "{\"value\": 1.2163543234531120e+25}", // JSON objecthash only uses 64-bit floats.
			ExpectedHashString:   "8cf0a1724e22238c",
		},

		// There's no float32 number that is equivalent to a float64 "1e+25".
		{
			Protos: []proto.Message{
				&api.DoubleMessage{Value: 1e+25},
			},
			EquivalentObject:     map[string]float64{"value": 1e+25},
			EquivalentJSONString: "{\"value\": 1e+25}",
			ExpectedHashString:   "b0c41f10e61fe56a",
		},

		//////////////////////
		//  Special values. //
		//////////////////////
		{
			Protos: []proto.Message{
				// Proto3 zero values are indistinguishable from unset values.
			},
			EquivalentObject:     map[string]float64{"value": 0},
			EquivalentJSONString: "{\"value\":0}",
			ExpectedHashString:   "0",
		},

		{
			Protos: []proto.Message{
				&api.DoubleMessage{Value: math.NaN()},
			},
			EquivalentObject: map[string]float64{"value": math.NaN()},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			ExpectedHashString: "e1ccffaf73d7415e",
		},
		{
			Protos: []proto.Message{
				&api.FloatMessage{Value: float32(math.NaN())},
			},
			EquivalentObject: map[string]float64{"value": math.NaN()},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			ExpectedHashString: "93d5337cc141c0de",
		},

		{
			Protos: []proto.Message{
				&api.DoubleMessage{Value: math.Inf(1)},

				&api.FloatMessage{Value: float32(math.Inf(1))},
			},
			EquivalentObject: map[string]float64{"value": math.Inf(1)},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			ExpectedHashString: "989db45457ac5739",
		},

		{
			Protos: []proto.Message{
				&api.DoubleMessage{Value: math.Inf(-1)},

				&api.FloatMessage{Value: float32(math.Inf(-1))},
			},
			EquivalentObject: map[string]float64{"value": math.Inf(-1)},
			// No equivalent JSON: JSON does not support special float values.
			// See: https://tools.ietf.org/html/rfc4627#section-2.4
			ExpectedHashString: "1057aef55fca7774",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
