package wal

import (
	"context"
	"io"
)

func CopyWithContext(ctx context.Context, dst io.Writer, src io.Reader, buffer []byte) (int64, error) {
	if len(buffer) == 0 {
		buffer = make([]byte, 32*1024)
	}
	var written int64
	for {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		n, readErr := src.Read(buffer)
		if n > 0 {
			m, writeErr := dst.Write(buffer[:n])
			written += int64(m)
			if writeErr != nil {
				return written, writeErr
			}
			if m != n {
				return written, io.ErrShortWrite
			}
		}
		if readErr == io.EOF {
			return written, nil
		}
		if readErr != nil {
			return written, readErr
		}
	}
}
