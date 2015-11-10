package ratelimit

import "io"

type limitedReader struct {
	r io.Reader
	b *Bucket
}

func (r *limitedReader) Read(buf []byte) (ret int, err error) {
	bytesToRead := int64(len(buf))
	n := r.b.Take(bytesToRead)
	if n == 0 {
		return 0, nil
	}
	tmpBuf := buf[0:n]
	ret, err = r.r.Read(tmpBuf)
	r.b.Return(n - int64(ret))
	return
}

// NewReader wraps an io.Reader and add transfer rate limitation on it.
func NewReader(reader io.Reader, bucket *Bucket) io.Reader {
	return &limitedReader{reader, bucket}
}
