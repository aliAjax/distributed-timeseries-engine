package wal

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"sync"
)

var ErrCorrupt = errors.New("corrupt wal record")
var ErrReplay = errors.New("wal replay failed")

type ReplayFailure struct {
	Kind error
	Err  error
}

func (e ReplayFailure) Error() string { return fmt.Sprintf("replay: %v", e.Err) }

func (e ReplayFailure) Unwrap() []error { return []error{e.Kind, e.Err} }

type Record struct {
	Sequence uint64 `json:"sequence"`
	Kind     string `json:"kind"`
	Payload  []byte `json:"payload"`
}
type Log struct {
	mu   sync.Mutex
	file *os.File
	next uint64
	path string
}

func Open(path string) (*Log, error) {
	f, e := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if e != nil {
		return nil, e
	}
	return &Log{file: f, next: 1, path: path}, nil
}
func (l *Log) Append(r Record, syncWrite bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if r.Sequence == 0 {
		r.Sequence = l.next
	}
	if r.Sequence >= l.next {
		l.next = r.Sequence + 1
	}
	b, e := json.Marshal(r)
	if e != nil {
		return e
	}
	var hdr [8]byte
	binary.BigEndian.PutUint32(hdr[:4], uint32(len(b)))
	binary.BigEndian.PutUint32(hdr[4:], crc32.ChecksumIEEE(b))
	if _, e = l.file.Write(hdr[:]); e != nil {
		return e
	}
	if _, e = l.file.Write(b); e != nil {
		return e
	}
	if syncWrite {
		return l.file.Sync()
	}
	return nil
}
func (l *Log) Replay(fn func(Record) error) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, e := l.file.Seek(0, 0); e != nil {
		return e
	}
	rd := bufio.NewReader(l.file)
	var good int64
	for {
		var hdr [8]byte
		_, e := io.ReadFull(rd, hdr[:])
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		n := binary.BigEndian.Uint32(hdr[:4])
		sum := binary.BigEndian.Uint32(hdr[4:])
		if n == 0 || n > 64<<20 {
			return ReplayFailure{Kind: ErrReplay, Err: ErrCorrupt}
		}
		b := make([]byte, n)
		if _, e = io.ReadFull(rd, b); e != nil {
			_ = l.truncate(good)
			break
		}
		if crc32.ChecksumIEEE(b) != sum {
			_ = l.truncate(good)
			break
		}
		var r Record
		if e = json.Unmarshal(b, &r); e != nil {
			return e
		}
		if e = fn(r); e != nil {
			return e
		}
		good += int64(8 + n)
		if r.Sequence >= l.next {
			l.next = r.Sequence + 1
		}
	}
	_, e := l.file.Seek(0, io.SeekEnd)
	return e
}
func (l *Log) truncate(n int64) error { return l.file.Truncate(n) }
func (l *Log) Close() error           { l.mu.Lock(); defer l.mu.Unlock(); return l.file.Close() }
func (l *Log) Path() string           { return l.path }
