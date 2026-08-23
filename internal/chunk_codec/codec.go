package chunk_codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"math"
)

var (
	ErrDecode        = errors.New("chunk decode failed")
	ErrShortChunk    = errors.New("short chunk")
	ErrInvalidLength = errors.New("invalid chunk length")
	ErrChecksum      = errors.New("chunk checksum")
	ErrHeader        = errors.New("chunk header")
	ErrPoint         = errors.New("chunk point")
)

type DecodeFailure struct {
	Kind error
	Err  error
}

func (e DecodeFailure) Error() string { return fmt.Sprintf("decode chunk: %v", e.Err) }

func (e DecodeFailure) Unwrap() []error { return []error{e.Kind, e.Err} }

const magic uint32 = 0x54534331

type Point struct {
	Timestamp int64
	Value     float64
	Quality   uint8
}

func Encode(points []Point) ([]byte, error) {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, magic)
	binary.Write(&b, binary.BigEndian, uint32(len(points)))
	var prev int64
	for _, p := range points {
		binary.Write(&b, binary.BigEndian, p.Timestamp-prev)
		binary.Write(&b, binary.BigEndian, p.Value)
		b.WriteByte(p.Quality)
		prev = p.Timestamp
	}
	raw := b.Bytes()
	var out bytes.Buffer
	binary.Write(&out, binary.BigEndian, uint32(len(raw)))
	binary.Write(&out, binary.BigEndian, crc32.ChecksumIEEE(raw))
	out.Write(raw)
	return out.Bytes(), nil
}
func Decode(data []byte) ([]Point, error) {
	if len(data) < 8 {
		return nil, DecodeFailure{Kind: ErrDecode, Err: ErrShortChunk}
	}
	n := binary.BigEndian.Uint32(data)
	sum := binary.BigEndian.Uint32(data[4:])
	if int(n) != len(data)-8 {
		return nil, DecodeFailure{Kind: ErrDecode, Err: ErrInvalidLength}
	}
	raw := data[8:]
	if crc32.ChecksumIEEE(raw) != sum {
		return nil, DecodeFailure{Kind: ErrDecode, Err: ErrChecksum}
	}
	rd := bytes.NewReader(raw)
	var m, c uint32
	binary.Read(rd, binary.BigEndian, &m)
	binary.Read(rd, binary.BigEndian, &c)
	if m != magic || c > 10_000_000 {
		return nil, DecodeFailure{Kind: ErrDecode, Err: ErrHeader}
	}
	out := make([]Point, 0, c)
	var prev int64
	for i := uint32(0); i < c; i++ {
		var d int64
		var v uint64
		var q uint8
		if binary.Read(rd, binary.BigEndian, &d) != nil || binary.Read(rd, binary.BigEndian, &v) != nil {
			return nil, DecodeFailure{Kind: ErrDecode, Err: ErrPoint}
		}
		q, e := rd.ReadByte()
		if e != nil {
			return nil, e
		}
		prev += d
		out = append(out, Point{prev, math.Float64frombits(v), q})
	}
	return out, nil
}
