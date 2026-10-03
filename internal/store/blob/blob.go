// Package blob holds the binary formats of the run store.
//
// Every blob starts with a kind byte and a version byte. Numbers are unsigned
// varints, line numbers are stored as the difference to the previous line, and
// strings carry a length prefix. The encoders are deterministic: the same data
// always gives the same bytes, so the same hash, which is what lets runs share
// the blobs of files that did not change.
package blob

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
)

// Kind is the first byte of every blob.
type Kind byte

const (
	// 1 was the source text, which the version control system keeps
	KindAnalysis Kind = 2 // methods, complexity and statements of one file version
	KindCoverage Kind = 3 // coverable lines, hits and per-report hits of one file
	KindManifest Kind = 4 // one entry per file of a run
	KindDiff     Kind = 5 // the changed lines of a run
)

func (k Kind) String() string {
	switch k {
	case KindAnalysis:
		return "analysis"
	case KindCoverage:
		return "coverage"
	case KindManifest:
		return "manifest"
	case KindDiff:
		return "diff"
	default:
		return fmt.Sprintf("kind(%d)", byte(k))
	}
}

// Valid reports whether k is a kind the store keeps.
func (k Kind) Valid() bool {
	return k >= KindAnalysis && k <= KindDiff
}

// HashSize is the length of a blob key: the first 16 bytes of SHA-256.
const HashSize = 16

type Hash [HashSize]byte

// Sum is the key of data: a blob is keyed by its encoded bytes.
func Sum(data []byte) Hash {
	full := sha256.Sum256(data)
	var h Hash
	copy(h[:], full[:HashSize])
	return h
}

func (h Hash) String() string { return hex.EncodeToString(h[:]) }

func (h Hash) IsZero() bool { return h == Hash{} }

func (h Hash) MarshalText() ([]byte, error) { return []byte(h.String()), nil }

func (h *Hash) UnmarshalText(text []byte) error {
	parsed, err := ParseHash(string(text))
	if err != nil {
		return err
	}
	*h = parsed
	return nil
}

func ParseHash(s string) (Hash, error) {
	var h Hash
	if len(s) != 2*HashSize {
		return h, fmt.Errorf("blob hash %q: want %d hex characters", s, 2*HashSize)
	}
	if _, err := hex.Decode(h[:], []byte(s)); err != nil {
		return h, fmt.Errorf("blob hash %q: %w", s, err)
	}
	return h, nil
}

// ErrCorrupt wraps every decode failure.
var ErrCorrupt = errors.New("corrupt blob")

func corrupt(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrCorrupt, fmt.Sprintf(format, args...))
}

// PeekKind returns the kind byte of an encoded blob.
func PeekKind(data []byte) (Kind, error) {
	if len(data) < 2 {
		return 0, corrupt("blob shorter than its header")
	}
	return Kind(data[0]), nil
}

type encoder struct {
	buf []byte
}

func newEncoder(kind Kind, version byte, sizeHint int) *encoder {
	e := &encoder{buf: make([]byte, 0, sizeHint+2)}
	e.buf = append(e.buf, byte(kind), version)
	return e
}

func (e *encoder) uvarint(v uint64) { e.buf = binary.AppendUvarint(e.buf, v) }

func (e *encoder) varint(v int64) { e.buf = binary.AppendVarint(e.buf, v) }

// count writes a non-negative int; negative values are clamped to zero.
func (e *encoder) count(v int) {
	if v < 0 {
		v = 0
	}
	e.uvarint(uint64(v))
}

func (e *encoder) bytes(b []byte) {
	e.uvarint(uint64(len(b)))
	e.buf = append(e.buf, b...)
}

func (e *encoder) string(s string) {
	e.uvarint(uint64(len(s)))
	e.buf = append(e.buf, s...)
}

func (e *encoder) hash(h Hash) { e.buf = append(e.buf, h[:]...) }

