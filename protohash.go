package protohash

import (
	"encoding/binary"
	"hash"
	"hash/fnv"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func New(opts ...HashOption) *ProtoHasher {
	ph := ProtoHasher{
		h: fnv.New64a(),
	}

	for _, opt := range opts {
		opt(&ph)
	}

	return &ph
}

type HashOption func(*ProtoHasher)

func WithHash64(h hash.Hash64) HashOption {
	return func(ph *ProtoHasher) {
		ph.h = h
	}
}

type ProtoHasher struct {
	h hash.Hash64
}

func (ph *ProtoHasher) HashMessage(msg proto.Message) (uint64, error) {
	if msg == nil {
		return 0, status.Error(codes.InvalidArgument, "msg is nil")
	}

	m := msg.ProtoReflect()
	if !m.IsValid() {
		return 0, status.Error(codes.FailedPrecondition, "msg is invalid")
	}

	return ph.hashMessage(m)
}

func (ph *ProtoHasher) hashMessage(msg protoreflect.Message) (uint64, error) {
	var (
		h uint64
		e error
	)

	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		hv, err := ph.hashField(fd, v)
		if err != nil {
			e = err
			return false
		}

		h, err = hashUpdateOrdered(ph.h, h, hv)
		if err != nil {
			e = err
			return false
		}

		return true
	})
	if e != nil {
		return 0, e
	}
	return h, nil
}

func (ph *ProtoHasher) hashField(fd protoreflect.FieldDescriptor, v protoreflect.Value) (uint64, error) {
	switch {
	case fd.IsList():
		return ph.hashList(fd, v.List())
	case fd.IsMap():
		return ph.hashMap(fd, v.Map())
	default:
		return ph.hashValue(fd, v)
	}
}

func (ph *ProtoHasher) hashMap(fd protoreflect.FieldDescriptor, v protoreflect.Map) (uint64, error) {
	var (
		h uint64
		e error
	)
	v.Range(func(k protoreflect.MapKey, vx protoreflect.Value) bool {
		hk, keyErr := ph.hashValue(fd.MapKey(), k.Value())
		if keyErr != nil {
			e = keyErr
			return false
		}

		hv, valErr := ph.hashValue(fd.MapValue(), vx)
		if valErr != nil {
			e = valErr
			return false
		}

		fieldHash, err := hashUpdateOrdered(ph.h, hk, hv)
		if err != nil {
			e = err
			return false
		}

		h = hashUpdateUnordered(h, fieldHash)

		return true
	})
	if e != nil {
		return 0, e
	}

	return hashFinishUnordered(ph.h, h)
}

func (ph *ProtoHasher) hashList(fd protoreflect.FieldDescriptor, v protoreflect.List) (uint64, error) {
	var h uint64
	for i := v.Len() - 1; i >= 0; i-- {
		hv, err := ph.hashValue(fd, v.Get(i))
		if err != nil {
			return 0, err
		}

		h, err = hashUpdateOrdered(ph.h, h, hv)
		if err != nil {
			return 0, err
		}
	}
	return h, nil
}

func (ph *ProtoHasher) hashValue(fd protoreflect.FieldDescriptor, v protoreflect.Value) (uint64, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		var tmp int8
		if v.Bool() {
			tmp = 1
		}
		ph.h.Reset()
		err := binary.Write(ph.h, binary.LittleEndian, tmp)
		return ph.h.Sum64(), err

	case protoreflect.EnumKind:
		ph.h.Reset()
		err := binary.Write(ph.h, binary.LittleEndian, v.Enum())
		return ph.h.Sum64(), err

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		ph.h.Reset()
		err := binary.Write(ph.h, binary.LittleEndian, v.Int())
		return ph.h.Sum64(), err

	case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		ph.h.Reset()
		err := binary.Write(ph.h, binary.LittleEndian, v.Uint())
		return ph.h.Sum64(), err

	case protoreflect.FloatKind, protoreflect.DoubleKind:
		ph.h.Reset()
		err := binary.Write(ph.h, binary.LittleEndian, v.Float())
		return ph.h.Sum64(), err

	case protoreflect.StringKind:
		ph.h.Reset()
		_, err := ph.h.Write([]byte(v.String()))
		return ph.h.Sum64(), err

	case protoreflect.BytesKind:
		ph.h.Reset()
		err := binary.Write(ph.h, binary.LittleEndian, v.Bytes())
		return ph.h.Sum64(), err

	case protoreflect.MessageKind, protoreflect.GroupKind:
		return ph.hashMessage(v.Message())

	default:
		return 0, errors.Errorf("unknown kind to hash: %s", fd.Kind())
	}
}

// hashUpdateUnordered
// Adaopted for protomsg from https://github.com/mitchellh/hashstructure
func hashUpdateUnordered(a, b uint64) uint64 {
	return a ^ b
}

// hashUpdateOrdered
// Adaopted for protomsg from https://github.com/mitchellh/hashstructure
func hashUpdateOrdered(h hash.Hash64, a, b uint64) (uint64, error) {
	// For ordered updates, use a real hash function
	h.Reset()

	// We just panic if the binary writes fail because we are writing
	// an int64 which should never be fail-able.
	e1 := binary.Write(h, binary.LittleEndian, a)
	e2 := binary.Write(h, binary.LittleEndian, b)
	if e1 != nil {
		return 0, e1
	}
	if e2 != nil {
		return 0, e2
	}

	return h.Sum64(), nil
}

// hashFinishUnordered
// Adaopted for protomsg from https://github.com/mitchellh/hashstructure
//
// After mixing a group of unique hashes with hashUpdateUnordered, it's always
// necessary to call hashFinishUnordered. Why? Because hashUpdateUnordered
// is a simple XOR, and calling hashUpdateUnordered on hashes produced by
// hashUpdateUnordered can effectively cancel out a previous change to the hash
// result if the same hash value appears later on. For example, consider:
//
//   hashUpdateUnordered(hashUpdateUnordered("A", "B"), hashUpdateUnordered("A", "C")) =
//   H("A") ^ H("B")) ^ (H("A") ^ H("C")) =
//   (H("A") ^ H("A")) ^ (H("B") ^ H(C)) =
//   H(B) ^ H(C) =
//   hashUpdateUnordered(hashUpdateUnordered("Z", "B"), hashUpdateUnordered("Z", "C"))
//
// hashFinishUnordered "hardens" the result, so that encountering partially
// overlapping input data later on in a different context won't cancel out.
func hashFinishUnordered(h hash.Hash64, a uint64) (uint64, error) {
	h.Reset()

	err := binary.Write(h, binary.LittleEndian, a)
	if err != nil {
		return 0, err
	}

	return h.Sum64(), nil
}
