package ratelimit

import "io"

type limitedWriter struct {
	w io.Writer
	b *Bucket
}

func (w *limitedWriter) Write(buf []byte) (written int, err error) {
	bytesToWrite := int64(len(buf))
	for bytesToWrite > 0 {
		n := w.b.Take(bytesToWrite)
		tmpBuf := buf[written : int64(written)+n]
		ret, err := w.w.Write(tmpBuf)
		if err != nil {
			return written, err
		}
		bytesToWrite -= int64(ret)
		written += ret
	}
	return
}

// NewReader wraps an io.Reader and add transfer rate limitation on it.
func NewWriter(writer io.Writer, bucket *Bucket) io.Writer {
	return &limitedWriter{writer, bucket}
}
