package protohash_test

import (
	"hash/fnv"
	"testing"

	"github.com/aserto-dev/go-protohash"
	"github.com/aserto-dev/go-protohash/tests"
)

func TestFunctional(t *testing.T) {
	h := fnv.New64a()
	ph := protohash.New(protohash.WithHash64(h))

	// t.Run("TestBadness", func(t *testing.T) { tests.TestBadness(t, ph) })
	t.Run("TestEmptyFields", func(t *testing.T) { tests.TestEmptyFields(t, ph) })
	t.Run("TestFloatFields", func(t *testing.T) { tests.TestFloatFields(t, ph) })
	t.Run("TestIntegerFields", func(t *testing.T) { tests.TestIntegerFields(t, ph) })
	t.Run("TestMaps", func(t *testing.T) { tests.TestMaps(t, ph) })
	t.Run("TestOneOfFields", func(t *testing.T) { tests.TestOneOfFields(t, ph) })
	t.Run("TestOtherTypes", func(t *testing.T) { tests.TestOtherTypes(t, ph) })
	// t.Run("TestProto2DefaultFieldValues", func(t *testing.T) { tests.TestProto2DefaultFieldValues(t, ph) })
	t.Run("TestRepeatedFields", func(t *testing.T) { tests.TestRepeatedFields(t, ph) })
	t.Run("TestStringFields", func(t *testing.T) { tests.TestStringFields(t, ph) })

	// Well-known types.
	// t.Run("TestTimestamps", func(t *testing.T) { wkt.TestTimestamps(t, ph) })
	// t.Run("TestUnsupportedWellKnownTypes", func(t *testing.T) { wkt.TestUnsupportedWellKnownTypes(t, ph) })
}
