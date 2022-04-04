package tests

import (
	"testing"

	"github.com/aserto-dev/go-protohash"
	"github.com/aserto-dev/go-protohash/tests/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestSingle(t *testing.T) {
	msg := &api.Single{
		BoolField:     true,
		BytesField:    []byte{0, 1, 2, 3, 4, 5},
		DoubleField:   float64(0),
		Fixed32Field:  uint32(1),
		Fixed64Field:  uint64(2),
		FloatField:    float32(1),
		Int32Field:    int32(1),
		Int64Field:    int64(1),
		Sfixed32Field: int32(0),
		Sfixed64Field: int64(1),
		Sint32Field:   int32(1),
		Sint64Field:   int64(1),
		Uint32Field:   uint32(1),
		Uint64Field:   uint64(1),
	}

	ph := protohash.New()
	hv, err := ph.HashMessage(msg)
	assert.NoError(t, err)
	t.Logf("hash %d", hv)
}

func TestOptionalSingleNil(t *testing.T) {
	msg := &api.OptionalSingle{
		BoolField:     nil,
		BytesField:    nil,
		DoubleField:   nil,
		Fixed32Field:  nil,
		Fixed64Field:  nil,
		FloatField:    nil,
		Int32Field:    nil,
		Int64Field:    nil,
		Sfixed32Field: nil,
		Sfixed64Field: nil,
		Sint32Field:   nil,
		Sint64Field:   nil,
		Uint32Field:   nil,
		Uint64Field:   nil,
	}

	ph := protohash.New()
	hv, err := ph.HashMessage(msg)
	assert.NoError(t, err)
	t.Logf("hash %d", hv)
}
