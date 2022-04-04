package tests

import (
	"testing"

	ph "github.com/aserto-dev/go-protohash"
	"github.com/aserto-dev/go-protohash/tests/api/v1"
	tc "github.com/aserto-dev/go-protohash/tests/internal"
	"google.golang.org/protobuf/proto"
)

// TestStringFields performs tests on how strings are handled.
func TestStringFields(t *testing.T, hasher *ph.ProtoHasher) {

	testCases := []tc.TestCase{
		{
			Protos: []proto.Message{
				&api.Simple{StringField: "你好"},
			},
			ExpectedHashString: "e2dd2a3d97f401ac",
		},

		{
			Protos: []proto.Message{
				&api.Simple{StringField: "\u03d3"},
			},
			EquivalentObject:     map[string]string{"string_field": "\u03d3"},
			EquivalentJSONString: "{\"string_field\":\"\u03d3\"}",
			ExpectedHashString:   "889bf3c60923cb21",
		},

		// Note that this is the same character as above, but hashes differently
		// unless unicode normalisation is applied.
		{
			Protos: []proto.Message{
				&api.Simple{StringField: "\u03d2\u0301"},
			},
			EquivalentObject:     map[string]string{"string_field": "\u03d2\u0301"},
			EquivalentJSONString: "{\"string_field\":\"\u03d2\u0301\"}",
			ExpectedHashString:   "6b17ebd06d7ba11d",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{StringField: []string{""}},
			},
			EquivalentObject:     map[string][]string{"string_field": {""}},
			EquivalentJSONString: "{\"string_field\":[\"\"]}",
			ExpectedHashString:   "bab48eecfa8cd51a",
		},

		{
			Protos: []proto.Message{
				&api.Repetitive{StringField: []string{"", "Test", "你好", "\u03d3"}},
			},
			EquivalentObject:     map[string][]string{"string_field": {"", "Test", "你好", "\u03d3"}},
			EquivalentJSONString: "{\"string_field\":[\"\",\"Test\",\"你好\",\"\u03d3\"]}",
			ExpectedHashString:   "ee1bfab42da1d7fe",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