// path writes a path sharing a prefix with the previous one, which keeps sorted
// path lists small.
func (e *encoder) path(prev, cur string) {
	shared := commonPrefix(prev, cur)
	e.uvarint(uint64(shared))
	e.string(cur[shared:])
}

type decoder struct {
	data []byte
	pos  int
}

func newDecoder(data []byte, kind Kind, maxVersion byte) (*decoder, byte, error) {
	if len(data) < 2 {
		return nil, 0, corrupt("blob shorter than its header")
	}
	if Kind(data[0]) != kind {
		return nil, 0, corrupt("kind %s, want %s", Kind(data[0]), kind)
	}
	version := data[1]
	if version == 0 || version > maxVersion {
		return nil, 0, corrupt("%s version %d is not supported (newest %d)", kind, version, maxVersion)
	}
	return &decoder{data: data, pos: 2}, version, nil
}

func (d *decoder) remaining() int { return len(d.data) - d.pos }

func (d *decoder) done() error {
	if d.pos != len(d.data) {
		return corrupt("%d trailing bytes", len(d.data)-d.pos)
	}
	return nil
}

func (d *decoder) uvarint() (uint64, error) {
	v, n := binary.Uvarint(d.data[d.pos:])
	if err := d.checkVarint(n); err != nil {
		return 0, err
	}
	d.pos += n
	return v, nil
}

func (d *decoder) varint() (int64, error) {
	v, n := binary.Varint(d.data[d.pos:])
	if err := d.checkVarint(n); err != nil {
		return 0, err
	}
	d.pos += n
	return v, nil
}

// checkVarint rejects malformed varints and padded ones like 0x80 0x00, which
// decode to the same number as a shorter encoding.
func (d *decoder) checkVarint(n int) error {
	if n <= 0 {
		return corrupt("bad varint at byte %d", d.pos)
	}
	if n > 1 && d.data[d.pos+n-1] == 0 {
		return corrupt("padded varint at byte %d", d.pos)
	}
	return nil
}

// count reads a non-negative int that fits the platform int.
func (d *decoder) count() (int, error) {
	v, err := d.uvarint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxInt32 {
		return 0, corrupt("value %d out of range", v)
	}
	return int(v), nil
}

// length reads the size of a list whose items take at least minItemSize bytes,
// so a corrupt length cannot make the decoder allocate more than the input.
func (d *decoder) length(minItemSize int) (int, error) {
	n, err := d.count()
	if err != nil {
		return 0, err
	}
	if n*minItemSize > d.remaining() {
		return 0, corrupt("list of %d items does not fit in %d bytes", n, d.remaining())
	}
	return n, nil
}

func (d *decoder) raw(n int) ([]byte, error) {
	if n < 0 || n > d.remaining() {
		return nil, corrupt("need %d bytes, have %d", n, d.remaining())
	}
	b := d.data[d.pos : d.pos+n]
	d.pos += n
	return b, nil
}

func (d *decoder) string() (string, error) {
	n, err := d.length(1)
	if err != nil {
		return "", err
	}
	b, err := d.raw(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (d *decoder) hash() (Hash, error) {
	var h Hash
	b, err := d.raw(HashSize)
	if err != nil {
		return h, err
	}
	copy(h[:], b)
	return h, nil
}

// path reads a path written by encoder.path and checks the list stays sorted
// and free of duplicates, so every list has exactly one encoding.
func (d *decoder) path(prev string, first bool) (string, error) {
	shared, err := d.count()
	if err != nil {
		return "", err
	}
	if shared > len(prev) {
		return "", corrupt("shared prefix %d longer than previous path", shared)
	}
	suffix, err := d.string()
	if err != nil {
		return "", err
	}
	cur := prev[:shared] + suffix
	if !first && cur <= prev {
		return "", corrupt("paths not sorted: %q after %q", cur, prev)
	}
	if shared != commonPrefix(prev, cur) {
		return "", corrupt("shared prefix of %q is not minimal", cur)
	}
	return cur, nil
}

func commonPrefix(a, b string) int {
	n := min(len(a), len(b))
	for i := range n {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}
