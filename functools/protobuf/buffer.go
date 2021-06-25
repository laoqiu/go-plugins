package protobuf

// Constants that identify the encoding of a value on the wire.
const (
	WireVarint     = 0
	WireFixed64    = 1
	WireBytes      = 2
	WireStartGroup = 3
	WireEndGroup   = 4
	WireFixed32    = 5
)

type Buffer struct {
	buf   []byte // encode/decode byte stream
	index int    // read point

	deterministic bool
}

// NewBuffer allocates a new Buffer and initializes its internal data to
// the contents of the argument slice.
func NewBuffer(e []byte) *Buffer {
	return &Buffer{buf: e}
}

// Reset resets the Buffer, ready for marshaling a new protocol buffer.
func (p *Buffer) Reset() {
	p.buf = p.buf[0:0] // for reading/writing
	p.index = 0        // for reading
}

// SetBuf replaces the internal buffer with the slice,
// ready for unmarshaling the contents of the slice.
func (p *Buffer) SetBuf(s []byte) {
	p.buf = s
	p.index = 0
}

// Bytes returns the contents of the Buffer.
func (p *Buffer) Bytes() []byte { return p.buf }

// SetDeterministic sets whether to use deterministic serialization.
//
// Deterministic serialization guarantees that for a given binary, equal
// messages will always be serialized to the same bytes. This implies:
//
//   - Repeated serialization of a message will return the same bytes.
//   - Different processes of the same binary (which may be executing on
//     different machines) will serialize equal messages to the same bytes.
//
// Note that the deterministic serialization is NOT canonical across
// languages. It is not guaranteed to remain stable over time. It is unstable
// across different builds with schema changes due to unknown fields.
// Users who need canonical serialization (e.g., persistent storage in a
// canonical form, fingerprinting, etc.) should define their own
// canonicalization specification and implement their own serializer rather
// than relying on this API.
//
// If deterministic serialization is requested, map entries will be sorted
// by keys in lexographical order. This is an implementation detail and
// subject to change.
func (p *Buffer) SetDeterministic(deterministic bool) {
	p.deterministic = deterministic
}
